package main

import (
	"encoding/json"
	"golang-ride-sharing/services/api-gateway/grpc_clients"
	"golang-ride-sharing/shared/contracts"
	"log"
	"net/http"
)


func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// payload validation
	if reqBody.UserID == "" {
		http.Error(w, "failed to parse JSON payload, userID is required", http.StatusBadRequest)
		return
	}

	// call TripService
	tripServiceClient, err := grpc_clients.NewtripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripServiceClient.Close()

	previewTripResponse, err := tripServiceClient.Client.PreviewTrip(r.Context(), reqBody.toProto())
	if err != nil {
		log.Printf("error in trip-service.PreviewTrip: %v", err)
		http.Error(w, "failed to preview a trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: previewTripResponse}
	writeJSON(w, http.StatusCreated, response)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	var reqBody startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// payload validation
	if reqBody.UserID == "" || reqBody.RideFareID == "" {
		http.Error(w, "failed to parse JSON payload, userID and rideFareID are required", http.StatusBadRequest)
		return
	}

	tripServiceClient, err := grpc_clients.NewtripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripServiceClient.Close()

	createTripResponse, err := tripServiceClient.Client.CreateTrip(r.Context(), reqBody.toProto())
		if err != nil {
		log.Printf("error in trip-service.CreateTrip: %v", err)
		http.Error(w, "failed to create a trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: createTripResponse}
	writeJSON(w, http.StatusCreated, response)
}
