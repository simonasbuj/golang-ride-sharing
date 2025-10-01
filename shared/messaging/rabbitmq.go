package messaging

import (
	"context"
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
		false, 		// durable
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
		},
	)

	return err
}
