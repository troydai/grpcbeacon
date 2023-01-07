package beacon

import (
	"context"
	"fmt"
	"time"

	api "github.com/troydai/grpcbeacon/api/protos"
)

type Server struct {
	api.UnimplementedBeaconServer
}

func (s *Server) Signal(_ context.Context, req *api.SignalReqeust) (*api.SignalResponse, error) {
	resp := &api.SignalResponse{
		Reply: fmt.Sprintf("Beacon signal at %s", time.Now().Format(time.RFC1123)),
	}

	return resp, nil
}
