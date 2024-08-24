package beacon

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	pb "github.com/troydai/grpcbeacon/gen/api/protos/beacon"
)

type service struct {
	pb.UnimplementedBeaconServiceServer

	details map[string]string
	logger  *zap.Logger
}

var _ pb.BeaconServiceServer = (*service)(nil)

func newService(hostName, beaconName string, logger *zap.Logger) *service {
	s := &service{
		details: make(map[string]string),
		logger:  logger,
	}

	s.details["Hostname"] = hostName
	s.details["BeaconName"] = beaconName

	return s
}

func (s *service) Signal(ctx context.Context, req *pb.SignalRequest) (*pb.SignalResponse, error) {
	logger := s.logger
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		logger = logger.With(zap.Any("metadata", md))
	}

	logger.Info("Signal received")
	resp := &pb.SignalResponse{
		Reply: fmt.Sprintf("Beacon signal at %s", time.Now().Format(time.RFC1123)),
	}

	if len(s.details) > 0 {
		resp.Details = s.details
	}

	return resp, nil
}
