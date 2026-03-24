package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {

	newChannel, err := conn.Channel()

	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("Could not created channel: %v", err)
	}

	q, err := newChannel.QueueDeclare(queueName, queueType == SimpleQueueDurable, queueType == SimpleQueueTransient, queueType == SimpleQueueTransient, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("Could not create a new queue: %v", err)
	}

	bindError := newChannel.QueueBind(q.Name, key, exchange, false, nil)
	if bindError != nil {
		return nil, amqp.Queue{}, fmt.Errorf("Could not bind the queue to the new channel: %v", bindError)
	}

	return newChannel, q, nil
}
