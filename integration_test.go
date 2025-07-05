package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/troydai/grpcbeacon/gen/go/troydai/grpcbeacon/v1"
	"github.com/troydai/grpcbeacon/internal/beacon"
	"github.com/troydai/grpcbeacon/internal/health"
	"github.com/troydai/grpcbeacon/internal/logging"
	"github.com/troydai/grpcbeacon/internal/rpc"
	"github.com/troydai/grpcbeacon/internal/settings"
)

func TestIntegration_ServerStartsAndEchoWorks(t *testing.T) {
	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())

	// Override configuration for testing
	testConfig := settings.Configuration{
		Name:    "test-beacon",
		Address: "127.0.0.1",
		Port:    port,
	}

	testEnv := settings.Environment{
		HostName: "test-host",
	}

	// Create the fx app with test configuration
	app := fxtest.New(t,
		// Override the settings providers with test values
		fx.Provide(func() settings.Configuration { return testConfig }),
		fx.Provide(func() settings.Environment { return testEnv }),

		// Include all the necessary modules
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

	// Test the gRPC service
	t.Run("Signal endpoint works", func(t *testing.T) {
		// Connect to the server
		conn, err := grpc.NewClient(
			fmt.Sprintf("127.0.0.1:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)
		defer func() { require.NoError(t, conn.Close()) }()

		client := pb.NewBeaconServiceClient(conn)

		// Test the Signal method (echo endpoint)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		testMessage := "Hello, gRPC Beacon!"
		req := &pb.SignalRequest{
			Message: testMessage,
		}

		resp, err := client.Signal(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Verify the response
		assert.NotEmpty(t, resp.Reply)
		assert.Contains(t, resp.Reply, "Beacon signal at")

		// Verify the details
		assert.NotNil(t, resp.Details)
		assert.Equal(t, "test-host", resp.Details["Hostname"])
		assert.Equal(t, "test-beacon", resp.Details["BeaconName"])
	})

	// Test multiple concurrent calls
	t.Run("Concurrent calls work", func(t *testing.T) {
		conn, err := grpc.NewClient(
			fmt.Sprintf("127.0.0.1:%d", port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)
		defer func() { require.NoError(t, conn.Close()) }()

		client := pb.NewBeaconServiceClient(conn)

		const numCalls = 5
		var wg sync.WaitGroup
		responses := make([]*pb.SignalResponse, numCalls)
		errors := make([]error, numCalls)

		for i := 0; i < numCalls; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				req := &pb.SignalRequest{
					Message: fmt.Sprintf("Message %d", idx),
				}

				responses[idx], errors[idx] = client.Signal(ctx, req)
			}(i)
		}

		wg.Wait()

		// Verify all calls succeeded
		for i := 0; i < numCalls; i++ {
			assert.NoError(t, errors[i], "Call %d failed", i)
			assert.NotNil(t, responses[i], "Response %d is nil", i)
			if responses[i] != nil {
				assert.NotEmpty(t, responses[i].Reply, "Reply %d is empty", i)
				assert.Contains(t, responses[i].Reply, "Beacon signal at", "Reply %d doesn't contain expected text", i)
			}
		}
	})

	// Clean shutdown
	stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, app.Stop(stopCtx))
}

func TestIntegration_ServerConfiguration(t *testing.T) {
	// Test that the server respects configuration changes
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())

	testConfig := settings.Configuration{
		Name:    "custom-beacon-name",
		Address: "127.0.0.1",
		Port:    port,
	}

	testEnv := settings.Environment{
		HostName: "custom-hostname",
	}

	app := fxtest.New(t,
		fx.Provide(func() settings.Configuration { return testConfig }),
		fx.Provide(func() settings.Environment { return testEnv }),
		logging.Module,
		rpc.Module,
		beacon.Module,
		health.Module,
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, app.Start(startCtx))

	time.Sleep(100 * time.Millisecond)

	// Test that configuration is reflected in the response
	conn, err := grpc.NewClient(
		fmt.Sprintf("127.0.0.1:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer func() { require.NoError(t, conn.Close()) }()

	client := pb.NewBeaconServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Signal(ctx, &pb.SignalRequest{Message: "test"})
	require.NoError(t, err)

	// Verify custom configuration is used
	assert.Equal(t, "custom-hostname", resp.Details["Hostname"])
	assert.Equal(t, "custom-beacon-name", resp.Details["BeaconName"])

	stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, app.Stop(stopCtx))
}
