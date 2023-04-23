package main

import (
	"go.uber.org/fx"

	"github.com/troydai/grpcbeacon/internal/beacon"
	"github.com/troydai/grpcbeacon/internal/health"
	"github.com/troydai/grpcbeacon/internal/logging"
	"github.com/troydai/grpcbeacon/internal/rpc"
	"github.com/troydai/grpcbeacon/internal/settings"
)

func main() {
	basic := fx.Options(
		settings.Module,
		rpc.Module,
		logging.Module,
	)

	services := fx.Options(
		beacon.Module,
		health.Module,
	)

	fx.New(basic, services).Run()
}
