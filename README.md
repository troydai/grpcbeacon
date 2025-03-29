# gRPC Beacon

A grpc service that respond to the request. Utility for demo purposes

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

## References

- Image registry: https://hub.docker.com/repository/docker/troydai/grpcbeacon
