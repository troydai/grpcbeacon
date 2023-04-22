package settings

import (
	"fmt"

	env "github.com/caarlos0/env/v6"
	"go.uber.org/fx"
)

var Module = fx.Provide(LoadEnvironment)

type Environment struct {
	HostName   string `env:"HOSTNAME"`
	BeaconName string `env:"BEACON_NAME"`
}

func LoadEnvironment() (Environment, error) {
	var e Environment
	if err := env.Parse(&e); err != nil {
		return e, fmt.Errorf("fail to parse environment variables: %w", err)
	}

	return e, nil
}
