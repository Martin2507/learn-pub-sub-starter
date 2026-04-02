package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {

	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(val)

	if err != nil {
		return err
	}

	channelError := ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{ContentType: "application/gob", Body: buffer.Bytes()})

	if channelError != nil {
		return channelError
	}

	return nil

}
