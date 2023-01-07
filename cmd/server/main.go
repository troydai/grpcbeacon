package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/troydai/grpcbeacon/api/protos"
	"github.com/troydai/grpcbeacon/internal/beacon"
)

func main() {
	server := grpc.NewServer()
	api.RegisterBeaconServer(server, &beacon.Server{})
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(fmt.Errorf("fail to start TCP listener: %w", err))
	}

	chServerStopped := make(chan struct{})
	chSystemSignal := make(chan os.Signal, 1)

	signal.Notify(chSystemSignal)

	go func() {
		select {
		case <-chServerStopped:
		case <-chSystemSignal:
			server.GracefulStop()
		}
	}()

	go func() {
		defer close(chServerStopped)
		server.Serve(lis)
	}()

	<-chServerStopped
}
