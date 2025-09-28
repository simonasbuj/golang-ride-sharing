package http

import (
	"context"
	"golang-ride-sharing/services/trip-service/internal/domain"
	"log"
	"net/http"
)

func HandleTripPreview(w http.ResponseWriter, r *http.Request, tripService domain.TripService) {
	fare := &domain.RideFareModel{
		UserID: "42069",
	}

	trip, err := tripService.CreateTrip(context.Background(), fare)
	if err != nil {
		log.Println(err)
	}

	log.Printf("new trip created: %+v", trip)
}