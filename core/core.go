package core

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/leighmacdonald/verimapcom/store"
	"github.com/leighmacdonald/verimapcom/web"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	ErrInvalidMission = errors.New("invalid mission")
	ErrDuplicate      = errors.New("duplicate entry")
)

type Core struct {
	uploadDir  string
	ctx        context.Context
	db         *pgxpool.Pool
	missions   map[int32]*store.Mission
	missionsMu *sync.RWMutex
	web        *web.Web
	http       *http.Server
	grpc       *grpc.Server
}

func (c *Core) Mission(missionID int32) (*store.Mission, error) {
	c.missionsMu.RLock()
	m, found := c.missions[missionID]
	c.missionsMu.RUnlock()
	if !found {
		return nil, ErrInvalidMission
	}
	return m, nil
}

func (c *Core) MissionCreate(p store.Person) (*store.Mission, error) {
	var m store.Mission
	m.PersonID = p.PersonID
	m.AgencyID = p.AgencyID
	m.MissionName = "Unnamed Mission"
	m.MissionState = 1
	if err := store.SaveMission(c.ctx, c.db, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (c *Core) setupDB() error {
	dsn := viper.GetString("dsn")
	if viper.GetBool("migrate") {
		if err := store.Migrate(dsn); err != nil {
			if err.Error() != "no change" {
				log.Fatalf("Could not do migrations: %v", err)
			}
			log.Infof("No migration performed")
		}
	}
	c.db = store.MustConnectDB(c.ctx)
	return nil
}

func (c *Core) setupHTTP() error {
	w := web.New(c.ctx, viper.GetString("redis"))
	if err := w.Setup(); err != nil {
		log.Fatalf("Could not run setup: %v", err)
	}
	opts := web.DefaultHTTPOpts()
	opts.Handler = w.Handler
	c.http = web.NewHTTPServer(opts)
	return nil
}

func (c *Core) setupGRPC() error {
	s := NewGRPCServer(c, Opts{
		Tls: false,
	})
	c.grpc = s
	return nil
}

func New(ctx context.Context) (*Core, error) {
	c := Core{
		uploadDir:  "./uploads",
		ctx:        ctx,
		missions:   make(map[int32]*store.Mission),
		missionsMu: &sync.RWMutex{},
	}
	if err := c.setupDB(); err != nil {
		return nil, errors.Wrapf(err, "Failed to setup database")
	}
	if err := c.setupHTTP(); err != nil {
		return nil, errors.Wrapf(err, "Failed to setup HTTP service")
	}
	if err := c.setupGRPC(); err != nil {
		return nil, errors.Wrapf(err, "failed to setup gRPC service")
	}

	return &c, nil
}

func (c *Core) listenHTTP() error {
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Infof("HTTP listening on %s", viper.GetString("listen_http"))
	if err := c.http.ListenAndServe(); err != nil {
		log.Errorf("gRPC Shutdown unclean: %v", err)
	}
	return nil
}

func (c *Core) listenGRPC() error {
	listenAddr := viper.GetString("listen_grpc")
	lis, err3 := net.Listen("tcp", listenAddr)
	if err3 != nil {
		log.Fatalf("gRPC failed to listen: %v", err3)
	}
	log.Infof("gRPC listening on %s", listenAddr)
	if err := c.grpc.Serve(lis); err != nil {
		log.Errorf("Failed to serve: %s", err)
	}
	return nil
}

func (c *Core) AddPosition(flightID int32, t time.Time, pos store.PositionZ) error {
	if err := store.FlightPositionInsert(c.ctx, c.db, flightID, t, pos.Lat, pos.Lon, pos.Elevation); err != nil {
		return errors.Wrapf(err, "Could not add position: %v", err)
	}
	return nil
}

func (c *Core) AddHotspot(flightID int32, t time.Time, hs *store.HotSpot) error {
	if err := store.FlightPositionInsert(c.ctx, c.db, flightID, t, hs.Lat, hs.Lon, hs.Delta); err != nil {
		return errors.Wrapf(err, "Could not add hotspot: %v", err)
	}
	return nil
}

func (c *Core) SendMessage(m *store.Message) error {
	if err := store.MessageAdd(c.ctx, c.db, m); err != nil {
		return err
	}
	return nil
}

func (c *Core) ListenAndServe() error {
	go func() {
		if err := c.listenHTTP(); err != nil {
			log.Errorf("%v", err)
		}
	}()
	if err := c.listenGRPC(); err != nil {
		log.Errorf("%v", err)
	}
	return nil
}

func (c *Core) Close() {
	c.web.Close()
}
