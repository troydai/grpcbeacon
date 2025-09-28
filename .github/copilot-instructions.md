# Copilot Instructions for gRPC Beacon

This document provides guidance for GitHub Copilot when working on the gRPC Beacon project. It helps Copilot understand the project structure, coding patterns, and best practices to provide better assistance.

## Project Overview

gRPC Beacon is a demonstration gRPC service written in Go that provides:
- **Beacon Service**: Simple echo/signal functionality with server metadata
- **Health Check Service**: Standard gRPC health checking
- **TLS Support**: Optional TLS configuration for secure connections
- **Comprehensive Testing**: Unit and integration tests with real server startup

## Architecture & Tech Stack

- **Language**: Go 1.24
- **Framework**: gRPC with Protocol Buffers
- **Dependency Injection**: Uber FX framework
- **Logging**: Uber Zap structured logging
- **Configuration**: TOML files + environment variables
- **Build Tools**: Make, buf (Protocol Buffers), Docker
- **Testing**: Go's testing package with comprehensive integration tests

## File Structure & Patterns

```
├── cmd/server/           # Main server application entry point
├── internal/             # Private application packages
│   ├── beacon/          # Core beacon service implementation
│   ├── health/          # Health check service
│   ├── logging/         # Logging configuration
│   ├── rpc/             # gRPC server setup and registration
│   └── settings/        # Configuration loading
├── proto/               # Protocol Buffer definitions
├── gen/                 # Generated code (protobuf)
├── demo/                # Demo configurations and certificates
└── *.go                # Integration tests at root level
```

## Coding Standards & Patterns

### Go Code Style
- Follow standard Go formatting (`gofmt -s`)
- Use `go vet` for static analysis
- All exported functions/types must have comments
- Error handling: Always check and handle errors appropriately
- Use structured logging with Zap logger

### gRPC Service Patterns
```go
// Service implementation pattern
type BeaconServer struct {
    logger *zap.Logger
    // Embed unimplemented server for forward compatibility
    pb.UnimplementedBeaconServiceServer
}

// Service method pattern with proper error handling
func (s *BeaconServer) Signal(ctx context.Context, req *pb.SignalRequest) (*pb.SignalResponse, error) {
    s.logger.Info("received signal", zap.String("message", req.Message))
    // Implementation...
    return response, nil
}
```

### Dependency Injection with FX
- Use Uber FX for dependency injection
- Create modules in `fx.go` files
- Provide dependencies using `fx.Provide()`
- Use lifecycle hooks for startup/shutdown

```go
// Module pattern
var Module = fx.Module("modulename",
    fx.Provide(ProvideService),
    fx.Invoke(RegisterService),
)
```

### Configuration Pattern
- Configuration structs use TOML tags
- Environment variables override config files
- Default values provided in struct tags or code

```go
type Config struct {
    Name    string `toml:"name" env:"BEACON_NAME" envDefault:"beacon"`
    Address string `toml:"address" env:"BEACON_ADDRESS" envDefault:"127.0.0.1"`
    Port    int    `toml:"port" env:"BEACON_PORT" envDefault:"8080"`
}
```

## Testing Guidelines

### Integration Tests
- Located at repository root (`integration_test.go`, `health_integration_test.go`)
- Start real gRPC server using FX framework
- Test actual client-server communication
- Include concurrent access testing
- Verify error handling and edge cases

```go
// Integration test pattern
func TestIntegration_ServiceName(t *testing.T) {
    // Create test configuration
    app := fxtest.New(t,
        // Include all required modules
        fx.Provide(testConfig),
        // Register test dependencies
    )
    defer app.RequireStop()
    app.RequireStart()
    
    // Test implementation...
}
```

### Unit Tests
- Test individual components in isolation
- Mock external dependencies
- Focus on business logic and error conditions

## Build & Development Workflow

