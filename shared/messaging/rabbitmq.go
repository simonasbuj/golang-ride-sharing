package messaging

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
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

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}

	if r.Channel != nil {
		r.Channel.Close()
	}
}

func (r *RabbitMQ) setupExchangesAndQueues() error {
	_, err := r.Channel.QueueDeclare(
		"hello", 	// name
		true, 		// durable
		false, 		// delete when used
		false, 		// exclusive
		false,		// no-wait
		nil,		// arguments
	)

	return err
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message string) error {
	err := r.Channel.PublishWithContext(
		ctx,
		"",						// exchange
		"hello",				// routing key aka queue name
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
