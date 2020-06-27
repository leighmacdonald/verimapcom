package core

import (
	"context"
	"fmt"
	"github.com/leighmacdonald/verimapcom/pb"
	"github.com/leighmacdonald/verimapcom/store"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (s *RPCServer) CreateMission(ctx context.Context, in *pb.CreateMissionRequest) (*pb.MissionReply, error) {
	p, err := s.personFromCtx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Invalid auth token: %v", err)
	}
	m := store.NewMission(p.Person)
	m.MissionName = in.Name
	if m.MissionName == "" {
		m.MissionName = fmt.Sprintf("Mission-%s", time.Now().Format(time.RFC822Z))
	}
	if err := store.SaveMission(s.core.ctx, s.core.db, &m); err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "Failed to save mission: %v", err)
	}
	if p.MissionID > 0 {
		log.Fatalf("XXX Why?")
	}
	p.Mu.Lock()
	p.MissionID = m.MissionID
	p.Mu.Unlock()
	return &pb.MissionReply{
		MissionId: m.MissionID,
		Message:   "Created successfully",
		Name:      m.MissionName,
	}, nil
}

func (s *RPCServer) OpenMission(ctx context.Context, in *pb.MissionRequest) (*pb.MissionReply, error) {
	p, err := s.personFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	m, err2 := s.core.Mission(in.MissionId)
	if err2 != nil {
		return nil, status.Errorf(codes.NotFound, "Failed to open mission")
	}
	p.Mu.Lock()
	p.MissionID = m.MissionID
	p.Mu.Unlock()
	return &pb.MissionReply{
		MissionId: m.MissionID,
		Message:   "Opened successfully",
	}, nil
}

func (s *RPCServer) CreateFlight(ctx context.Context, req *pb.CreateFlightRequest) (*pb.CreateFlightResponse, error) {
	p, err := s.personFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if p.MissionID <= 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "Invalid mission id")
	}
	flight := store.Flight{
		MissionID:   p.MissionID,
		FlightState: 1,
		Summary:     req.Description,
		CreatedOn:   time.Now(),
	}
	if err := store.FlightSave(s.core.ctx, s.core.db, &flight); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save flight: %v", err)
	}
	return &pb.CreateFlightResponse{
		FlightId: flight.FlightID,
	}, nil
}

func (s *RPCServer) SendMessage(ctx context.Context, req *pb.ChatMessageRequest) (*pb.StatusReply, error) {
	p, err := s.personFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if p.MissionID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Mission ID must be valid")
	}
	e := store.NewMissionEvent(store.EvtMessage, p.MissionID)
	e.Payload = map[string]interface{}{
		"message":   req.Message,
		"person_id": p.PersonID,
	}
	if err := store.MissionEventAdd(s.core.ctx, s.core.db, &e); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add event to db: %v", err)
	}
	return &pb.StatusReply{Status: 0}, nil
}
