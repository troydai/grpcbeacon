package echoserver

import (
	"context"

	api "github.com/troydai/grpcecho/protos"
)

type Server struct {
	api.UnimplementedServiceServer
}

func (s *Server) Echo(_ context.Context, req *api.EchoReqeust) (*api.EchoResponse, error) {
	resp := &api.EchoResponse{
		Reply: req.Greet,
	}

	return resp, nil
}
