# gRPC Beacon API Documentation

## Overview

gRPC Beacon is a gRPC service designed for demonstration and testing purposes. It provides a simple beacon service that responds to signal requests with timestamps and server details. The service includes health checks, TLS support, and comprehensive logging.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [gRPC API Reference](#grpc-api-reference)
3. [Configuration](#configuration)
4. [Health Check Service](#health-check-service)
5. [TLS Configuration](#tls-configuration)
6. [Logging](#logging)
7. [Build and Deployment](#build-and-deployment)
8. [Client Examples](#client-examples)
9. [Server Components](#server-components)
10. [Error Handling](#error-handling)

## Architecture Overview

The gRPC Beacon service is built using Go with the following key components:

- **gRPC Service**: Main beacon service with Signal RPC method
- **Health Check**: Standard gRPC health check service
- **Configuration**: Environment variables and TOML configuration file support
- **TLS**: Optional TLS/SSL encryption
- **Logging**: Structured logging with Uber Zap
- **Dependency Injection**: Uber FX framework for dependency management

### Module Structure

```
github.com/troydai/grpcbeacon/
├── cmd/server/          # Main server entry point
├── internal/            # Internal packages
│   ├── beacon/         # Main beacon service
│   ├── health/         # Health check service
│   ├── rpc/            # gRPC server configuration
│   ├── settings/       # Configuration management
│   └── logging/        # Logging configuration
├── proto/              # Protocol Buffer definitions
└── gen/go/             # Generated Go code
```

## gRPC API Reference

### Service: BeaconService

**Package**: `troydai.grpcbeacon.v1`

The BeaconService provides a simple signal endpoint for testing and demonstration purposes.

#### Method: Signal

**Full Method Name**: `/troydai.grpcbeacon.v1.BeaconService/Signal`

**Description**: Accepts a signal request and responds with a timestamp and server details.

**Request Message**: `SignalRequest`
```protobuf
message SignalRequest {
  string message = 1;
}
```

**Response Message**: `SignalResponse`
```protobuf
message SignalResponse {
  string reply = 1;
  map<string, string> details = 10;
}
```

**Go Client Interface**:
```go
type BeaconServiceClient interface {
    Signal(ctx context.Context, in *SignalRequest, opts ...grpc.CallOption) (*SignalResponse, error)
}
```

**Go Server Interface**:
```go
type BeaconServiceServer interface {
    Signal(context.Context, *SignalRequest) (*SignalResponse, error)
    mustEmbedUnimplementedBeaconServiceServer()
}
```

### Message Types

#### SignalRequest

| Field | Type | Description |
|-------|------|-------------|
| message | string | Optional message to include in the request |

#### SignalResponse

| Field | Type | Description |
|-------|------|-------------|
| reply | string | Server response with timestamp |
| details | map<string,string> | Server details including hostname and beacon name |

**Example Response**:
```json
{
  "reply": "Beacon signal at Mon, 02 Jan 2006 15:04:05 MST",
  "details": {
    "Hostname": "server-001",
    "BeaconName": "red cliff"
  }
}
```

## Configuration

The server supports configuration through environment variables and TOML configuration files.

### Environment Variables

| Variable | Type | Description | Default |
|----------|------|-------------|---------|
| `HOSTNAME` | string | Server hostname | System hostname |

### Configuration File

The server accepts a TOML configuration file. Default location: `/etc/beacon-svc/beacon.toml`

Use the `-config` flag to specify a custom configuration file:
```bash
./server -config=/path/to/config.toml
```

#### Configuration Structure

```toml
# Server configuration
name = "beacon-server"      # Beacon name identifier
address = "127.0.0.1"       # Server bind address
port = 8080                 # Server port

# Logging configuration
[logging]
Development = true          # Enable development mode logging

# TLS configuration
[tls]
Enabled = true                          # Enable TLS
KeyFilePath = "certs/server.key.pem"    # Path to private key
CertFilePath = "certs/server.crt.pem"   # Path to certificate
```

### Default Configuration

If no configuration file is found, the server uses these defaults:
```go
Configuration{
    Name:    "red cliff",
    Address: "127.0.0.1",
    Port:    8080,
    Logging: nil,           // Production logging
    TLS:     nil,           // TLS disabled
}
```

### Configuration Types

#### Configuration
```go
type Configuration struct {
    Name    string            `toml:"name"`
    Address string            `toml:"address"`
    Port    int               `toml:"port"`
    Logging *Logging          `toml:"logging"`
    TLS     *TLSConfiguration `toml:"tls"`
}
```

#### Logging
```go
type Logging struct {
    Development bool
}
```

#### TLSConfiguration
```go
type TLSConfiguration struct {
    Enabled      bool
    KeyFilePath  string
    CertFilePath string
}
```

## Health Check Service

The server implements the standard gRPC health check protocol.

### Service: grpc.health.v1.Health

**Methods**:
- `Check(HealthCheckRequest) returns (HealthCheckResponse)`

**Supported Services**:
- `""` (empty string) - Overall server health
- `"liveness"` - Liveness probe
- `"readiness"` - Readiness probe  
- `"Beacon"` - Beacon service health

**Example Usage**:
```bash
# Check overall health
grpcurl --plaintext localhost:8080 grpc.health.v1.Health/Check

# Check specific service
grpcurl --plaintext -d '{"service":"Beacon"}' localhost:8080 grpc.health.v1.Health/Check
```

**Response**:
```json
{
  "status": "SERVING"
}
```

## TLS Configuration

The server supports optional TLS encryption for secure communication.

### Enabling TLS

1. **Configuration File**:
```toml
[tls]
Enabled = true
KeyFilePath = "path/to/private.key"
CertFilePath = "path/to/certificate.crt"
```

2. **Certificate Requirements**:
   - Private key file must be readable
   - Certificate file must be readable
   - Both files must be valid PEM format

### TLS Functions

#### DetermineTLSOption
```go
func DetermineTLSOption(cfg settings.Configuration) (grpc.ServerOption, error)
```

**Description**: Determines the appropriate TLS server option based on configuration.

**Parameters**:
- `cfg`: Configuration object containing TLS settings

**Returns**:
- `grpc.ServerOption`: TLS server option if enabled, nil otherwise
- `error`: Error if TLS configuration is invalid

**Example**:
```go
tlsOpt, err := DetermineTLSOption(config)
if err != nil {
    return fmt.Errorf("TLS configuration error: %w", err)
}
```

## Logging

The server uses structured logging with Uber Zap logger.

### Logger Configuration

#### NewLogger
```go
func NewLogger(c settings.Configuration) (*zap.Logger, error)
```

**Description**: Creates a new Zap logger based on configuration.

**Parameters**:
- `c`: Configuration object with logging settings

**Returns**:
- `*zap.Logger`: Configured logger instance
- `error`: Error if logger creation fails

**Behavior**:
- **Development Mode**: Pretty-printed, colored output with debug level
- **Production Mode**: JSON-formatted output with info level

### Logging in Services

The beacon service logs all incoming requests with metadata:

```go
func (s *service) Signal(ctx context.Context, req *pb.SignalRequest) (*pb.SignalResponse, error) {
    logger := s.logger
    if md, ok := metadata.FromIncomingContext(ctx); ok {
        logger = logger.With(zap.Any("metadata", md))
    }
    
    logger.Info("Signal received")
    // ... implementation
}
```

## Build and Deployment

### Build Commands

#### Build Server
```bash
make bin
```

#### Run Server
```bash
make run
```

#### Run Tests
```bash
make test
```

#### Generate Protocol Buffers
```bash
make gen
```

### Docker Deployment

#### Build Docker Image
```bash
docker build -t grpcbeacon .
```

#### Run Container
```bash
docker run -p 8080:8080 grpcbeacon
```

### Manual Build

```bash
# Generate protobuf code
buf generate

# Build server
go build -o bin/server cmd/server/main.go

# Run server
./bin/server -config=demo/demo.conf
```

## Client Examples

### Go Client

```go
package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    
    pb "github.com/troydai/grpcbeacon/gen/go/troydai/grpcbeacon/v1"
)

func main() {
    // Connect to server
    conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    // Create client
    client := pb.NewBeaconServiceClient(conn)

    // Make request
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    resp, err := client.Signal(ctx, &pb.SignalRequest{
        Message: "Hello from client",
    })
    if err != nil {
        log.Fatalf("Signal failed: %v", err)
    }

    log.Printf("Response: %s", resp.Reply)
    log.Printf("Details: %v", resp.Details)
}
```

### gRPCurl Examples

#### Basic Signal Request
```bash
grpcurl --plaintext localhost:8080 troydai.grpcbeacon.v1.BeaconService.Signal
```

#### Signal with Message
```bash
grpcurl --plaintext -d '{"message": "test"}' localhost:8080 troydai.grpcbeacon.v1.BeaconService.Signal
```

#### TLS Request
```bash
grpcurl --cacert demo/certs/root.crt.pem localhost:8080 troydai.grpcbeacon.v1.BeaconService.Signal
```

#### Health Check
```bash
grpcurl --plaintext localhost:8080 grpc.health.v1.Health/Check
```

### Python Client

```python
import grpc
from gen.python.troydai.grpcbeacon.v1 import api_pb2, api_pb2_grpc

def main():
    # Create channel
    with grpc.insecure_channel('localhost:8080') as channel:
        # Create client
        stub = api_pb2_grpc.BeaconServiceStub(channel)
        
        # Make request
        request = api_pb2.SignalRequest(message="Hello from Python")
        response = stub.Signal(request)
        
        print(f"Response: {response.reply}")
        print(f"Details: {response.details}")

if __name__ == "__main__":
    main()
```

## Server Components

### Beacon Service

#### Service Implementation
```go
type service struct {
    pb.UnimplementedBeaconServiceServer
    details map[string]string
    logger  *zap.Logger
}
```

#### Constructor
```go
func newService(hostName, beaconName string, logger *zap.Logger) *service
```

**Parameters**:
- `hostName`: Server hostname
- `beaconName`: Beacon identifier from configuration
- `logger`: Zap logger instance

### FX Module Registration

#### Beacon Module
```go
var Module = fx.Provide(ProvideRegister)

func ProvideRegister(param Param) Result {
    // Creates and configures beacon service
    // Returns gRPC register function
}
```

#### Module Dependencies
```go
type Param struct {
    fx.In
    Env    settings.Environment
    Config settings.Configuration
    Logger *zap.Logger
}
```

### gRPC Server

#### Server Registration
```go
func RegisterRPCServer(param Param) error
```

**Description**: Configures and starts the gRPC server with all registered services.

**Features**:
- Automatic service registration
- TLS support
- Graceful shutdown
- gRPC reflection

#### GRPCRegister Interface
```go
type GRPCRegister interface {
    Register(*grpc.Server) error
}
```

**Implementation**:
```go
func GRPCRegisterFromFn(fn func(*grpc.Server) error) GRPCRegister
```

## Error Handling

### Configuration Errors

#### Environment Loading
```go
func LoadEnvironment() (Environment, error)
```

**Possible Errors**:
- Environment variable parsing errors
- Missing required environment variables

#### Configuration Loading
```go
func LoadConfig() (Configuration, error)
```

**Possible Errors**:
- File not found (falls back to defaults)
- Invalid TOML format
- File permission errors

### TLS Errors

#### Certificate Loading
```go
func DetermineTLSOption(cfg settings.Configuration) (grpc.ServerOption, error)
```

**Possible Errors**:
- File not found
- Invalid certificate format
- Permission denied
- Certificate/key mismatch

### Service Errors

#### Health Check Errors
```go
func (s *healthcheck) Check(ctx context.Context, req *healthapi.HealthCheckRequest) (*healthapi.HealthCheckResponse, error)
```

**Possible Errors**:
- `codes.NotFound`: Unknown service name

#### Signal Errors
The Signal method currently does not return errors, but logs all requests for debugging.

## Performance Considerations

### Resource Usage
- **Memory**: Minimal memory footprint
- **CPU**: Low CPU usage for simple request/response
- **Connections**: Supports concurrent connections via gRPC

### Scalability
- **Horizontal**: Multiple instances can run behind a load balancer
- **Vertical**: Single instance handles multiple concurrent requests
- **Monitoring**: Health checks enable proper load balancer integration

## Security Considerations

### TLS Configuration
- Use strong certificates from trusted CAs
- Regularly rotate certificates
- Monitor certificate expiration

### Network Security
- Run behind firewall/security groups
- Use TLS for production deployments
- Implement proper authentication if needed

## Troubleshooting

### Common Issues

#### Server Won't Start
1. Check configuration file syntax
2. Verify port availability
3. Check TLS certificate paths
4. Review server logs

#### TLS Issues
1. Verify certificate and key files exist
2. Check file permissions
3. Validate certificate format
4. Ensure certificate/key match

#### Client Connection Issues
1. Verify server address and port
2. Check firewall settings
3. Confirm TLS configuration matches
4. Test with gRPCurl first

### Debug Commands

#### Check Server Health
```bash
grpcurl --plaintext localhost:8080 grpc.health.v1.Health/Check
```

#### Test TLS Connection
```bash
grpcurl --cacert <ca-cert> <server>:<port> troydai.grpcbeacon.v1.BeaconService.Signal
```

#### View Server Logs
```bash
# With development logging
./server -config=demo/demo.conf
```

## Version Information

- **Protocol Buffer Version**: proto3
- **Go Version**: 1.24+
- **gRPC Version**: 1.65.0+
- **Service Version**: v1

## References

- [gRPC Documentation](https://grpc.io/docs/)
- [Protocol Buffers Documentation](https://developers.google.com/protocol-buffers)
- [Uber FX Documentation](https://uber-go.github.io/fx/)
- [Zap Logger Documentation](https://pkg.go.dev/go.uber.org/zap)