package beacon

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/troydai/grpcbeacon/gen/go/troydai/grpcbeacon/v1"
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
			pb.RegisterBeaconServiceServer(s, svc)
			return nil
		}),
	}
}
