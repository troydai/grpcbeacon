# gRPC Beacon API Summary

## Public gRPC APIs

### BeaconService (`troydai.grpcbeacon.v1.BeaconService`)
- **Signal**: `Signal(SignalRequest) returns (SignalResponse)`
  - Accepts signal requests and returns timestamp with server details
  - Method: `/troydai.grpcbeacon.v1.BeaconService/Signal`

### Health Service (`grpc.health.v1.Health`)
- **Check**: `Check(HealthCheckRequest) returns (HealthCheckResponse)`
  - Supports services: `""`, `"liveness"`, `"readiness"`, `"Beacon"`

## Message Types

### SignalRequest
- `message` (string): Optional message to include in request

### SignalResponse
- `reply` (string): Server response with timestamp
- `details` (map<string,string>): Server details (hostname, beacon name)

## Configuration Functions

### Settings Package (`internal/settings`)
- **LoadEnvironment()**: `() -> (Environment, error)`
  - Loads environment variables (HOSTNAME)
- **LoadConfig()**: `() -> (Configuration, error)`
  - Loads TOML configuration file with defaults

### Configuration Types
- **Environment**: `{ HostName string }`
- **Configuration**: `{ Name, Address string; Port int; Logging *Logging; TLS *TLSConfiguration }`
- **Logging**: `{ Development bool }`
- **TLSConfiguration**: `{ Enabled bool; KeyFilePath, CertFilePath string }`

## Server Components

### Beacon Service (`internal/beacon`)
- **newService()**: `(hostName, beaconName string, logger *zap.Logger) -> *service`
- **ProvideRegister()**: `(Param) -> Result`
  - FX module provider for beacon service

### RPC Server (`internal/rpc`)
- **RegisterRPCServer()**: `(Param) -> error`
  - Configures and starts gRPC server
- **DetermineTLSOption()**: `(Configuration) -> (grpc.ServerOption, error)`
  - Configures TLS settings
- **GRPCRegisterFromFn()**: `(func(*grpc.Server) error) -> GRPCRegister`
  - Creates gRPC register from function

### Health Service (`internal/health`)
- **ProvideHealthCheckService()**: `() -> Result`
  - FX module provider for health checks

### Logging (`internal/logging`)
- **NewLogger()**: `(Configuration) -> (*zap.Logger, error)`
  - Creates production or development logger

## Client Interfaces

### Go Client
```go
type BeaconServiceClient interface {
    Signal(ctx context.Context, in *SignalRequest, opts ...grpc.CallOption) (*SignalResponse, error)
}
```

### Server Interface
```go
type BeaconServiceServer interface {
    Signal(context.Context, *SignalRequest) (*SignalResponse, error)
    mustEmbedUnimplementedBeaconServiceServer()
}
```

## Build Commands

- `make bin`: Build server binary
- `make run`: Run server with demo config
- `make test`: Run tests
- `make gen`: Generate protobuf code

## Command Line Usage

### Server
```bash
./server -config=/path/to/config.toml
```

### Client Examples
```bash
# Basic signal
grpcurl --plaintext localhost:8080 troydai.grpcbeacon.v1.BeaconService.Signal

# With message
grpcurl --plaintext -d '{"message": "test"}' localhost:8080 troydai.grpcbeacon.v1.BeaconService.Signal

# Health check
grpcurl --plaintext localhost:8080 grpc.health.v1.Health/Check

# TLS
grpcurl --cacert demo/certs/root.crt.pem localhost:8080 troydai.grpcbeacon.v1.BeaconService.Signal
```

## Docker Commands

```bash
# Build image
docker build -t grpcbeacon .

# Run container
docker run -p 8080:8080 grpcbeacon
```

## Error Types

- **Configuration Errors**: Environment parsing, TOML format, file permissions
- **TLS Errors**: Certificate loading, format validation, file access
- **Service Errors**: `codes.NotFound` for unknown health check services

## Default Configuration

```toml
name = "red cliff"
address = "127.0.0.1"
port = 8080
# logging.Development = false (production mode)
# tls.Enabled = false (TLS disabled)
```

## FX Modules

- **settings.Module**: Configuration loading
- **rpc.Module**: gRPC server setup
- **logging.Module**: Logger configuration
- **beacon.Module**: Beacon service
- **health.Module**: Health check service