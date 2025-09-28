package service

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-ride-sharing/services/trip-service/internal/domain"
	"golang-ride-sharing/shared/types"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type tripService struct {
	repo domain.TripRepository
}

func NewTripService(repo domain.TripRepository) *tripService {
	return &tripService{
		repo: repo,
	}
}

func (s *tripService) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	trip := &domain.TripModel{
		ID: primitive.NewObjectID(),
		UserID: fare.UserID,
		Status: "pending",
		RideFare: fare,
	}
	return s.repo.CreateTrip(ctx, trip)
}

func (s *tripService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routes from project-osrm API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}

	var routeResp types.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %v", err)
	}

	return &routeResp, nil

}
