package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang-ride-sharing/services/payment-service/internal/events"
	"golang-ride-sharing/services/payment-service/internal/infrastructure/stripe"
	"golang-ride-sharing/services/payment-service/internal/service"
	"golang-ride-sharing/services/payment-service/types"
	"golang-ride-sharing/shared/env"
	"golang-ride-sharing/shared/messaging"
	"golang-ride-sharing/shared/tracing"
)

func main() {
	// GrpcAddr := env.GetString("GRPC_ADDR", ":9004")
	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")

	// initialzie tracing
	tracerCfg := tracing.Config{
		ServiceName: 	"payment-service",
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

	// Setup graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	appURL := env.GetString("APP_URL", "http://localhost:3000")

	// Stripe config and processor
	stripeCfg := &types.PaymentConfig{
		StripeSecretKey: env.GetString("STRIPE_SECRET_KEY", ""),
		SuccessURL:      env.GetString("STRIPE_SUCCESS_URL", appURL+"?payment=success"),
		CancelURL:       env.GetString("STRIPE_CANCEL_URL", appURL+"?payment=cancel"),
	}

	if stripeCfg.StripeSecretKey == "" {
		log.Fatalf("STRIPE_SECRET_KEY is not set")
		return
	}

	paymentProcessor := stripe.NewStripeClient(stripeCfg)

	svc := service.NewPaymentService(paymentProcessor)
	


	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	log.Println("Starting RabbitMQ connection")

	// trip Consumer
	tripConsumer := events.NewTripConsumer(rabbitmq, svc)
	go tripConsumer.Listen()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutting down payment service...")
}
