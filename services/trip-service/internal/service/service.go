package service

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-ride-sharing/services/trip-service/internal/domain"
	trip_types "golang-ride-sharing/services/trip-service/internal/types"
	"golang-ride-sharing/shared/proto/trip"
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
		Driver: &trip.TripDriver{},
	}
	return s.repo.CreateTrip(ctx, trip)
}

func (s *tripService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*trip_types.OsrmApiResponse, error) {
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

	var routeResp trip_types.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %v", err)
	}

	return &routeResp, nil

}

func (s *tripService) EstimatePackagesPriceWithRoute(ctx context.Context, route *trip_types.OsrmApiResponse) []*domain.RideFareModel {
	baseFares := getBaseFares()

	estimatedFares := make([]*domain.RideFareModel, len(baseFares))

	for i, f := range baseFares {
		estimatedFares[i] = s.estimateFareRoute(f, route)
	}
	return estimatedFares
}

func (s *tripService) GenerateTripFares(ctx context.Context, rideFares []*domain.RideFareModel, userID string, route *trip_types.OsrmApiResponse) ([]*domain.RideFareModel, error){
	fares := make([]*domain.RideFareModel, len(rideFares))

	for i, f := range rideFares {
		id := primitive.NewObjectID()

		fare := &domain.RideFareModel{
			ID: id,
			UserID: userID,
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug: f.PackageSlug,
			Route: route,
		}

		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("failed to save trip fare %s", err)
		}

		fares[i] = fare
	}

	return fares, nil
}

func (s *tripService) estimateFareRoute(f *domain.RideFareModel, route *trip_types.OsrmApiResponse) *domain.RideFareModel {
	pricingCfg := trip_types.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents

	distanceKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	distanceFare := distanceKm * pricingCfg.PricePerKm
	timeFare := durationInMinutes * pricingCfg.PricingPerMinute

	totalPrice := carPackagePrice + distanceFare + timeFare

	return &domain.RideFareModel{
		TotalPriceInCents: totalPrice,
		PackageSlug: f.PackageSlug,
	}
}

func (s *tripService) GetAndValidateRideFare(ctx context.Context, rideFareID, userID string) (*domain.RideFareModel, error) {
	rideFare, err := s.repo.GetRideFareByID(ctx, rideFareID)
	if err != nil {
		return nil, err
	}

	if rideFare.UserID != userID {
		return nil, fmt.Errorf("fare with id %s doesn't belong to userID %s", rideFareID, userID)
	}

	return rideFare, nil
}


func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug: "suv",
			TotalPriceInCents: 200,
		},
		{
			PackageSlug: "sedan",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug: "van",
			TotalPriceInCents: 400,
		},
		{
			PackageSlug: "luxury",
			TotalPriceInCents: 1000,
		},
	}
}


