package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {

	return func(am gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(am)

		if outcome == gamelogic.MoveOutcomeSafe {

			return pubsub.Ack
		}

		if outcome == gamelogic.MoveOutcomeMakeWar {
			err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				routing.WarRecognitionsPrefix+"."+am.Player.Username,
				gamelogic.RecognitionOfWar{
					Attacker: am.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)

			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return pubsub.NackRequeue
			}

			return pubsub.Ack
		}

		if outcome == gamelogic.MoveOutcomeSamePlayer {
			return pubsub.NackDiscard
		}

		return pubsub.NackDiscard
	}
}
