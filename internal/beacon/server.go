package beacon

import (
	"context"
	"fmt"
	"time"

	api "github.com/troydai/grpcbeacon/gen/api/protos"
	"github.com/troydai/grpcbeacon/internal/settings"
	"go.uber.org/zap"
)

type Server struct {
	api.UnimplementedBeaconServer

	details map[string]string
	logger  *zap.Logger
}

var _ api.BeaconServer = (*Server)(nil)

func NewServer(env settings.Environment, logger *zap.Logger) *Server {
	s := &Server{
		details: make(map[string]string),
		logger:  logger,
	}

	s.details["Hostname"] = env.HostName
	s.details["Flockname"] = env.FlockName

	return s
}

func (s *Server) Signal(_ context.Context, req *api.SignalReqeust) (*api.SignalResponse, error) {
	s.logger.Info("Signal received")

	resp := &api.SignalResponse{
		Reply: fmt.Sprintf("Beacon signal at %s", time.Now().Format(time.RFC1123)),
	}

	if len(s.details) > 0 {
		resp.Details = s.details
	}

	return resp, nil
}
