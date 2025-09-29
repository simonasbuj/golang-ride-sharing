package grpc

import (
	"context"
	"golang-ride-sharing/services/trip-service/internal/domain"
	pb "golang-ride-sharing/shared/proto/trip"
	"golang-ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type grpcHandler struct {
	pb.UnimplementedTripServiceServer

	service domain.TripService
}

func NewGrpcHandler(server *grpc.Server, service domain.TripService) *grpcHandler {
	handler := &grpcHandler{
		service: service,
	}

	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *grpcHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	
	startLocation := req.GetStartLocation()
	endLocation := req.GetEndLocation()

	pickup := &types.Coordinate{
		Latitude: startLocation.Latitiude,
		Longitude: startLocation.Longitude,
	}
	destination := &types.Coordinate{
		Latitude: endLocation.Latitiude,
		Longitude: endLocation.Longitude,
	}

	trip, err := h.service.GetRoute(ctx, pickup, destination)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	response := &pb.PreviewTripResponse{
		TripID: "fake-id-hardcoded",
		Route: trip.ToProto(),
		RideFares: []*pb.RideFare{},
	}
	return response, nil

}