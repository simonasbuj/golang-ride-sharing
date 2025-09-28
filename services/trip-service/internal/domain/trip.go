package domain

import (
	"context"

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
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
}
