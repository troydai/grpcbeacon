package settings_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/troydai/grpcbeacon/internal/settings"
)

const _testSample1 = `
name = "white peak"
address = "127.0.0.1"
port = 6899

[logging]
Development = true

[tls]
Enabled = true
KeyFilePath = "/path/to/key"
CertFilePath = "/path/to/cert"
`

const _testSample2 = `
name = "white peak"
address = "127.0.0.1"
port = 6899
`

func TestDataModel(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expectation func(*testing.T, settings.Configuration)
	}{
		{
			name:  "full configuration",
			input: _testSample1,
			expectation: func(t *testing.T, c settings.Configuration) {

				require.NotNil(t, c.Logging)
				assert.True(t, c.Logging.Development)

				require.NotNil(t, c.TLS)
				assert.True(t, c.TLS.Enabled)
				assert.Equal(t, "/path/to/key", c.TLS.KeyFilePath)
				assert.Equal(t, "/path/to/cert", c.TLS.CertFilePath)
			},
		},
		{
			name:  "minimal configuration",
			input: _testSample2,
			expectation: func(t *testing.T, c settings.Configuration) {
				assert.Nil(t, c.Logging)
				assert.Nil(t, c.TLS)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var c settings.Configuration
			_, err := toml.Decode(tc.input, &c)
			require.NoError(t, err)

			tc.expectation(t, c)
		})
	}
}
