package messaging

import (
	"encoding/json"
	"golang-ride-sharing/shared/contracts"
	"log"
)

type QueueConsumer struct {
	rabbitmq  *RabbitMQ
	connMgr   *ConnectionManager
	queueName string
}

func NewQueueConsumer(rabbitmq *RabbitMQ, connMgr *ConnectionManager, queueName string) *QueueConsumer {
	return &QueueConsumer{
		rabbitmq:  rabbitmq,
		connMgr:   connMgr,
		queueName: queueName,
	}
}

func (qc *QueueConsumer) Start() error {
	msgs, err := qc.rabbitmq.Channel.Consume(
		qc.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var msgBody contracts.AmqpMessage
			if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
				log.Println("failed to unmarshal message: %+v, error: %v", msg.Body, err)
				continue
			}

			userID := msgBody.OwnerID

			var payload any
			if msgBody.Data != nil {
				if err := json.Unmarshal(msgBody.Data, &payload); err != nil {
					log.Println("failed to unmarshal payload: %+v, error: %v", msgBody.Data, err)
				}
			}

			clientMsg := contracts.WSMessage{
				Type: msg.RoutingKey,
				Data: payload,
			}

			if err := qc.connMgr.SendMessage(userID, clientMsg); err != nil {
				log.Printf("failed to send message to user: %s: %v", userID, err)
			}

		}
	}()

	return nil
}
