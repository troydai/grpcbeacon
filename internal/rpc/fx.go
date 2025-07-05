package rpc

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/troydai/grpcbeacon/internal/settings"
)

var Module = fx.Invoke(RegisterRPCServer)

type (
	Param struct {
		fx.In

		Lifecycle     fx.Lifecycle
		Logger        *zap.Logger
		GRPCRegisters []GRPCRegister `group:"grpc_registers"`
		Config        settings.Configuration
	}

	GRPCRegister interface {
		Register(*grpc.Server) error
	}
)

func GRPCRegisterFromFn(fn func(*grpc.Server) error) GRPCRegister {
	return &genericGRPCRegister{fn: fn}
}

func RegisterRPCServer(param Param) error {
	if len(param.GRPCRegisters) == 0 {
		return fmt.Errorf("no grpc register found")
	}

	var serverOptions []grpc.ServerOption

	tlsOpt, err := DetermineTLSOption(param.Config)
	if err != nil {
		return fmt.Errorf("fail to determine TLS option: %w", err)
	}
	if tlsOpt != nil {
		serverOptions = append(serverOptions, tlsOpt)
	}

	s := grpc.NewServer(serverOptions...)
	reflection.Register(s)
	for _, r := range param.GRPCRegisters {
		if err := r.Register(s); err != nil {
			return fmt.Errorf("fail to register grpc server: %w", err)
		}
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", param.Config.Address, param.Config.Port))
	if err != nil {
		return fmt.Errorf("fail to start TCP listener: %w", err)
	}

	serverStopped := make(chan struct{})
	param.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				defer close(serverStopped)
				if err := s.Serve(lis); err != nil {
					param.Logger.Error("gRPC server failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			go func() {
				s.GracefulStop()
			}()

			select {
			case <-ctx.Done():
				return fmt.Errorf("fail to stop daemon in time: %w", ctx.Err())
			case <-serverStopped:
				return nil
			}
		},
	})

	return nil
}
