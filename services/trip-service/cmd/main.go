package main

import (
	"context"
	http_handlers "golang-ride-sharing/services/trip-service/internal/infrastructure/http"
	"golang-ride-sharing/services/trip-service/internal/infrastructure/repository"
	"golang-ride-sharing/services/trip-service/internal/service"
	"golang-ride-sharing/shared/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8083")
)

func main() {
	// dependency injections
	inmemoryRepo := repository.NewInmemoryRepository()
	tripService := service.NewTripService(inmemoryRepo)

	httpHandler := http_handlers.NewHttpHandler(tripService)

	// start http server
	log.Printf("starting HTTP server on port %s", httpAddr)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /preview", httpHandler.HandleTripPreview)

	server := &http.Server{
		Addr:    httpAddr,
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
	case err := <-serverErrors:
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
