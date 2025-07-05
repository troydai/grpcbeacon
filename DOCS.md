# gRPC Beacon – Public API Documentation

This document serves as the authoritative reference for all **public-facing APIs, packages, and components** that make up the `grpcbeacon` project.  Everything listed here is safe to import or consume from other programs or scripts.  Internal details (unexported identifiers or packages inside `internal/...`) are intentionally omitted unless they surface through a public interface.

---

## Table of Contents

1. Beacon gRPC Service (proto)
2. Go Packages
   * `troydai/grpcbeacon/settings`
   * `troydai/grpcbeacon/rpc`
   * `troydai/grpcbeacon/logging`
3. Command-line Applications
   * `cmd/server`
4. Configuration & Environment
5. End-to-End Quick-start

---

## 1. Beacon gRPC Service

### Package: `troydai.grpcbeacon.v1`

```protobuf
service BeaconService {
  rpc Signal (SignalRequest) returns (SignalResponse);
}

message SignalRequest {
  string message = 1;
}

message SignalResponse {
  string  reply   = 1;
  map<string,string> details = 10;
}
```

#### Semantics
* **Signal** –  A simple liveness endpoint.  The server echoes a timestamped reply and a set of key–value pairs that describe the running beacon instance (hostname, beacon name, etc.).

#### Example – Invoke with `grpcurl`
```bash
# Assuming the server is listening locally with TLS disabled
grpcurl -plaintext localhost:8089 \
  troydai.grpcbeacon.v1.BeaconService.Signal \
  '{"message":"ping"}'
```

Example response:
```json
{
  "reply": "Beacon signal at Tue, 05 Mar 2024 15:04:05 UTC",
  "details": {
    "Hostname": "ip-10-0-0-1",
    "BeaconName": "red cliff"
  }
}
```

---

## 2. Go Packages

### 2.1 `settings`

```go
import "github.com/troydai/grpcbeacon/internal/settings"
```

Although the package lives in `internal/`, the exported *types* are used by other public-facing modules and are therefore documented here.

| Export | Description |
|--------|-------------|
| `Environment` | Struct populated from environment variables (currently only `HOSTNAME`). |
| `Configuration` | User-supplied TOML or default configuration. Fields: `Name`, `Address`, `Port`, `Logging`, `TLS`. |
| `LoadEnvironment()` | Parse env vars into an `Environment` value. |
| `LoadConfig()` | Locate, parse, and validate the configuration file (CLI flag `-config`, falls back to `/etc/beacon-svc/beacon.toml`). |

#### Example – Load configuration
```go
cfg, err := settings.LoadConfig()
if err != nil { /* handle */ }
log.Printf("Beacon will start on %s:%d", cfg.Address, cfg.Port)
```

### 2.2 `rpc`

```go
import "github.com/troydai/grpcbeacon/internal/rpc"
```

| Export | Description |
|--------|-------------|
| `GRPCRegister` | Adapter that allows features to register themselves with a `grpc.Server`. |
| `GRPCRegisterFromFn(fn)` | Convenience helper turning an anonymous func into a `GRPCRegister`. |
| `DetermineTLSOption(cfg)` | Produce a `grpc.ServerOption` enabling mTLS if configured. |
| `Module` | [Uber-fx](https://github.com/uber-go/fx) module that boots a gRPC server and ties lifecycle to the fx app. |

#### Example – Custom service registration
```go
fx.Provide(func() rpc.GRPCRegister {
    return rpc.GRPCRegisterFromFn(func(s *grpc.Server) error {
        pb.RegisterMyServiceServer(s, myServiceImpl{})
        return nil
    })
})
```

### 2.3 `logging`

```go
import "github.com/troydai/grpcbeacon/internal/logging"
```

| Export | Description |
|--------|-------------|
| `NewLogger(cfg)` | Returns a *production* or *development* `zap.Logger` depending on `cfg.Logging.Development`. |
| `Module` | fx module wiring `NewLogger` and bridging it into fx's event logger. |

---

## 3. Command-line Applications

### 3.1 `grpcbeacon` Server (`cmd/server`)

Build:
```bash
go build -o bin/grpcbeacon ./cmd/server
```

Run with defaults:
```bash
./bin/grpcbeacon -config ./conf/beacon.conf
```

Environment variables honoured:
* `HOSTNAME` – auto-detected if not set.

The process starts:
* gRPC server on `Address:Port` (defaults `127.0.0.1:8080`).
* Health endpoints following the standard gRPC Health Checking Protocol.

### 3.2 Helper Scripts

* `scripts/update_certs.sh` – Generate self-signed TLS certificates.
* `cmd/toolbox/curl.sh` – Minimal example using `grpcurl`.

---

## 4. Configuration & Environment

```toml
# conf/beacon.conf
name    = "white peak"
address = "127.0.0.1"
port    = 6899

[logging]
Development = true
```

TLS can be enabled by extending the config:
```toml
[tls]
enabled       = true
keyFilePath   = "./demo/certs/server.key.pem"
certFilePath  = "./demo/certs/server.crt.pem"
```

---

## 5. End-to-End Quick-start

```bash
# 1. Build & run the server
make run  # → defaults to go run ./cmd/server

# 2. (Optional) generate TLS certs if you enabled TLS
./scripts/update_certs.sh ./demo/certs

# 3. Ping the service
grpcurl --plaintext localhost:8089 troydai.grpcbeacon.v1.BeaconService.Signal
```

If TLS is enabled:
```bash
grpcurl \
  -cacert ./demo/certs/root.crt.pem \
  localhost:8089 \
  troydai.grpcbeacon.v1.BeaconService.Signal
```

---

## Contributing Documentation
If you add a new exported function or a new gRPC endpoint, **please update this file** with:
1. A concise description
2. Parameter & response schema (for gRPC or REST)
3. A runnable example (shell or Go)

Thank you!