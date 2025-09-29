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
		log.Printf("failed to preview a trip: %v", err)
		http.Error(w, "failed to parse JSON payload, userID is required", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: previewTripResponse}
	writeJSON(w, http.StatusCreated, response)
}
