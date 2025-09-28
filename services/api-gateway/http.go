package main

import (
	"encoding/json"
	"golang-ride-sharing/shared/contracts"
	"net/http"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	_, _ = http.Get(
		"http://trip-service:8083/preview",
	)


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

	// TODO: call TripService
	response := contracts.APIResponse{Data: "all gucci"}
	writeJSON(w, http.StatusCreated, response)
}
