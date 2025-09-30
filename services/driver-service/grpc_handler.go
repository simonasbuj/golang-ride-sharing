package main

import (
	"context"
	pb "golang-ride-sharing/shared/proto/driver"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	pb.UnimplementedDriverServiceServer

	service *DriverService
}

func NewGrpcHandler(server *grpc.Server, service *DriverService) *grpcHandler {
	handler := &grpcHandler{
		service: service,
	}

	pb.RegisterDriverServiceServer(server, handler)
	return handler
}

func (h *grpcHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	log.Printf("registering driver with id: %s, packageSlug: %s", req.DriverID, req.PackageSlug)
	driver , err := h.service.RegisterDriver(req.DriverID, req.PackageSlug)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error in DriverService.RegisterDriver: %s", err)
	}

	return &pb.RegisterDriverResponse{
		Driver: driver,
	}, nil
}

func (h *grpcHandler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	log.Printf("unregistering driver with id: %s", req.DriverID)
	h.service.UnregisterDriver(req.DriverID)
	
	return &pb.RegisterDriverResponse{}, nil
}
