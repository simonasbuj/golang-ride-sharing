package repository

import (
	"context"
	"fmt"
	"golang-ride-sharing/services/trip-service/internal/domain"
	"sync"

	pb "golang-ride-sharing/shared/proto/trip"
	pbd "golang-ride-sharing/shared/proto/driver"
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

func (r *inmemoryRepository) GetTripByID(ctx context.Context, id string) (*domain.TripModel, error) {
	trip, ok := r.trips[id]
	if !ok {
		return nil, fmt.Errorf("trip with id %s not found")
	}

	return trip, nil
}

func (r *inmemoryRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) (*domain.TripModel, error) {
	r.Lock()
	defer r.Unlock()

	trip, err := r.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, err
	}
	
	trip.Status = status

	if driver != nil {
		trip.Driver = &pb.TripDriver{
			Id: 			driver.Id,
			Name:			driver.Name,
			CarPlateNumber: driver.CarPlate,
			ProfilePicture: driver.ProfilePicture,
		}
	}

	return trip, nil
}
