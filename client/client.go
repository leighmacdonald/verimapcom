package client

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hpcloud/tail"
	"github.com/leighmacdonald/verimapcom/gs"
	"github.com/leighmacdonald/verimapcom/pb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Opts struct {
	RootDir    string
	ListenAddr string
}

type ProjectConfig struct {
	ProjectRoot string `yaml:"project_root"`
	ProjectID   uint32 `yaml:"project_id"`
}

type Client struct {
	Opts
	done        chan bool
	newFiles    chan string
	hotSpots    chan TimedLine
	positions   chan TimedLine
	Ctx         context.Context
	client      pb.RPCClient
	clientConn  *grpc.ClientConn
	tailedFiles []*tail.Tail
	projectID   int32
}

func (c *Client) Disconnect() {
	if err := c.clientConn.Close(); err != nil {
		log.Error("Failed to cleanly close connection: %v", err)
	}
	log.Infof("Disconnected from gRPC")
}

func (c *Client) Connect() error {
	conn, err3 := newGRPCConn(gRPCOpts{
		Tls:          false,
		ServerAddr:   c.ListenAddr,
		CaFile:       "",
		HostOverride: "",
	})
	if err3 != nil {
		return errors.Wrap(err3, "Failed to create client connection")
	}
	c.clientConn = conn
	c.client = pb.NewRPCClient(conn)
	log.Println("Connected to gRPC")
	return nil
}

// 0                                 1           2             3         4        5        6         7
// frame_4_time_1560859859.3455.tiff,56.23031677,-117.44821497,572.43188,-0.39258,-2.01074,327.12402,-0.39254
func (c *Client) processPosition() error {
	stream, err := c.client.SendPosition(c.Ctx, grpc.UseCompressor(gzip.Name))
	if err != nil {
		return errors.Wrap(err, "Failed to setup stream")
	}
	var unsent []*pb.PositionEvent
	for {
		select {
		case line := <-c.positions:
			log.Debugf(line.Text)
			columns := strings.Split(line.Text, ",")
			pcs := strings.Split(columns[0], "_")
			if len(columns) != 8 {
				continue
			}
			if len(pcs) != 5 {
				continue
			}
			tStr := strings.ReplaceAll(pcs[4], ".raw", "")
			pcsTime := strings.Split(tStr, ".")
			if len(pcsTime) != 2 {
				continue
			}
			t0, err := strconv.ParseInt(pcsTime[0], 10, 64)
			if err != nil {
				continue
			}
			t1, err := strconv.ParseInt(pcsTime[1], 10, 64)
			if err != nil {
				continue
			}
			lat, err := strconv.ParseFloat(columns[1], 64)
			if err != nil {
				continue
			}
			lon, err := strconv.ParseFloat(columns[2], 64)
			if err != nil {
				continue
			}
			elevation, err := strconv.ParseFloat(columns[3], 64)
			if err != nil {
				continue
			}
			req := &pb.PositionEvent{
				MissionId: c.projectID,
				At:        &timestamp.Timestamp{Seconds: t0, Nanos: int32(t1 * 1000000)},
				Location:  &pb.Location{Lat: lat, Lon: lon},
				Elevation: int32(elevation),
			}
			if err := stream.Send(req); err != nil {
				log.Errorf("Could not send position: %s", err)
				unsent = append(unsent, req)
			}
		}
	}
}

func (c *Client) processHotcluster() error {
	stream, err := c.client.SendHotspot(c.Ctx)
	if err != nil {
		return err
	}
	var unsent []*pb.HotSpotEvent
	for {
		select {
		case line := <-c.hotSpots:
			log.Debugf(line.Text)
			columns := strings.Split(line.Text, ",")
			id, err := strconv.ParseInt(columns[0], 10, 64)
			if err != nil {
				continue
			}
			lat, err := strconv.ParseFloat(columns[1], 64)
			if err != nil {
				continue
			}
			lon, err := strconv.ParseFloat(columns[2], 64)
			if err != nil {
				continue
			}
			delta, err := strconv.ParseFloat(columns[3], 64)
			if err != nil {
				continue
			}
			req := &pb.HotSpotEvent{
				MissionId: c.projectID,
				Id:        id,
				Location:  &pb.Location{Lat: lat, Lon: lon},
				Delta:     float32(delta),
			}
			if err := stream.Send(req); err != nil {
				log.Errorf("Could not send position: %s", err)
				unsent = append(unsent, req)
			}
		}
	}
}

