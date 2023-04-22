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

type (
	Param struct {
		fx.In

		Env    settings.Environment
		Config settings.Configuration
		Logger *zap.Logger
	}

	Result struct {
		fx.Out

		Register rpc.GRPCRegister `group:"grpc_registers"`
	}
)

func ProvideRegister(param Param) Result {
	hostName := param.Env.HostName
	beaconName := param.Config.Name

	svc := newService(hostName, beaconName, param.Logger)

	return Result{
		Register: rpc.GRPCRegisterFromFn(func(s *grpc.Server) error {
			api.RegisterBeaconServer(s, svc)
			return nil
		}),
	}
}
