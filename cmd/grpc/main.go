package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	}

	log.Println("Starting gRPC server...")

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Enable reflection for development (helps with tools like grpcurl)
	reflection.Register(grpcServer)
	log.Println("gRPC reflection enabled")

	// Create listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	// Start server in goroutine
	go func() {
		log.Println("gRPC server starting on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	fmt.Println("gRPC server is running on port 50051")
	fmt.Println("Note: This is a basic gRPC server setup.")
	fmt.Println("To complete the setup:")
	fmt.Println("1. Define your .proto files in the proto/ directory")
	fmt.Println("2. Generate Go code using protoc")
	fmt.Println("3. Implement your gRPC service handlers")
	fmt.Println("4. Register your services with the server")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a channel to signal shutdown completion
	done := make(chan bool, 1)

	go func() {
		grpcServer.GracefulStop()
		done <- true
	}()

	select {
	case <-done:
		log.Println("gRPC server shutdown completed")
	case <-ctx.Done():
		log.Println("gRPC server shutdown timeout, forcing stop")
		grpcServer.Stop()
	}

	log.Println("gRPC server stopped")
}
