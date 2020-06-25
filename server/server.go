package server

import (
	"context"
	"fmt"
	"github.com/leighmacdonald/verimapcom/gs"
	"github.com/leighmacdonald/verimapcom/pb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type server struct {
	pb.UnimplementedRPCServer
	projectDir string
	projects   map[int32]*gs.Project
	projectsMu *sync.RWMutex
}

func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	return &pb.PingReply{Ok: true}, nil
}

func (s *server) OpenProject(ctx context.Context, in *pb.ProjectRequest) (*pb.ProjectReply, error) {
	project, err := gs.OpenProject(s.projectDir, in.MissionId)
	if err != nil {
		return nil, err
	}
	s.projectsMu.Lock()
	s.projects[project.ProjectID] = project
	s.projectsMu.Unlock()
	var m string
	if in.MissionId == 0 {
		m = fmt.Sprintf("Created new project: %d", project.ProjectID)
	} else {
		m = fmt.Sprintf("Opened existing project: %d", project.ProjectID)
	}
	return &pb.ProjectReply{
		MissionId: project.ProjectID,
		Message:   m,
	}, nil
}

func (s *server) SendPosition(in pb.RPC_SendPositionServer) error {
	for {
		resp, err := in.Recv()
		if err == io.EOF {
			log.Infof("position received EOF")
			break
		}
		if err != nil {
			log.Errorf("Could not receive position: %v", err)
			break
		}
		pos := gs.Position{
			At:        time.Unix(resp.At.Seconds, int64(resp.At.Nanos/1000000)),
			Lat:       resp.Location.Lat,
			Lon:       resp.Location.Lon,
			Elevation: float64(resp.Elevation),
		}
		s.projectsMu.Lock()
		s.projects[resp.MissionId].AddPosition(pos)
		s.projectsMu.Unlock()
	}
	return nil
}

func (s *server) SendHotspot(in pb.RPC_SendHotspotServer) error {
	for {
		resp, err := in.Recv()
		if err == io.EOF {
			log.Infof("hotspot received EOF")
			break
		}
		if err != nil {
			log.Errorf("Could not receive hotspot: %v", err)
			break
		}
		hs := gs.HotSpot{
			ID:    resp.Id,
			Lat:   resp.Location.Lat,
			Lon:   resp.Location.Lon,
			Delta: float64(resp.Delta),
		}
		if resp.MissionId == 0 {
			return errors.New("Invalid mission id")
		}
		s.projectsMu.RLock()
		p, found := s.projects[resp.MissionId]
		s.projectsMu.RUnlock()
		if found {
			p.AddHotspot(hs)
		}

	}
	return nil
}

func (s *server) SendFile(ctx context.Context, in *pb.FileUpload) (*pb.FileReply, error) {
	s.projectsMu.RLock()
	_, ok := s.projects[in.MissionId]
	s.projectsMu.RUnlock()
	if !ok {
		return &pb.FileReply{
			Status: 2,
		}, errors.Errorf("Invalid project: %d", in.MissionId)
	}

	op := path.Join(
		s.projectDir,
		fmt.Sprintf("project_%d", in.MissionId),
		"ir_export",
		path.Base(in.FileName))
	if err := ioutil.WriteFile(op, in.Data, 0766); err != nil {
		return &pb.FileReply{
			Status: 2,
		}, errors.Errorf("Error writing file: %v", err)
	}
	log.Infof("Wrote client file: %s", op)
	return &pb.FileReply{
		Status: 1,
	}, nil

}

type Opts struct {
	Tls      bool
	CertFile string
	KeyFile  string
	DataRoot string
}

func matchProject(path string) int32 {
	rx := regexp.MustCompile(`^project_(\d+)$`)
	m := rx.FindStringSubmatch(path)
	if len(m) > 0 {
		v, err := strconv.ParseInt(m[1], 10, 32)
		if err != nil {
			return 0
		}
		return int32(v)
	}
	return 0
}

func NewServer(opts Opts) *grpc.Server {
	var serverOpts []grpc.ServerOption
	if opts.Tls {
		creds, err := credentials.NewServerTLSFromFile(opts.CertFile, opts.KeyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		serverOpts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	if opts.DataRoot == "" {
		opts.DataRoot = "./"
	}
	projects := make(map[int32]*gs.Project)
	if err := filepath.Walk(opts.DataRoot, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			projectID := matchProject(info.Name())
			if projectID == 0 {
				return nil
			}
			np, err := gs.NewProject(projectID)
			if err != nil {
				return errors.Wrap(err, "Failed to open project")
			}
			projects[np.ProjectID] = np
		}
		return nil
	}); err != nil {
		log.Errorf("Error reading project data root: %v", err)
	}
	grpcServer := grpc.NewServer(serverOpts...)
	pb.RegisterRPCServer(grpcServer, &server{
		projects:   projects,
		projectsMu: &sync.RWMutex{},
		projectDir: opts.DataRoot,
	})
	return grpcServer
}
