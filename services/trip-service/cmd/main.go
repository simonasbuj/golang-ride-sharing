package main

import (
	"context"
	"golang-ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"golang-ride-sharing/services/trip-service/internal/infrastructure/repository"
	"golang-ride-sharing/services/trip-service/internal/service"
	"golang-ride-sharing/shared/env"
	"golang-ride-sharing/shared/messaging"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var (
	grpcAddr = env.GetString("GRPC_ADDR", ":9093")
)

func main() {
	// env vars
	rabbitMqUri := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")

	// dependency injections
	inmemoryRepo := repository.NewInmemoryRepository()
	tripService := service.NewTripService(inmemoryRepo)


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
	rabbitMQ, err := messaging.NewRabbitMQ(rabbitMqUri)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
		return
	}
	defer rabbitMQ.Close()

	// starting grpc server
	grpcServer := grpcserver.NewServer()
	grpc.NewGrpcHandler(grpcServer, tripService)

	log.Printf("starting gRPC server Trip Service on port %s", lis.Addr().String())

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
