package grpc

import (
	"context"
	"golang-ride-sharing/services/trip-service/internal/domain"
	"golang-ride-sharing/services/trip-service/internal/infrastructure/events"
	pb "golang-ride-sharing/shared/proto/trip"
	"golang-ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type grpcHandler struct {
	pb.UnimplementedTripServiceServer

	service 	domain.TripService
	publisher 	*events.TripEventPublisher
}

func NewGrpcHandler(server *grpc.Server, service domain.TripService, publisher 	*events.TripEventPublisher) *grpcHandler {
	handler := &grpcHandler{
		service: 	service,
		publisher: 	publisher,
	}

	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *grpcHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	
	startLocation := req.GetStartLocation()
	endLocation := req.GetEndLocation()

	pickup := &types.Coordinate{
		Latitude: startLocation.Latitude,
		Longitude: startLocation.Longitude,
	}
	destination := &types.Coordinate{
		Latitude: endLocation.Latitude,
		Longitude: endLocation.Longitude,
	}

	route, err := h.service.GetRoute(ctx, pickup, destination)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	estimatedFares := h.service.EstimatePackagesPriceWithRoute(ctx, route)
	fares, err := h.service.GenerateTripFares(ctx, estimatedFares, req.GetUserID(), route)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate the ride fares: %v", err)
	}

	response := &pb.PreviewTripResponse{
		TripID: "fake-id-hardcoded",
		Route: route.ToProto(),
		RideFares: domain.ToRideFaresProto(fares),
	}
	return response, nil

}

func (h *grpcHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	rideFare, err := h.service.GetAndValidateRideFare(ctx, req.GetRideFareID(), req.GetUserID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get and validate fare: %v ", err)
	}

	trip, err := h.service.CreateTrip(ctx, rideFare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create a trip: %v", err)
	}

	if err := h.publisher.PublishTripCreated(ctx, trip); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish trip create event: %v", err)
	}

	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil
}
