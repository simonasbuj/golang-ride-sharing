package events

import (
	"context"
	"encoding/json"
	"golang-ride-sharing/services/trip-service/internal/domain"
	"golang-ride-sharing/shared/contracts"
	"golang-ride-sharing/shared/messaging"
)


type TripEventPublisher struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripEventPublisher (rabbitmq *messaging.RabbitMQ) *TripEventPublisher {
	return &TripEventPublisher{
		rabbitmq: rabbitmq,
	}
}

func (p *TripEventPublisher) PublishTripCreated(ctx context.Context, trip *domain.TripModel) error {
	payload := messaging.TripEventData{
		Trip: trip.ToProto(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	message := contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data: body,
	}

	err = p.rabbitmq.PublishMessage(ctx, contracts.TripEventCreated, message)
	
	return err
}
