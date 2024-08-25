package logging_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/troydai/grpcbeacon/internal/logging"
	"github.com/troydai/grpcbeacon/internal/settings"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		l, err := logging.NewLogger(settings.Configuration{})
		assert.NoError(t, err)
		assert.NotNil(t, l)

		assert.False(t, l.Core().Enabled(zapcore.DebugLevel))
		assert.True(t, l.Core().Enabled(zapcore.InfoLevel))
	})

	t.Run("debug level", func(t *testing.T) {
		l, err := logging.NewLogger(settings.Configuration{
			Logging: &settings.Logging{Development: true},
		})
		assert.NoError(t, err)
		assert.NotNil(t, l)

		assert.True(t, l.Core().Enabled(zapcore.DebugLevel))
		assert.True(t, l.Core().Enabled(zapcore.InfoLevel))
	})
}
