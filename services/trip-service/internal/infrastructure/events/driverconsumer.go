package events

import (
	"context"
	"encoding/json"
	"log"

	"golang-ride-sharing/services/trip-service/internal/domain"
	"golang-ride-sharing/shared/contracts"
	"golang-ride-sharing/shared/messaging"

	amqp "github.com/rabbitmq/amqp091-go"
)


type driverResponseConsumer struct {
	rabbitmq 	*messaging.RabbitMQ
	service 	domain.TripService
}

func NewDriverResponseConsumer(rabbitmq *messaging.RabbitMQ, service domain.TripService) *driverResponseConsumer {
	return &driverResponseConsumer{
		rabbitmq: 	rabbitmq,
		service: 	service,
	}
}

func (c *driverResponseConsumer) Listen() error  {
	err := c.rabbitmq.ConsumeMessages(messaging.DriverTripResponseQueue, func(ctx context.Context, msg amqp.Delivery) error {
		var driverResponseMsg contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &driverResponseMsg); err != nil {
			log.Printf("ERROR: failed to unmarshal message: %v, error: %v", msg, err)
			return err
		}

		var payload messaging.DriverTripResponseData
		if err := json.Unmarshal(driverResponseMsg.Data, &payload); err != nil {
			log.Printf("ERROR: failed to unmarshal message payload: %v, error: %v", driverResponseMsg.Data, err)
			return err
		}

		log.Printf("rider received a driver response message: %+v", payload)

		switch msg.RoutingKey {
		case contracts.DriverCmdTripAccept:
			if err := c.handelTripAccepted(ctx, payload); err != nil {
				log.Printf("ERROR: failed to handle driver trip accepted event: %+v", err)
				return err
			}
		case contracts.DriverCmdTripDecline:
			if err := c.handelTripDeclined(ctx, payload); err != nil {
				log.Printf("ERROR: failed to handle driver trip declined event: %v", err)
				return err
			}
		}

		log.Printf("unknown driver response event type: %s, event data: %+v", msg.RoutingKey, driverResponseMsg)

		return nil
	})

	return err
}

func (c *driverResponseConsumer) handelTripAccepted(ctx context.Context, payload messaging.DriverTripResponseData) error {
	log.Printf("handling trip accepted event")



	return nil
}

func (c *driverResponseConsumer) handelTripDeclined(ctx context.Context, payload messaging.DriverTripResponseData) error {
	log.Printf("handling trip declined event")

	_, err := c.service.GetTripByID(ctx, payload.TripID)
	if err != nil {
		return err
	}

	updatedTrip, err := c.service.UpdateTrip(ctx, payload.TripID, "accepted", payload.Driver)
	if err != nil {
		log.Printf("ERROR: failed to update the trip: %v", err)
		return err
	}

	marshalledTrip, err := json.Marshal(updatedTrip)
	if err != nil {
		return err
	}

	// notify the rider that a driver has been assigned
	if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventDriverAssigned, contracts.AmqpMessage{
		OwnerID: updatedTrip.UserID,
		Data: marshalledTrip,
	}); err != nil {
		return err
	}

	return nil
}