func (c *Client) Ping() error {
	t0 := time.Now()
	ctx, cancel := context.WithDeadline(c.Ctx, time.Now().Add(time.Second*10))
	defer cancel()
	_, err := c.client.Ping(ctx, &pb.PingRequest{At: &timestamp.Timestamp{
		Seconds: t0.Unix(),
		Nanos:   0,
	}})
	if err != nil {
		return err
	}
	log.Println("Ping OK")
	return nil
}

func (c *Client) OpenProject(project *gs.Project) error {
	rep, err := c.client.OpenProject(c.Ctx, &pb.ProjectRequest{
		MissionId: project.ProjectID,
		Name:      project.Name,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to create or open project")
	}
	c.projectID = rep.MissionId
	project.ProjectID = rep.MissionId
	log.Infof("Server Msg: %s", rep.Message)
	return nil
}

func (c *Client) SendFile(ctx context.Context, filePath string) error {
	d, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	req := &pb.FileUpload{
		MissionId: c.projectID,
		FileName:  filePath,
		FileSize:  int64(len(d)),
		Data:      d,
	}
	_, err = c.client.SendFile(ctx, req, []grpc.CallOption{
		grpc.UseCompressor(gzip.Name),
	}...)
	if err != nil {
		log.Errorf("Failed to send file payload: %s", err)
		// TODO queue in retry buff
		log.Infof("TODO Add to queue: %v", req.FileName)
	}
	return nil
}

func (c *Client) Start() error {
	go monitorDirectory(c.Ctx, c.RootDir, c.newFiles)
	go func() {
		if err := c.processHotcluster(); err != nil {
			log.Errorf("hot cluster handler exited: %s", err)
		}
	}()
	go func() {
		if err := c.processPosition(); err != nil {
			log.Errorf("hot cluster handler exited: %s", err)
		}
	}()
	for {
		select {
		case newFile := <-c.newFiles:
			fileName := filepath.Base(newFile)
			if strings.HasPrefix(fileName, "stage_position_file") {
				t, err := tail.TailFile(newFile, tail.Config{
					Poll:      true,
					Follow:    true,
					MustExist: false,
					ReOpen:    true,
					Logger:    log.StandardLogger(),
				})
				if err != nil {
					return errors.Wrap(err, "Failed to open filePath, is the path correct?")
				}
				go watchFile(t, c.positions)
				c.tailedFiles = append(c.tailedFiles, t)
			} else if strings.HasPrefix(fileName, "hotcluster") {
				t, err := tail.TailFile(newFile, tail.Config{
					Poll:      true,
					Follow:    true,
					MustExist: false,
					ReOpen:    true,
					Logger:    log.StandardLogger(),
				})
				if err != nil {
					return errors.Wrap(err, "Failed to open filePath, is the path correct?")
				}
				go watchFile(t, c.hotSpots)
				c.tailedFiles = append(c.tailedFiles, t)
			} else if strings.Contains(fileName, ".png") {
				go func() {
					lastSize := int64(-1)
					t := time.NewTicker(time.Second * 2)
					for {
						select {
						case <-t.C:
							newSize, err := getSize(newFile)
							if err != nil {
								log.Errorf("Skipping file on error %s: %v", newFile, err)
								return
							}
							if lastSize == newSize {
								if err := c.SendFile(c.Ctx, newFile); err != nil {
									log.Errorf("Failed to send file: %v", err)
									return
								}
								log.Infof("Sent file successfully %d: %s", lastSize, newFile)
								return
							}
							lastSize = newSize
						}
					}

				}()
			}

		case <-c.done:
			for _, tf := range c.tailedFiles {
				if err := tf.Stop(); err != nil {
					log.Errorf("Failed to stop file watch: %s", tf.Filename)
				}
			}
			return nil
		}
	}
}

func getSize(fileName string) (int64, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	s, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return s.Size(), nil
}

func New(opts Opts) *Client {
	app := &Client{
		Opts:      opts,
		newFiles:  make(chan string),
		done:      make(chan bool),
		hotSpots:  make(chan TimedLine),
		positions: make(chan TimedLine),
		Ctx:       context.Background(),
	}
	return app
}

type gRPCOpts struct {
	Tls          bool
	ServerAddr   string
	CaFile       string
	HostOverride string
}

func newGRPCConn(opts gRPCOpts) (*grpc.ClientConn, error) {
	var clientOpts []grpc.DialOption
	if opts.Tls {
		clientCredentials, err := credentials.NewClientTLSFromFile(opts.CaFile, opts.HostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		clientOpts = append(clientOpts, grpc.WithTransportCredentials(clientCredentials))
	} else {
		clientOpts = append(clientOpts, grpc.WithInsecure())
	}
	clientOpts = append(clientOpts, grpc.WithBlock())
	conn, err := grpc.Dial(opts.ServerAddr, clientOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return conn, nil
}
