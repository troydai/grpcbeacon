package main

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	healthpb "github.com/troydai/grpcbeacon/gen/go/grpc/health/v1"
	"github.com/troydai/grpcbeacon/internal/beacon"
	"github.com/troydai/grpcbeacon/internal/health"
	"github.com/troydai/grpcbeacon/internal/logging"
	"github.com/troydai/grpcbeacon/internal/rpc"
	"github.com/troydai/grpcbeacon/internal/settings"
)

func TestIntegration_HealthCheck(t *testing.T) {
	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())

	// Override configuration for testing
	testConfig := settings.Configuration{
		Name:    "health-test-beacon",
		Address: "127.0.0.1",
		Port:    port,
	}

	testEnv := settings.Environment{
		HostName: "health-test-host",
	}

	// Create the fx app with test configuration
	app := fxtest.New(t,
		fx.Provide(func() settings.Configuration { return testConfig }),
		fx.Provide(func() settings.Environment { return testEnv }),
		logging.Module,
		rpc.Module,
		beacon.Module,
		health.Module,
	)

	// Start the application
	startCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	require.NoError(t, app.Start(startCtx))

	// Wait a bit for the server to be ready
	time.Sleep(100 * time.Millisecond)

	// Test the health check service
	t.Run("Health check endpoint works", func(t *testing.T) {
		// Connect to the server
		conn, err := grpc.NewClient(
			fmt.Sprintf("127.0.0.1:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)
		defer func() { require.NoError(t, conn.Close()) }()

		client := healthpb.NewHealthClient(conn)

		// Test the Check method
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &healthpb.HealthCheckRequest{
			Service: "", // Empty service name for overall health
		}

		resp, err := client.Check(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Verify the response indicates the service is serving
		assert.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.Status)
	})

	t.Run("Health check for specific service", func(t *testing.T) {
		conn, err := grpc.NewClient(
			fmt.Sprintf("127.0.0.1:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)
		defer func() { require.NoError(t, conn.Close()) }()

		client := healthpb.NewHealthClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Test health check for the beacon service specifically (using the correct service name)
		req := &healthpb.HealthCheckRequest{
			Service: "Beacon",
		}

		resp, err := client.Check(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Should also be serving
		assert.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.Status)
	})

	t.Run("Health check for liveness", func(t *testing.T) {
		conn, err := grpc.NewClient(
			fmt.Sprintf("127.0.0.1:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)
		defer func() { require.NoError(t, conn.Close()) }()

		client := healthpb.NewHealthClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &healthpb.HealthCheckRequest{
			Service: "liveness",
		}

		resp, err := client.Check(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.Status)
	})

	t.Run("Health check for readiness", func(t *testing.T) {
		conn, err := grpc.NewClient(
			fmt.Sprintf("127.0.0.1:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)
		defer func() { require.NoError(t, conn.Close()) }()

		client := healthpb.NewHealthClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &healthpb.HealthCheckRequest{
			Service: "readiness",
		}

		resp, err := client.Check(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.Status)
	})

	t.Run("Health check for unknown service returns not found", func(t *testing.T) {
		conn, err := grpc.NewClient(
			fmt.Sprintf("127.0.0.1:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)
		defer func() { require.NoError(t, conn.Close()) }()

		client := healthpb.NewHealthClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &healthpb.HealthCheckRequest{
			Service: "unknown-service",
		}

		_, err = client.Check(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "service unknown-service not found")
	})

	// Clean shutdown
	stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, app.Stop(stopCtx))
}
