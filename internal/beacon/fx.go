package beacon

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	api "github.com/troydai/grpcbeacon/gen/api/protos"
	"github.com/troydai/grpcbeacon/internal/rpc"
	"github.com/troydai/grpcbeacon/internal/settings"
)

var Module = fx.Provide(ProvideRegister)

type Result struct {
	fx.Out

	Register rpc.GRPCRegister `group:"grpc_registers"`
}

func ProvideRegister(env settings.Environment, logger *zap.Logger) Result {
	svc := newService(env, logger)

	return Result{
		Register: rpc.GRPCRegisterFromFn(func(s *grpc.Server) error {
			api.RegisterBeaconServer(s, svc)
			return nil
		}),
	}
}
