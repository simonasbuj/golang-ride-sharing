package main

import (
	"bytes"
	"log"
	"encoding/json"
	"golang-ride-sharing/shared/contracts"
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
	jsonBody, _ := json.Marshal(reqBody)
	reader := bytes.NewReader(jsonBody)

	svcResponse, err := http.Post(
		"http://trip-service:8083/preview",
		"application/json",
		reader,
	)
	if err != nil {
		log.Println(err)
		return
	}

	defer svcResponse.Body.Close()

	var svcResponseBody any
	if err := json.NewDecoder(svcResponse.Body).Decode(&svcResponseBody); err != nil {
		log.Println(err)
		http.Error(w, "failed to parse JSON response from trip service", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: svcResponseBody}
	writeJSON(w, http.StatusCreated, response)
}
