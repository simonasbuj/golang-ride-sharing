package main

import "golang-ride-sharing/shared/types"

type previewTripRequest struct {
	UserID 		string				`json:"userID"`
	Pickup 		types.Coordinate	`json:"pickup"`
	Destination	types.Coordinate	`json:"destination"`
}