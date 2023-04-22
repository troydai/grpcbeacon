package settings

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	env "github.com/caarlos0/env/v6"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(LoadEnvironment, LoadConfig)

const _defaultConfigPath = "/etc/beacon-svc/beacon.toml"

type (
	Environment struct {
		HostName string `env:"HOSTNAME"`
	}

	Configuration struct {
		Name    string
		Address string
		Port    int
	}
)

func LoadEnvironment() (Environment, error) {
	var e Environment
	if err := env.Parse(&e); err != nil {
		return e, fmt.Errorf("fail to parse environment variables: %w", err)
	}

	return e, nil
}

func LoadConfig(logger *zap.Logger) (Configuration, error) {
	var configFlag string
	flag.StringVar(&configFlag, "config", "", "path to config file")
	flag.Parse()

	if configFlag == "" {
		configFlag = _defaultConfigPath
	}

	logger.Info(
		"attempt to read config file",
		zap.String("path", configFlag),
		zap.Any("os.Args", os.Args),
	)

	if _, err := os.Stat(configFlag); err != nil {
		if os.IsNotExist(err) {
			logger.Info("fall back to default config")
			return defaultConfig(), nil
		}
		return Configuration{}, fmt.Errorf("fail to stat config file: %w", err)
	}

	var config Configuration
	if _, err := toml.DecodeFile(configFlag, &config); err != nil {
		return Configuration{}, fmt.Errorf("fail to decode config file: %w", err)
	}

	return config, nil
}

func defaultConfig() Configuration {
	return Configuration{
		Name:    "red cliff",
		Address: "127.0.0.1",
		Port:    8080,
	}
}
