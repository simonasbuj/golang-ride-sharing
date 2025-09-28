package http

import (
	"encoding/json"
	"golang-ride-sharing/services/trip-service/internal/domain"
	"golang-ride-sharing/shared/types"
	"log"
	"net/http"
)

type httpHandler struct {
	svc domain.TripService
}

func NewHttpHandler(svc domain.TripService) *httpHandler {
	return &httpHandler{
		svc: svc,
	}
}

func (h *httpHandler) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}


	fare := &domain.RideFareModel{
		UserID: "42069",
	}

	trip, err := h.svc.CreateTrip(r.Context(), fare)
	if err != nil {
		log.Println(err)
	}

	log.Printf("new trip created: %+v", trip)
	writeJSON(w, http.StatusOK, trip)
}


type previewTripRequest struct {
	UserID 		string				`json:"userID"`
	Pickup 		types.Coordinate	`json:"pickup"`
	Destination	types.Coordinate	`json:"destination"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	return json.NewEncoder(w).Encode(data)
}