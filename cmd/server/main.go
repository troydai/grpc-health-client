package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	settings, err := parseArguments()
	if err != nil {
		logger.Error("failed to parse arguments", "error", err)
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", settings.Address())
	if err != nil {
		logger.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the health server
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	logger.Debug("Server listening on", "port", settings.Address())
	if err := server.Serve(lis); err != nil {
		logger.Error("failed to serve", "error", err)
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
