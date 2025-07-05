# Integration Tests

This document describes the comprehensive integration tests created for the gRPC Beacon service to ensure the server starts correctly and all endpoints work as expected.

## Overview

The integration tests verify the complete functionality of the gRPC Beacon service by:
- Starting the actual server with all dependencies
- Creating real gRPC client connections  
- Testing all available endpoints
- Validating responses and error conditions
- Ensuring proper server lifecycle management

## Test Files

### 1. `integration_test.go` - Main Service Tests

**Primary Test: `TestIntegration_ServerStartsAndEchoWorks`**
- ✅ Starts the gRPC server on a random available port
- ✅ Tests the `Signal` endpoint (echo functionality)
- ✅ Verifies response content and metadata
- ✅ Tests concurrent calls to ensure thread safety
- ✅ Validates proper server startup and shutdown

**Configuration Test: `TestIntegration_ServerConfiguration`**
- ✅ Tests custom configuration is properly applied
- ✅ Verifies hostname and beacon name are reflected in responses
- ✅ Ensures configuration dependency injection works correctly

### 2. `health_integration_test.go` - Health Check Tests

**Health Check Test: `TestIntegration_HealthCheck`**
- ✅ Tests general health check endpoint (`""` service)
- ✅ Tests specific service health checks (`"Beacon"`, `"liveness"`, `"readiness"`)
- ✅ Validates error handling for unknown services
- ✅ Ensures all health endpoints return `SERVING` status

## Test Features

### Server Management
- **Dynamic Port Allocation**: Tests use random available ports to avoid conflicts
- **Proper Lifecycle**: Tests start and stop servers cleanly with timeouts
- **Dependency Injection**: Uses Uber FX with test-specific configurations

### Realistic Testing
- **Real gRPC Connections**: Tests use actual gRPC clients, not mocks
- **Network Communication**: Full network stack testing over TCP
- **Concurrent Access**: Multi-goroutine testing for race conditions

### Comprehensive Coverage
- **Signal Endpoint**: The main "echo" functionality with timestamp responses
- **Health Endpoints**: All health check variants including error cases
- **Configuration**: Custom settings and environment variables
- **Error Handling**: Invalid requests and unknown services

## Running the Tests

```bash
# Run all integration tests
go test -v -run TestIntegration

# Run only main service tests
go test -v -run TestIntegration_ServerStartsAndEchoWorks

# Run only health check tests  
go test -v -run TestIntegration_HealthCheck

# Run with verbose logging to see server startup
go test -v -run TestIntegration 2>&1 | grep -E "(RUN|PASS|FAIL|Signal received)"
```

## Test Architecture

### Dependencies Used
- **testify**: Assertions and test utilities (`assert`, `require`)
- **fx/fxtest**: Dependency injection testing framework
- **grpc**: Real gRPC client/server communication
- **net**: Dynamic port allocation for test isolation

### Test Structure
```go
// 1. Find available port
listener, err := net.Listen("tcp", "127.0.0.1:0")
port := listener.Addr().(*net.TCPAddr).Port
listener.Close()

// 2. Create test configuration
testConfig := settings.Configuration{
    Name:    "test-beacon",
    Address: "127.0.0.1", 
    Port:    port,
}

// 3. Start application with test config
app := fxtest.New(t,
    fx.Provide(func() settings.Configuration { return testConfig }),
    // ... all modules
)
app.Start(ctx)

// 4. Test endpoints
conn, _ := grpc.NewClient(fmt.Sprintf("127.0.0.1:%d", port))
client := pb.NewBeaconServiceClient(conn)
resp, err := client.Signal(ctx, &pb.SignalRequest{Message: "test"})

// 5. Clean shutdown
app.Stop(ctx)
```

## Expected Behavior

### Signal Endpoint
- **Input**: `SignalRequest{Message: "any string"}`
- **Output**: `SignalResponse{Reply: "Beacon signal at <timestamp>", Details: {...}}`
- **Details Include**: Hostname, BeaconName from configuration

### Health Endpoints
- **General Health** (`""`): Returns `SERVING`
- **Beacon Health** (`"Beacon"`): Returns `SERVING`  
- **Liveness** (`"liveness"`): Returns `SERVING`
- **Readiness** (`"readiness"`): Returns `SERVING`
- **Unknown Service**: Returns `NotFound` error

## Benefits

1. **Confidence**: Tests verify the complete system works end-to-end
2. **Regression Prevention**: Catches breaking changes in server startup or endpoints
3. **Documentation**: Tests serve as executable documentation of API behavior
4. **CI/CD Ready**: Tests can run in any environment with available ports
5. **Performance Validation**: Concurrent testing ensures scalability

## Integration with CI/CD

These tests are designed to run in CI/CD pipelines:
- No external dependencies (only localhost networking)
- Dynamic port allocation prevents conflicts
- Deterministic test outcomes
- Fast execution (typically under 1 second total)
- Clear pass/fail indicators

The integration tests ensure that the gRPC Beacon service is production-ready and all endpoints function correctly under realistic conditions.