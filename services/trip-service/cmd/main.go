package main

import (
	"context"
	"golang-ride-sharing/services/trip-service/internal/domain"
	"golang-ride-sharing/services/trip-service/internal/infrastructure/repository"
	"golang-ride-sharing/services/trip-service/internal/service"
	"log"
	"time"
)

func main() {
	inmemoryRepo := repository.NewInmemoryRepository()
	tripService := service.NewTripService(inmemoryRepo)

	fare := &domain.RideFareModel{
		UserID: "42069",
	}
	trip, err := tripService.CreateTrip(context.Background(), fare)
	if err != nil {
		log.Println(err)
	}

	log.Println(trip)

	// TODO: delete this abomination
	for {
		time.Sleep(time.Second * 5)
		log.Println(time.Now())
	}
}
