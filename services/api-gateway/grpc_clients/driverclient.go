package grpc_clients

import (
	"golang-ride-sharing/shared/env"
	pb "golang-ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type driverServiceClient struct {
	Client 	pb.DriverServiceClient
	conn 	*grpc.ClientConn
}

func NewDriverServiceClient() (*driverServiceClient, error) {
	tripServiceURL := env.GetString("DRIVER_SERVICE_URL", "driver-service:9092")

	conn, err := grpc.NewClient(tripServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewDriverServiceClient(conn)
	
	return &driverServiceClient{
		Client: client,
		conn:	conn,
	}, nil
}

func (c *driverServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
