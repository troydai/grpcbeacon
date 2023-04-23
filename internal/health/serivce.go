package health

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	healthapi "github.com/troydai/grpcbeacon/gen/api/protos/grpchealth"
	"github.com/troydai/grpcbeacon/internal/rpc"
)

var (
	_healthResp = &healthapi.HealthCheckResponse{Status: healthapi.HealthCheckResponse_SERVING}
)

type healthcheck struct {
	healthapi.UnimplementedHealthServer
}

var _ healthapi.HealthServer = (*healthcheck)(nil)
var _ rpc.GRPCRegister = (*healthcheck)(nil)

func (s *healthcheck) Check(_ context.Context, req *healthapi.HealthCheckRequest) (*healthapi.HealthCheckResponse, error) {
	switch req.Service {
	case "":
		return _healthResp, nil
	case "Beacon":
		return _healthResp, nil
	}

	return nil, grpc.Errorf(codes.NotFound, "service %s not found", req.Service)
}

func (s *healthcheck) Register(server *grpc.Server) error {
	if s == nil {
		return fmt.Errorf("health check service is nil")
	}

	healthapi.RegisterHealthServer(server, s)
	return nil
}
