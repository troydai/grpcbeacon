package logging

import (
	"github.com/troydai/grpcbeacon/internal/settings"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(NewLogger),
	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger}
	}),
)

func NewLogger(c settings.Configuration) (*zap.Logger, error) {
	if c.Logging != nil && c.Logging.Development {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
}
