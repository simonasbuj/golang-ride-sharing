package repository

import (
	"context"
	"golang-ride-sharing/services/trip-service/internal/domain"
)


type inmemoryRepository struct {
	trips 		map[string]*domain.TripModel
	rideFares 	map[string]*domain.RideFareModel
}

func NewInmemoryRepository() *inmemoryRepository {
	return &inmemoryRepository{
		trips: 		make(map[string]*domain.TripModel),
		rideFares: 	make(map[string]*domain.RideFareModel),
	}
}


func (r *inmemoryRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

