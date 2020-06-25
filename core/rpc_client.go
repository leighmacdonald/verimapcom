package core

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/leighmacdonald/verimapcom/pb"
	"github.com/leighmacdonald/verimapcom/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *RPCServer) ClientStreamMissionEvents(req *pb.MissionRequest, srv pb.RPC_ClientStreamMissionEventsServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStreamMissionEvents not implemented")
}
func (s *RPCServer) ClientStreamPositions(req *pb.MissionRequest, srv pb.RPC_ClientStreamPositionsServer) error {
	ps, err := store.FlightPositionsSince(s.core.ctx, s.core.db, req.MissionId, req.StartIdx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to lookup positions")
	}
	for _, p := range ps {
		if err := srv.Send(&pb.PositionEvent{
			At: &timestamp.Timestamp{
				Seconds: p.CreatedOn.Unix(),
				Nanos:   0,
			},
			Id: p.ID,
			Location: &pb.Location{
				Lat: p.Lat,
				Lon: p.Lon,
			},
			Elevation: p.Elevation,
			FlightId:  p.FlightID,
		}); err != nil {
			return status.Errorf(codes.Internal, "Failed to send position event: %v", err)
		}
	}
	// TODO watch a channel for new events
	return nil
}
func (s *RPCServer) ClientStreamHotSpots(req *pb.MissionRequest, srv pb.RPC_ClientStreamHotSpotsServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStreamHotSpots not implemented")
}

func (s *RPCServer) ClientSendMessage(ctx context.Context, req *pb.ChatMessageRequest) (*pb.StatusReply, error) {
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
