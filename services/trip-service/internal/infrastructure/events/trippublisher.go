package events

import (
	"context"
	"encoding/json"
	"golang-ride-sharing/services/trip-service/internal/domain"
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
	
	body, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	err = p.rabbitmq.PublishMessage(ctx, "hello", string(body))
	
	return err
}
