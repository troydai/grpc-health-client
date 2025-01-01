package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	port := resolvePort()
	if err := probe(port); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func probe(port int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		fmt.Sprintf("127.0.0.1:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Create health check client
	healthClient := healthpb.NewHealthClient(conn)

	// Perform health check
	resp, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		return fmt.Errorf("health check failed: %v", resp.GetStatus())
	}

	return nil
}

func resolvePort() int {
	port := os.Getenv("GRPC_HEALTH_PORT")
	if port == "" {
		return 50051
	}

	portn, err := strconv.Atoi(port)
	if err != nil {
		return 50051
	}

	return portn
}
