package grpc_clients

import (
	"golang-ride-sharing/shared/env"
	pb "golang-ride-sharing/shared/proto/trip"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type tripServiceClient struct {
	Client 	pb.TripServiceClient
	conn 	*grpc.ClientConn
}

func NewtripServiceClient() (*tripServiceClient, error) {
	tripServiceURL := env.GetString("TRIP_SERVICE_URL", "trip-service:9093")

	conn, err := grpc.NewClient(tripServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewTripServiceClient(conn)
	
	return &tripServiceClient{
		Client: client,
		conn:	conn,
	}, nil
}

func (c *tripServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