### Essential Commands
```bash
# FIRST: Install Protocol Buffer tools (required for code generation)
make tools

# Generate protobuf code (requires network access to buf.build)
make gen

# Build server binary (requires protobuf code generation first)
make bin

# Run with demo configuration
make run

# Run all tests (unit + integration)
make test
# or: go test -v ./...

# Run only integration tests
go test -v -run TestIntegration

# Format code
gofmt -s -w .

# Static analysis
go vet ./...
```

### Important Prerequisites
1. **Go 1.24+** installed and configured
2. **Network access** to buf.build for protobuf generation
3. **GOPATH/GOBIN** properly configured (buf installs to `~/go/bin/`)
4. **PATH** should include Go's bin directory

### Development Setup Steps
```bash
# 1. Install tools first
make tools

# 2. Ensure buf is in PATH (if needed)
export PATH="$PATH:$(go env GOPATH)/bin"

# 3. Generate protobuf code
make gen

# 4. Now you can build and test
make bin
go test -v ./...
```

### Docker Development
```bash
# Build container image
docker build -t grpcbeacon .

# Run container
docker run -p 8080:8080 grpcbeacon

# Use toolbox container for testing
docker build --target toolbox -t grpcbeacon-toolbox .
```

## Protocol Buffers

- Service definitions in `proto/troydai/grpcbeacon/v1/`
- Use `buf` tool for generation and linting
- Generated code goes in `gen/` directory
- Follow protobuf style guide for message and service naming

## Common Patterns to Follow

### Error Handling
```go
// Always check errors
if err != nil {
    s.logger.Error("operation failed", zap.Error(err))
    return nil, status.Errorf(codes.Internal, "internal error: %v", err)
}
```

### Logging
```go
// Use structured logging
s.logger.Info("processing request", 
    zap.String("method", "Signal"),
    zap.String("message", req.Message),
)
```

### Server Registration
```go
// Register services with gRPC server
func RegisterService(server *grpc.Server, service *Service) {
    pb.RegisterBeaconServiceServer(server, service)
}
```

## Testing Client Connections

Use `grpcurl` for manual testing:
```bash
# Basic health check
grpcurl --plaintext localhost:8080 grpc.health.v1.Health/Check

# Signal request
grpcurl --plaintext -d '{"message": "test"}' localhost:8080 troydai.grpcbeacon.v1.BeaconService.Signal

# With TLS
grpcurl --cacert demo/certs/root.crt.pem localhost:8080 troydai.grpcbeacon.v1.BeaconService.Signal
```

## CI/CD Integration

- GitHub Actions workflows test, build, and publish
- All tests must pass before merge
- Security scanning with gosec
- Multi-platform Docker builds
- Code coverage reporting

## Common Pitfalls to Avoid

1. **Forgetting to check errors** - All error returns must be handled
2. **Not using structured logging** - Use Zap with proper fields
3. **Missing protobuf regeneration** - Run `make gen` after proto changes
4. **Integration test port conflicts** - Tests handle dynamic port allocation
5. **Missing FX lifecycle hooks** - Properly start/stop services
6. **TLS certificate issues** - Use provided demo certs or regenerate with scripts

## Additional Resources

- **API Documentation**: See `API_DOCUMENTATION.md` for complete API reference
- **Integration Tests**: See `INTEGRATION_TESTS.md` for testing details
- **CI/CD Setup**: See `CI_CD_SETUP.md` for pipeline information
- **Project Summary**: See `SUMMARY.md` for recent updates and features

## When Making Changes

1. Update protocol buffers first if API changes are needed
2. Run `make gen` to regenerate code
3. Update or add tests (both unit and integration)
4. Ensure all tests pass: `go test -v ./...`
5. Format code: `gofmt -s -w .`
6. Run static analysis: `go vet ./...`
7. Update documentation if needed
8. Verify Docker build works: `docker build .`

This guidance helps Copilot understand the project's patterns and suggest appropriate code that follows established conventions.