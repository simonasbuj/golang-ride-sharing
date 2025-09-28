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
	
	// start http server
	log.Printf("starting HTTP server on port %s", httpAddr)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /preview", func(w http.ResponseWriter, r *http.Request) {
		http_handlers.HandleTripPreview(w, r, tripService)
	})

	server := &http.Server{
		Addr: 		httpAddr,
		Handler: 	mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server error: %v", err)
	}

}

