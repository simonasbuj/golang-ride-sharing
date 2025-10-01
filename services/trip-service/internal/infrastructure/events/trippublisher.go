package events

import (
	"context"
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

func (p *TripEventPublisher) PublishTripCreated(ctx context.Context) error {
	err := p.rabbitmq.PublishMessage(ctx, "hello", "my hardcoded message boi")
	
	return err
}
