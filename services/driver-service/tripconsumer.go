package main

import (
	"context"
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
	err := c.rabbitmq.ConsumeMessages("hello", func(ctx context.Context, msg amqp.Delivery) error {
		log.Printf("driver received a message: %s", msg.Body)
		return nil
	})

	return err
}