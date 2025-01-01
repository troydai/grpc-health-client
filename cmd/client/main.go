package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	_defaultServerAddr = "localhost:50051"
)

func main() {
	// Set up the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	settings, err := parseArguments()
	if err != nil {
		logger.Error("failed to parse arguments", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		settings.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("failed to connect", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create health check client
	healthClient := healthpb.NewHealthClient(conn)

	// Perform health check
	resp, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		logger.Error("health check failed", "error", err)
		os.Exit(1)
	}

	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		logger.Error("health check failed", "status", resp.GetStatus())
		os.Exit(1)
	}
}

type settings struct {
	port int
}

func (s settings) Address() string {
	return fmt.Sprintf("localhost:%d", s.port)
}

func parseArguments() (settings, error) {
	s := settings{
		port: 50051, // default port
	}

	args := os.Args[1:] // skip program name
	if len(args) == 0 {
		return s, nil
	}

	if len(args) > 1 {
		return settings{}, fmt.Errorf("usage: %s [port]", os.Args[0])
	}

	port, err := strconv.Atoi(args[0])
	if err != nil {
		return settings{}, fmt.Errorf("invalid port number: %v", err)
	}

	s.port = port
	return s, nil
}
