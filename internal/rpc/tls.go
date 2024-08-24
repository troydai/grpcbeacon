package rpc

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/troydai/grpcbeacon/internal/settings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/* facilitate the TLS */

func DetermineTLSOption(cfg settings.Configuration) (grpc.ServerOption, error) {
	if cfg.TLS == nil || !cfg.TLS.Enabled {
		return nil, nil
	}

	keyFilepath, err := resolveFilePath(cfg.TLS.KeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("fail to resolve key file path: %w", err)
	}

	certFilePath, err := resolveFilePath(cfg.TLS.CertFilePath)
	if err != nil {
		return nil, fmt.Errorf("fail to resolve cert file path: %w", err)
	}

	cred, err := credentials.NewServerTLSFromFile(certFilePath, keyFilepath)
	if err != nil {
		return nil, fmt.Errorf("fail to create credentials: %w", err)
	}

	return grpc.Creds(cred), nil
}

func resolveFilePath(filepath string) (string, error) {
	if filepath == "" {
		return "", errors.New("path is empty")
	}

	if !path.IsAbs(filepath) {
		filepath = path.Join(os.Getenv("PWD"), filepath)
	}

	stat, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", filepath)
	}
	if err != nil {
		return "", fmt.Errorf("fail to stat file: %w", err)
	}
	if stat.IsDir() {
		return "", fmt.Errorf("file is a directory: %s", filepath)
	}

	return filepath, nil
}
