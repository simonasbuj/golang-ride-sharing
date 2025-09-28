package main

import (
	"golang-ride-sharing/shared/types"
	pb "golang-ride-sharing/shared/proto/trip"
)

type previewTripRequest struct {
	UserID 		string				`json:"userID"`
	Pickup 		types.Coordinate	`json:"pickup"`
	Destination	types.Coordinate	`json:"destination"`
}

func (p *previewTripRequest) toProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserID: p.UserID,
		StartLocation: &pb.Coordinate{Latitiude: p.Pickup.Latitude, Longitude: p.Pickup.Longitude},
		EndLocation: &pb.Coordinate{Latitiude: p.Destination.Latitude, Longitude: p.Destination.Longitude},
	}
}