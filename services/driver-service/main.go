package main

import (
	"context"
	"golang-ride-sharing/shared/env"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang-ride-sharing/shared/messaging"
	grpcserver "google.golang.org/grpc"
)

var (
	grpcAddr = env.GetString("GRPC_ADDR", ":9092")
)

func main() {
	// env vars
	rabbitMqUri := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")
	// start grpc server with graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", grpcAddr, err)
	}

	service := NewDriverService()

	// RabbitMQ connection
	rabbitMQ, err := messaging.NewRabbitMQ(rabbitMqUri)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
		return
	}
	defer rabbitMQ.Close()

	// start grpc server
	grpcServer := grpcserver.NewServer()
	NewGrpcHandler(grpcServer, service)

	log.Printf("starting gRPC server Driver Service on port %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for shutdown signal
	<- ctx.Done()
	log.Println("shutting down the server gracefully...")
	grpcServer.GracefulStop()
}
