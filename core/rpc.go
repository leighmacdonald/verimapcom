package core

import (
	"context"
	"github.com/leighmacdonald/verimapcom/pb"
	"github.com/leighmacdonald/verimapcom/store"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
		return nil, status.Errorf(codes.PermissionDenied, "Invalid auth")
	}
	s.sessionsMu.RLock()
	cp, found := s.sessions[tokens[0]]
	s.sessionsMu.RUnlock()
	if found {
		return cp, nil
	}
	if err := store.LoadPersonByToken(s.core.ctx, s.core.db, tokens[0], &p.Person); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Invalid auth")
	}

	s.sessionsMu.Lock()
	s.sessions[tokens[0]] = &p
	s.sessionsMu.Unlock()

	return &p, nil
}

type Opts struct {
	Tls      bool
	CertFile string
	KeyFile  string
	DataRoot string
}

func NewGRPCServer(core *Core, opts Opts) *grpc.Server {
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
