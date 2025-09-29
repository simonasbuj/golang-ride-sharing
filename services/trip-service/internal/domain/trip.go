package domain

import (
	"context"
	"golang-ride-sharing/shared/types"
	trip_types "golang-ride-sharing/services/trip-service/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type TripModel struct {
	ID 			primitive.ObjectID	`json:"ID"`
	UserID		string				`json:"userID"`
	Status 		string				`json:"status"`
	RideFare	*RideFareModel		`json:"rideFare"`
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, fare *RideFareModel) error
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*trip_types.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(ctx context.Context, route *trip_types.OsrmApiResponse) []*RideFareModel
	GenerateTripFares(ctx context.Context, fares []*RideFareModel, userID string) ([]*RideFareModel, error)
}
