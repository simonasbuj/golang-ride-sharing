package domain

import (
	"context"
	"golang-ride-sharing/shared/types"
	trip_types "golang-ride-sharing/services/trip-service/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	pb "golang-ride-sharing/shared/proto/trip"
	pbd "golang-ride-sharing/shared/proto/driver"
)



type TripModel struct {
	ID 			primitive.ObjectID	`json:"ID"`
	UserID		string				`json:"userID"`
	Status 		string				`json:"status"`
	RideFare	*RideFareModel		`json:"rideFare"`
	Driver		*pb.TripDriver		`json:"driver"`
}

func (t *TripModel) ToProto() *pb.Trip {
	return &pb.Trip{
		Id: 			t.ID.Hex(),
		UserID: 		t.UserID,
		SelectedFare: 	t.RideFare.toProto(),
		Status: 		t.Status,
		Driver: 		t.Driver,
		Route: 			t.RideFare.Route.ToProto(),
	}
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, fare *RideFareModel) error
	GetRideFareByID(ctx context.Context, rideFareID string) (*RideFareModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*trip_types.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(ctx context.Context, route *trip_types.OsrmApiResponse) []*RideFareModel
	GenerateTripFares(ctx context.Context, fares []*RideFareModel, userID string, route *trip_types.OsrmApiResponse) ([]*RideFareModel, error)
	GetAndValidateRideFare(ctx context.Context, rideFareID, userID string) (*RideFareModel, error)
	GetTripByID(ctx context.Context, id string) (*TripModel, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) (*TripModel, error)
}
