package types

import (
	pb "golang-ride-sharing/shared/proto/trip"
)


type OsrmApiResponse struct {
	Routes []struct{
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func (o *OsrmApiResponse) ToProto() *pb.Route {
	route := o.Routes[0]
	geometry := route.Geometry.Coordinates
	coordinates := make([]*pb.Coordinate, len(geometry))

	for i, coord := range geometry {
		coordinates[i] = &pb.Coordinate{
			Latitude: coord[0],
			Longitude: coord[1],
		}
	}

	return &pb.Route{
		Distance: route.Distance,
		Duration: route.Duration,
		Geometry: []*pb.Geometry{
			{
				Coordinates: coordinates,
			},
		},
	}
}

type PricingConfig struct {
	PricePerKm 			float64
	PricingPerMinute	float64
}

func DefaultPricingConfig() *PricingConfig {
	return &PricingConfig{
		PricePerKm: 		1.5,
		PricingPerMinute: 	0.25,
	}
}
