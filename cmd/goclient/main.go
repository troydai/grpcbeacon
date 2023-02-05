package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/troydai/grpcbeacon/gen/api/protos"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	address := resolveAddress()
	fmt.Printf("Connecting to: %s\n", address)

	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail top dial the server: %v", err)
	}
	defer conn.Close()

	client := api.NewBeaconClient(conn)
	resp, err := client.Signal(ctx, &api.SignalReqeust{})
	if err != nil {
		log.Fatalf("fail to signal the server: %v", err)
	}
	fmt.Printf("Server replied: %s", resp)
}

func resolveAddress() string {
	if len(os.Args) < 2 {
		return "localhost:8080"
	}

	return os.Args[1]
}
