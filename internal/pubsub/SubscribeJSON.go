package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) Acktype) error {

	newChannel, _, pubsubError := DeclareAndBind(conn, exchange, queueName, key, queueType)

	if pubsubError != nil {
		return pubsubError
	}

	deliveries, err := newChannel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {

		defer newChannel.Close()

		for msg := range deliveries {

			var target T
			if err := json.Unmarshal(msg.Body, &target); err != nil {
				continue
			}

			result := handler(target)

			switch result {

			case Ack:
				{
					log.Printf("Ack")
					msg.Ack(false)
				}

			case NackRequeue:
				{
					log.Printf("NackRequeue")
					msg.Nack(false, true)
				}

			case NackDiscard:
				{
					log.Printf("NackDiscard")
					msg.Nack(false, false)
				}

			default:
				{
					log.Printf("Unknown return type")
					return
				}
			}
		}
	}()

	return nil
}
