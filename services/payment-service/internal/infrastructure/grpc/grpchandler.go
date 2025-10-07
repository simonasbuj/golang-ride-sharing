package grpc

import (
	pb "golang-ride-sharing/shared/proto/trip"

	"google.golang.org/grpc"
)


type grpcHandler struct {
	pb.UnimplementedTripServiceServer
}

func NewGrpcHandler(server *grpc.Server) *grpcHandler {
	handler := &grpcHandler{}

	pb.RegisterTripServiceServer(server, handler)
	return handler
}
