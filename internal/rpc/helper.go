package rpc

import "google.golang.org/grpc"

type genericGRPCRegister struct {
	fn func(*grpc.Server) error
}

var _ GRPCRegister = (*genericGRPCRegister)(nil)

func (r *genericGRPCRegister) Register(s *grpc.Server) error {
	return r.fn(s)
}
