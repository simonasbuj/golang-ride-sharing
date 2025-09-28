package main

import (
	"context"
	"golang-ride-sharing/services/trip-service/internal/domain"
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
	// start http server
	log.Printf("starting HTTP server on port %s", httpAddr)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /preview", handleTripPreview)

	server := &http.Server{
		Addr: 		httpAddr,
		Handler: 	mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server error: %v", err)
	}

}

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	inmemoryRepo := repository.NewInmemoryRepository()
	tripService := service.NewTripService(inmemoryRepo)

	fare := &domain.RideFareModel{
		UserID: "42069",
	}

	trip, err := tripService.CreateTrip(context.Background(), fare)
	if err != nil {
		log.Println(err)
	}

	log.Printf("new trip created: %+v", trip)
}