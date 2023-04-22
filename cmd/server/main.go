package main

import (
	"go.uber.org/fx"

	"github.com/troydai/grpcbeacon/internal/beacon"
	"github.com/troydai/grpcbeacon/internal/logging"
	"github.com/troydai/grpcbeacon/internal/rpc"
	"github.com/troydai/grpcbeacon/internal/settings"
)

func main() {
	app := fx.New(
		settings.Module,
		beacon.Module,
		logging.Module,
		rpc.Module,
	)

	app.Run()
}
