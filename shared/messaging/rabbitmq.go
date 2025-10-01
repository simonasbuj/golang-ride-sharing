package messaging

import (
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn *amqp.Connection
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	log.Print("starting RabbitMQ connection")

	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Conn: conn,
	}, nil
}

func (r *RabbitMQ) Close() {
	if r.Conn != nil {
		r.Conn.Close()
	}
}
