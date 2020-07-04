package core

import (
	"context"
	"github.com/leighmacdonald/verimapcom/pb"
	"github.com/leighmacdonald/verimapcom/store"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"sync"
)

type RPCServer struct {
	pb.UnimplementedRPCServer
	core       *Core
	sessions   map[string]*store.ContextualPerson
	sessionsMu *sync.RWMutex
}

func (s *RPCServer) personFromCtx(ctx context.Context) (*store.ContextualPerson, error) {
	var p store.ContextualPerson
	p.Mu = &sync.RWMutex{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Invalid auth")
	}
	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		log.Errorf("Missing auth")
		//return nil, status.Errorf(codes.PermissionDenied, "Invalid auth")
	}
	s.sessionsMu.RLock()
	cp, found := s.sessions["xxxx"]
	s.sessionsMu.RUnlock()
	if found {
		return cp, nil
	}
	if err := store.LoadPersonByToken(s.core.ctx, s.core.db, "xxxx", &p.Person); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Invalid auth")
	}

	s.sessionsMu.Lock()
	s.sessions["xxxx"] = &p
	s.sessionsMu.Unlock()

	return &p, nil
}

type Opts struct {
	Tls      bool
	CertFile string
	KeyFile  string
	DataRoot string
}

type TokenCredential string

func (c TokenCredential) Info() credentials.ProtocolInfo {
	panic("implement me")
}

func (c TokenCredential) Clone() credentials.TransportCredentials {
	panic("implement me")
}

func (c TokenCredential) OverrideServerName(s string) error {
	panic("implement me")
}

func (c TokenCredential) ClientHandshake(ctx context.Context, s string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	panic("implement me")
}

func (c TokenCredential) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	panic("implement me")
}

func (c TokenCredential) GetRequestMetadata(ctx context.Context) (map[string]string, error) {
	return map[string]string{
		"authorization": "xxxx",
	}, nil
}

func NewGRPCServer(core *Core, opts Opts) *grpc.Server {
	var serverOpts []grpc.ServerOption
	//serverOpts = append(serverOpts, grpc.WithPerRPCCredentials())

	tlsEnabled := viper.GetBool("tls_enabled")
	if tlsEnabled && opts.Tls {
		log.Fatalf("TLS NOT SUPPORTED")
		creds, err := credentials.NewServerTLSFromFile(opts.CertFile, opts.KeyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		serverOpts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	if opts.DataRoot == "" {
		opts.DataRoot = "./"
	}

	grpcServer := grpc.NewServer(serverOpts...)
	pb.RegisterRPCServer(grpcServer, &RPCServer{
		core:       core,
		sessions:   make(map[string]*store.ContextualPerson),
		sessionsMu: &sync.RWMutex{},
	})
	return grpcServer
}

//func tsToTime(ts timestamp.Timestamp) time.Time {
//	return time.Unix(ts.Seconds, int64(ts.Nanos))
//}
