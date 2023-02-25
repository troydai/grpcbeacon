package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/troydai/grpcbeacon/gen/api/protos"
	"github.com/troydai/grpcbeacon/internal/beacon"
	"github.com/troydai/grpcbeacon/internal/settings"
)

func main() {
	env, err := settings.LoadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(fmt.Errorf("fail to create logger: %w", err))
	}

	server := grpc.NewServer()
	api.RegisterBeaconServer(server, beacon.NewServer(env))
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Error("fail to start TCP listener", zap.Error(err))
	}

	chServerStopped := make(chan struct{})
	chSystemSignal := make(chan os.Signal, 1)

	signal.Notify(chSystemSignal, os.Interrupt)

	go func() {
		logger.Info("Listening system signals")

		select {
		case <-chServerStopped:
		case s := <-chSystemSignal:
			logger.Info("Stopping server after recieving system signal", zap.String("signal", s.String()))
			server.GracefulStop()
		}
	}()

	go func() {
		defer close(chServerStopped)
		logger.Info("Starting server")
		server.Serve(lis)
	}()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				logger.Info("Server is running")
			case <-chServerStopped:
				logger.Info("Ticker stopped")
				return
			}
		}
	}()

	<-chServerStopped
}
