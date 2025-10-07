package main

import (
	"context"
	"golang-ride-sharing/shared/env"
	"golang-ride-sharing/shared/messaging"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang-ride-sharing/services/payment-service/internal/infrastructure/grpc"
	grpcserver "google.golang.org/grpc"
)

var (
	grpcAddr = env.GetString("GRPC_ADDR", ":9094")
)

func main() {
	// env vars
	rabbitmqUri := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")

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

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitmqUri)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
		return
	}
	defer rabbitmq.Close()

	// starting grpc server
	grpcServer := grpcserver.NewServer()
	grpc.NewGrpcHandler(grpcServer)

	log.Printf("starting gRPC server Payment Service on port %s", lis.Addr().String())

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
