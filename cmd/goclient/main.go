// The client keeps sending signal to the server at specified interval.
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/troydai/cron"
	api "github.com/troydai/grpcbeacon/gen/api/protos"
)

const (
	_envServerAddress = "SERVER_ADDRESS"
	_envInterval      = "CLIENT_INTERVAL"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("fail to create logger: %v", err)
	}

	serverAddr, err := resolveServerAddr()
	if err != nil {
		logger.Fatal("fail to resolve server address", zap.Error(err))
	}

	interval, err := resolveInterval()
	if err != nil {
		logger.Fatal("fail to resolve interval", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	term, err := cron.Start(ctx, createSignalJob(serverAddr, logger), cron.WithInterval(interval))
	if err != nil {
		logger.Fatal("fail to start cron", zap.Error(err))
	}

	go func() {
		chSystemSignal := make(chan os.Signal, 1)
		signal.Notify(chSystemSignal, os.Interrupt)

		select {
		case <-term:
		case <-chSystemSignal:
			cancel()
		}
	}()

	<-term
}

func createSignalJob(serverAddr string, logger *zap.Logger) cron.Job {
	return func(ctx context.Context) bool {
		logger.Info("Connecting to server", zap.String("server", serverAddr))

		// TODO: make whether to maintain the connection configuration
		conn, err := grpc.DialContext(ctx, serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error("fail to connect to server", zap.Error(err))
		}
		defer conn.Close()

		client := api.NewBeaconClient(conn)
		resp, err := client.Signal(ctx, &api.SignalReqeust{})
		if err != nil {
			logger.Error("fail to signal the server", zap.Error(err))
		}

		logger.Info("Server response", zap.String("reply", resp.GetReply()), zap.Any("details", resp.GetDetails()))

		return true
	}
}

func resolveServerAddr() (string, error) {
	addr := os.Getenv(_envServerAddress)
	if addr == "" {
		return "", fmt.Errorf("server address is not specified")
	}

	if _, err := url.Parse(addr); err != nil {
		return "", err
	}

	return addr, nil
}

func resolveInterval() (time.Duration, error) {
	interval := os.Getenv(_envInterval)
	if interval == "" {
		// default interval is set here
		return 5 * time.Second, nil
	}

	d, err := time.ParseDuration(interval)
	if err != nil {
		return 0, err
	}

	return d, nil
}
