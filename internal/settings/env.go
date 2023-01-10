package settings

import (
	"fmt"

	env "github.com/caarlos0/env/v6"
)

type Environment struct {
	HostName  string `env:"HOSTNAME"`
	FlockName string `env:"FLOCKNAME"`
}

func LoadEnvironment() (Environment, error) {
	var e Environment
	if err := env.Parse(&e); err != nil {
		return e, fmt.Errorf("fail to parse environment variables: %w", err)
	}

	return e, nil
}
