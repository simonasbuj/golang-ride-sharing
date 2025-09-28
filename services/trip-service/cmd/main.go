package main

import (
	http_handlers "golang-ride-sharing/services/trip-service/internal/infrastructure/http"
	"golang-ride-sharing/services/trip-service/internal/infrastructure/repository"
	"golang-ride-sharing/services/trip-service/internal/service"
	"golang-ride-sharing/shared/env"
	"log"
	"net/http"
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
		Addr: 		httpAddr,
		Handler: 	mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server error: %v", err)
	}

}

