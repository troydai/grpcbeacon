package beacon

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	beaconapi "github.com/troydai/grpcbeacon/gen/api/protos/beacon"
)

type service struct {
	beaconapi.UnimplementedBeaconServer

	details map[string]string
	logger  *zap.Logger
}

var _ beaconapi.BeaconServer = (*service)(nil)

func newService(hostName, beaconName string, logger *zap.Logger) *service {
	s := &service{
		details: make(map[string]string),
		logger:  logger,
	}

	s.details["Hostname"] = hostName
	s.details["BeaconName"] = beaconName

	return s
}

func (s *service) Signal(ctx context.Context, req *beaconapi.SignalRequest) (*beaconapi.SignalResponse, error) {
	logger := s.logger
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		logger = logger.With(zap.Any("metadata", md))
	}

	logger.Info("Signal received")
	resp := &beaconapi.SignalResponse{
		Reply: fmt.Sprintf("Beacon signal at %s", time.Now().Format(time.RFC1123)),
	}

	if len(s.details) > 0 {
		resp.Details = s.details
	}

	return resp, nil
}
