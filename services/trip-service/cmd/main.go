package main

import (
	"context"
	"golang-ride-sharing/services/trip-service/internal/infrastructure/events"
	"golang-ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"golang-ride-sharing/services/trip-service/internal/infrastructure/repository"
	"golang-ride-sharing/services/trip-service/internal/service"
	"golang-ride-sharing/shared/db"
	"golang-ride-sharing/shared/env"
	"golang-ride-sharing/shared/messaging"
	"golang-ride-sharing/shared/tracing"
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
	rabbitmqUri := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")

	// initialzie tracing
	tracerCfg := tracing.Config{
		ServiceName: 	"trip-service",
		Environment: 	env.GetString("ENVIRONMENT", "development"),
		JaegerEndpoint: env.GetString("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces"),
	}

	sh, err := tracing.InitTracer(tracerCfg)
	if err != nil {
		log.Fatalf("failed to start tracer: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer sh(ctx)

	// init mongodb
	mongoClient, err := db.NewMongoClient(ctx, db.NewMongoDefaultConfig())
	if err != nil {
		log.Fatalf("failed to initialize mongodb, err: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	mongoDB := db.GetDatabase(mongoClient, db.NewMongoDefaultConfig())
	log.Printf("connected to mongo db: %s", mongoDB.Name())

	// dependency injections
	mongodbRepo := repository.NewMongoRepository(mongoDB)
	tripService := service.NewTripService(mongodbRepo)


	// start grpc server with graceful shutdown
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

	publisher := events.NewTripEventPublisher(rabbitmq)

	// start driver consumer
	driverResponseConsumer := events.NewDriverResponseConsumer(rabbitmq, tripService)
	go driverResponseConsumer.Listen()

	// start payment consumer
	paymentConsumer := events.NewPaymentConsumer(rabbitmq, tripService)
	go paymentConsumer.Listen()

	// starting grpc server
	grpcServer := grpcserver.NewServer()
	grpc.NewGrpcHandler(grpcServer, tripService, publisher)

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
