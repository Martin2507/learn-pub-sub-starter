package pubsub

import (
	"bytes"
	"encoding/gob"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](conn *amqp.Connection, exchange, queueName, key string, simpleQueueType SimpleQueueType, handler func(T) Acktype) error {

	newChannel, queue, pubsubError := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)

	if pubsubError != nil {
		return pubsubError
	}

	deliveries, err := newChannel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {

		defer newChannel.Close()

		for msg := range deliveries {

			var temp T

			buffer := bytes.NewBuffer(msg.Body)

			dec := gob.NewDecoder(buffer)
			err := dec.Decode(&temp)

			if err != nil {
				log.Printf("Unable to decode incoming message: %s", err)
				msg.Nack(false, false)
				continue
			}

			result := handler(temp)

			switch result {

			case Ack:
				{
					msg.Ack(false)
				}

			case NackRequeue:
				{
					msg.Nack(false, true)
				}

			case NackDiscard:
				{
					msg.Nack(false, false)
				}

			default:
				{
					log.Printf("Unknown return type")
					continue
				}

			}

		}

	}()

	return nil

}
