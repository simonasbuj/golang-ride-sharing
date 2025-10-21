package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-ride-sharing/shared/env"
	"golang-ride-sharing/shared/messaging"
	"golang-ride-sharing/shared/tracing"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting API Gateway v2")

	// env vars
	rabbitMqUri := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")

	// initialzie tracing
	tracerCfg := tracing.Config{
		ServiceName: 	"api-gateway",
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

	// init rabbitmq
	rabbitmq, err := messaging.NewRabbitMQ(rabbitMqUri)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
		return
	}
	defer rabbitmq.Close()

	mux := http.NewServeMux()

	mux.Handle("POST /trip/preview",  tracing.WrapHandlerFunc(enableCORS(handleTripPreview), "/trip/preview"))
	mux.Handle("POST /trip/start",  tracing.WrapHandlerFunc(enableCORS(handleTripStart), "/trip/start"))
	mux.Handle("/ws/drivers", tracing.WrapHandlerFunc(func(w http.ResponseWriter, r *http.Request){ handleDriversWebSocket(w, r, rabbitmq) }, "/ws/drivers"))
	mux.Handle("/ws/riders", tracing.WrapHandlerFunc(func(w http.ResponseWriter, r *http.Request){ handleRidersWebSocket(w, r, rabbitmq) } , "/ws/riders"))
	mux.Handle("/webhook/stripe",  tracing.WrapHandlerFunc(func(w http.ResponseWriter, r *http.Request){ handelStripeWebhook(w, r, rabbitmq) }, "/webhook/stripe"))

	server := &http.Server{
		Addr: httpAddr,
		Handler: mux,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("server listening on %s", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err:= <-serverErrors:
		log.Printf("error starting the server: %v", err)

	case sig := <-shutdown:
		log.Printf("server is shutting down due to %v signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("could not stop server gracefully: %v", err)
			server.Close()
		}
	}
}
