package repository

import (
	"context"
	"fmt"
	"golang-ride-sharing/services/trip-service/internal/domain"
	"sync"
)


type inmemoryRepository struct {
	trips 		map[string]*domain.TripModel
	rideFares 	map[string]*domain.RideFareModel

	sync.Mutex
}

func NewInmemoryRepository() *inmemoryRepository {
	return &inmemoryRepository{
		trips: 		make(map[string]*domain.TripModel),
		rideFares: 	make(map[string]*domain.RideFareModel),
	}
}


func (r *inmemoryRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.Lock()
	defer r.Unlock()

	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (r *inmemoryRepository) SaveRideFare(ctx context.Context, fare *domain.RideFareModel) error {
	r.Lock()
	defer r.Unlock()

	r.rideFares[fare.ID.Hex()] = fare
	return nil
}

func (r *inmemoryRepository) GetRideFareByID(ctx context.Context, rideFareID string) (*domain.RideFareModel, error) {
	rideFare, ok := r.rideFares[rideFareID]
	if !ok {
		return nil, fmt.Errorf("ride fare with ID %s not found", rideFareID)
	}
	
	return rideFare, nil
}
