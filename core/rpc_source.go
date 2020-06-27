package core

import (
	"context"
	"fmt"
	"github.com/leighmacdonald/verimapcom/pb"
	"github.com/leighmacdonald/verimapcom/store"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"path"
	"time"
)

func (s *RPCServer) SourceSendFile(ctx context.Context, in *pb.FileUpload) (*pb.FileReply, error) {
	p, err := s.personFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if p.MissionID <= 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "Must open mission")
	}
	missionDir := fmt.Sprintf("project_%d", p.MissionID)
	op := path.Join(s.core.uploadDir, missionDir, "ir_export", path.Base(in.FileName))
	if err := ioutil.WriteFile(op, in.Data, 0766); err != nil {
		return &pb.FileReply{
			Status: 2,
		}, status.Errorf(codes.Internal, "Error writing file: %v", err)
	}
	log.Infof("Wrote client file: %s", op)
	return &pb.FileReply{
		Status: 1,
	}, nil
}

func (s *RPCServer) SourceInsertPositions(srv pb.RPC_SourceInsertPositionsServer) error {
	for {
		resp, err := srv.Recv()
		if err == io.EOF {
			log.Infof("position received EOF")
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "Could not receive position: %v", err)
		}
		t := time.Unix(resp.At.Seconds, int64(resp.At.Nanos))
		pos := store.PositionZ{
			Position: store.Position{
				Lat: resp.Location.Lat,
				Lon: resp.Location.Lon,
			},
			Elevation: resp.Elevation,
		}
		if err := s.core.AddPosition(resp.FlightId, t, pos); err != nil {
			return status.Errorf(codes.Internal, "Failed to insert position: %v", err)
		}
	}
	return nil
}

func (s *RPCServer) SourceInsertHotspots(srv pb.RPC_SourceInsertHotspotsServer) error {
	for {
		resp, err := srv.Recv()
		if err == io.EOF {
			log.Infof("hotspot received EOF")
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "Could not receive hotspot: %v", err)
		}
		if resp.FlightId == 0 {
			return status.Errorf(codes.FailedPrecondition, "Invalid flight id")
		}
		hs := store.HotSpot{
			Position: store.Position{
				Lat: resp.Location.Lat,
				Lon: resp.Location.Lon,
			},
			ID:    resp.Id,
			Delta: resp.Delta,
		}
		log.Print(hs)
	}
	return nil
}
