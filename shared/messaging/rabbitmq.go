package messaging

import (
	"context"
	"fmt"
	"golang-ride-sharing/shared/contracts"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)


const (
	TripExchange = "trip"
)


type RabbitMQ struct {
	conn 	*amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	log.Print("starting RabbitMQ connection")

	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	log.Print("openning RabbitMQ channel")
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	rmq := &RabbitMQ{
		conn: 		conn,
		Channel: 	ch,
	}

	if err := rmq.setupExchangesAndQueues(); err != nil {
		rmq.Close()
		return nil, err
	}

	return rmq, nil
}

func (r *RabbitMQ) declareAndBoundQueue(queueName string, messageTpes []string, exchangeName string) error {
	q, err := r.Channel.QueueDeclare(
		queueName, 	// queue name
		true, 		// durable
		false, 		// delete when used
		false, 		// exclusive
		false,		// no-wait
		nil,		// arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s", queueName)
	}

	for _, msg := range messageTpes {
		err = r.Channel.QueueBind(
			q.Name,				// queue name
			msg, 				// routing key
			exchangeName,		// exchange name
			false,				// no-wait
			nil,				// arguments
		)
		if err != nil {
			return fmt.Errorf("failed to bind quueue %s: %v", queueName, err)
		}
	}

	return  nil
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}

	if r.Channel != nil {
		r.Channel.Close()
	}
}

func (r *RabbitMQ) setupExchangesAndQueues() error {
	err := r.Channel.ExchangeDeclare(
		TripExchange, 	// exhcange name
		"topic",		// routing type
		true,			// durable
		false,			// auto-deleted
		false,			//internal
		false,			// no-wait
		nil,			// arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange %s: %v", TripExchange, err)
	}

	if err := r.declareAndBoundQueue(
		"find_available_drivers",
		[]string{
			contracts.TripEventCreated,
			contracts.TripEventDriverNotInterested,
		},
		TripExchange,
	); err != nil { return err }

	return err
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message string) error {
	log.Printf("publishing message with routingKey: %s", routingKey)
	err := r.Channel.PublishWithContext(
		ctx,
		TripExchange,			// exchange
		routingKey,				// routing key aka queue name
		false,					// mandatory
		false,					// immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
			DeliveryMode: amqp.Persistent,
		},
	)

	return err
}

type MessageHandler func(context.Context, amqp.Delivery) error

func (r *RabbitMQ) ConsumeMessages(queueName string, handler MessageHandler) error {
	// set qos to fair dispatch
	err := r.Channel.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	msgs, err := r.Channel.Consume(
		queueName,
		"",			// consumer
		false,		// auto-ack
		false,		// exclusive
		false,		// no-local
		false,		// no-wait
		nil,		// args

	)
	if err != nil {
		return err
	}

	ctx := context.Background()

	go func() {
		for msg := range msgs {
			if err := handler(ctx, msg); err != nil {
				log.Printf("ERROR: failed to handle the message, message: %s, error: %v", msg.Body, err)
				if nackErr := msg.Nack(false, false); nackErr != nil {
					log.Printf("ERRPR: failed to Nack message: %v", nackErr)
				}

				continue
			}

			// acknowledge msg if handler succeded otherwise gonna stay in unack
			if ackErr := msg.Ack(false); ackErr != nil {
				log.Printf("ERROR: failed to Ack message: %. Message body: %s", ackErr, msg.Body)
			}
		}
	}()

	return nil
}
