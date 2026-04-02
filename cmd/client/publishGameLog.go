package main

import (
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func publishGameLog(ch *amqp.Channel, username string, str string) error {

	log := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     str,
		Username:    username,
	}

	err := pubsub.PublishGob(
		ch,
		routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+username,
		log)

	if err != nil {
		return err
	}

	return nil
}
