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
	rabbitmq 	*messaging.RabbitMQ
	service 	*DriverService
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ, service *DriverService) *tripConsumer {
	return &tripConsumer{
		rabbitmq: 	rabbitmq,
		service:	service,
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

		switch msg.RoutingKey {
		case contracts.TripEventCreated, contracts.TripEventDriverNotInterested:
			return c.handleFindAndNotifyDrivers(ctx, payload)
		}

		log.Printf("unknown trip event: %+v", tripEvent)

		return nil
	})

	return err
}

func (c *tripConsumer) handleFindAndNotifyDrivers(ctx context.Context, payload messaging.TripEventData) error {
	suitableDriverIDs := c.service.FindAvailableDrivers(payload.Trip.SelectedFare.PackageSlug)

	log.Printf("found suitable drivers: %v", suitableDriverIDs)

	if len(suitableDriverIDs) == 0 {
		// Notify client that no drivers are avaialble
		if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventNoDriversFound, contracts.AmqpMessage{
			OwnerID: payload.Trip.Id,
		}); err != nil {
			log.Printf("ERROR: failed to publish message to '%s' exchange: %v", contracts.TripEventNoDriversFound, err)
			return err
		}
		return nil
	}

	suitableDriverID := suitableDriverIDs[0]

	// notify the driver about potential trip
	marshalledEvent, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := c.rabbitmq.PublishMessage(ctx, contracts.DriverCmdTripRequest, contracts.AmqpMessage{
		OwnerID: suitableDriverID,
		Data: marshalledEvent,
	}); err != nil {
		log.Printf("ERROR: failed to publish message to '%s' exchange: %v", contracts.DriverCmdTripRequest, err)
		return err
	}

	return nil
}
