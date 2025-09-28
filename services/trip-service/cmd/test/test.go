package main

import (
	"context"
	"golang-ride-sharing/services/trip-service/internal/infrastructure/repository"
	"golang-ride-sharing/services/trip-service/internal/service"
	"golang-ride-sharing/shared/types"
	"log"
)


func main() {
	// dependency injections
	inmemoryRepo := repository.NewInmemoryRepository()
	tripService := service.NewTripService(inmemoryRepo)
	
	route, err := tripService.GetRoute(
		context.Background(), 
		&types.Coordinate{Latitude: 54.84695790910313, Longitude: 25.472619082527498,},
		&types.Coordinate{Latitude: 54.846111615547315, Longitude: 25.470966841844525,},
	)
	if err != nil {
		log.Println(err)
	}
	log.Print(route)
}

