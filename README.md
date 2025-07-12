# gRPC Beacon

A grpc service that respond to the request. Utility for demo purposes

## Prerequisites

- **Go 1.21+**: [Install Go](https://golang.org/doc/install)
- **buf**: Protocol buffer compiler (auto-installed by setup script)

## Quick Start

1. **Setup development environment:**
   ```bash
   make setup
   ```

2. **Generate protocol buffer code:**
   ```bash
   make gen
   ```

3. **Build the server:**
   ```bash
   make bin
   ```

## Run locally

In the case when the self signed cert is expired, use the following command
to create new cert.

```bash
./scripts/update_certs.sh ./demo/certs
```

Start server

```bash
make run
```

Query

```bash
grpcurl --cacert ./demo/certs/root.crt.pem localhost:8089 troydai.grpcbeacon.v1.BeaconService.Signal
```

## Available Make Targets

- `make setup` - Setup development environment and install required tools
- `make gen` - Generate Go code from protocol buffer definitions
- `make bin` - Build the server binary
- `make run` - Start the server locally
- `make test` - Run tests
- `make check-tools` - Verify required tools are installed

## Troubleshooting

### `make gen` fails with "buf: command not found"

Run the setup script to install required tools:
```bash
make setup
```

### Architecture issues during build

The Makefile automatically handles architecture translation (x86_64 → amd64). If you encounter issues, verify your Go installation:
```bash
go version
go env GOOS GOARCH
```

## References

- Image registry: https://hub.docker.com/repository/docker/troydai/grpcbeacon
