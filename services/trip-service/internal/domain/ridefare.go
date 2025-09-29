package domain

import (
	"time"

	pb "golang-ride-sharing/shared/proto/trip"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID	`json:"ID"`
	UserID            string				`json:"userID"`
	PackageSlug       string 				`json:"packageSlug"`	// ex: van, luxury, sedan
	TotalPriceInCents float64				`json:"totalPriceInCents"`
	ExpiresAt         time.Time				`json:"expiresAt"`
}

func (r *RideFareModel) toProto() *pb.RideFare {
	return &pb.RideFare{
		Id: r.ID.Hex(),
		UserID: r.UserID,
		PackageSlug: r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	var protoFares []*pb.RideFare

	for _, fare := range fares {
		protoFares = append(protoFares, fare.toProto())
	}

	return protoFares
}