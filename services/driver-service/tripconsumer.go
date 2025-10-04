package main

import (
	"context"
	"encoding/json"
	"golang-ride-sharing/shared/contracts"
	"golang-ride-sharing/shared/messaging"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type tripConsumer struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ) *tripConsumer {
	return &tripConsumer{
		rabbitmq: rabbitmq,
	}
}

func (c *tripConsumer) Listen() error {
	err := c.rabbitmq.ConsumeMessages(messaging.FindAvailableDriversQueue, func(ctx context.Context, msg amqp.Delivery) error {
		var tripEvent contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
			log.Printf("ERROR: failed to unmarshal message: %v, error: %v", msg, err)
			return err
		}

		var payload messaging.TripEventData
		if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
			log.Printf("ERROR: failed to unmarshal message payload: %v, error: %v", tripEvent.Data, err)
			return err
		}

		log.Printf("driver received a message: %+v", payload)
		return nil
	})

	return err
}
