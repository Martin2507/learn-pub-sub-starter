package pubsub

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T)) error {

	_, _, pubsubError := DeclareAndBind(conn, exchange, queueName, key, queueType)

	if pubsubError != nil {
		return pubsubError
	}

	newChannel, err := conn.Channel()
	if err != nil {
		return err
	}

	deliveries, err := newChannel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return nil
	}

	go func() {

		for msg := range deliveries {

			var target T
			if err := json.Unmarshal(msg.Body, &target); err != nil {
				continue
			}

			handler(target)

			msg.Ack(false)
		}
	}()

	return nil
}
