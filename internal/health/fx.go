package health

import (
	"go.uber.org/fx"

	"github.com/troydai/grpcbeacon/internal/rpc"
)

var Module = fx.Options(fx.Provide(ProvideHealthCheckService))

type Result struct {
	fx.Out

	Register rpc.GRPCRegister `group:"grpc_registers"`
}

func ProvideHealthCheckService() Result {
	return Result{Register: &healthcheck{}}
}
