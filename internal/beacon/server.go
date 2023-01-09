package beacon

import (
	"context"
	"fmt"
	"os"
	"time"

	api "github.com/troydai/grpcbeacon/api/protos"
)

type Server struct {
	api.UnimplementedBeaconServer

	details map[string]string
}

var _ api.BeaconServer = (*Server)(nil)

func NewServer() *Server {
	s := &Server{
		details: make(map[string]string),
	}

	s.details["Hostname"] = os.Getenv("HOSTNAME")

	return s
}

func (s *Server) Signal(_ context.Context, req *api.SignalReqeust) (*api.SignalResponse, error) {
	resp := &api.SignalResponse{
		Reply: fmt.Sprintf("Beacon signal at %s", time.Now().Format(time.RFC1123)),
	}

	if len(s.details) > 0 {
		resp.Details = s.details
	}

	return resp, nil
}
