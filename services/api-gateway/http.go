package main

import (
	"encoding/json"
	"fmt"
	"golang-ride-sharing/services/api-gateway/grpc_clients"
	"golang-ride-sharing/shared/contracts"
	"golang-ride-sharing/shared/env"
	"golang-ride-sharing/shared/messaging"
	"golang-ride-sharing/shared/tracing"
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)


var tracer = tracing.GetTracer("api-gateway")

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleTripPreview")
	defer span.End()

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
	tripServiceClient, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripServiceClient.Close()

	previewTripResponse, err := tripServiceClient.Client.PreviewTrip(ctx, reqBody.toProto())
	if err != nil {
		log.Printf("error in trip-service.PreviewTrip: %v", err)
		http.Error(w, "failed to preview a trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: previewTripResponse}
	writeJSON(w, http.StatusCreated, response)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleTripStart")
	defer span.End()

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

	tripServiceClient, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripServiceClient.Close()

	createTripResponse, err := tripServiceClient.Client.CreateTrip(ctx, reqBody.toProto())
		if err != nil {
		log.Printf("error in trip-service.CreateTrip: %v", err)
		http.Error(w, fmt.Sprintf("failed to create a trip: %s", err), http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: createTripResponse}
	writeJSON(w, http.StatusCreated, response)
}

func handelStripeWebhook(w http.ResponseWriter, r *http.Request, rabbitmq *messaging.RabbitMQ) {
	ctx, span := tracer.Start(r.Context(), "handleTripStart")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	webhookKey := env.GetString("STRIPE_WEBHOOK_KEY", "")
	if webhookKey == "" {
		log.Println("STRIPE_WEBHOOK_KEY not set and is required in stripe webhook")
		return 
	}

	event, err := webhook.ConstructEventWithOptions(
		body,
		r.Header.Get("Stripe-Signature"),
		webhookKey,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		log.Printf("error verifying webhook signature %v", err)
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	log.Printf("received stripe event: %v", event)

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession

		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		payload := messaging.PaymentStatusUpdateData{
			TripID:   session.Metadata["trip_id"],
			UserID:   session.Metadata["user_id"],
			DriverID: session.Metadata["driver_id"],
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshalling payload: %v", err)
			http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
			return
		}

		message := contracts.AmqpMessage{
			OwnerID: session.Metadata["user_id"],
			Data:    payloadBytes,
		}

		if err := rabbitmq.PublishMessage(
			ctx,
			contracts.PaymentEventSuccess,
			message,
		); err != nil {
			log.Printf("Error publishing payment event: %v", err)
			http.Error(w, "Failed to publish payment event", http.StatusInternalServerError)
			return
		}
	}
}
