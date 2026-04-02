package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/rabbitmq/amqp091-go"
)

func handlerConsumeWarMessages(gs *gamelogic.GameState, ch *amqp091.Channel) func(rw gamelogic.RecognitionOfWar) pubsub.Acktype {

	return func(rw gamelogic.RecognitionOfWar) pubsub.Acktype {

		defer fmt.Print("> ")

		outcome, winner, loser := gs.HandleWar(rw)

		switch outcome {

		case gamelogic.WarOutcomeNotInvolved:
			{
				return pubsub.NackRequeue
			}

		case gamelogic.WarOutcomeNoUnits:
			{
				return pubsub.NackDiscard
			}

		case gamelogic.WarOutcomeOpponentWon:
			{
				formattedString := (winner + " won a war against " + loser)
				err := publishGameLog(ch, gs.GetUsername(), formattedString)

				if err != nil {
					return pubsub.NackRequeue
				}

				return pubsub.Ack
			}

		case gamelogic.WarOutcomeYouWon:
			{
				formattedString := (winner + " won a war against " + loser)
				err := publishGameLog(ch, gs.GetUsername(), formattedString)

				if err != nil {
					return pubsub.NackRequeue
				}

				return pubsub.Ack
			}

		case gamelogic.WarOutcomeDraw:
			{
				formattedString := ("A war between " + winner + " and " + loser + " resulted in a draw")
				err := publishGameLog(ch, gs.GetUsername(), formattedString)

				if err != nil {
					return pubsub.NackRequeue
				}

				return pubsub.Ack
			}

		default:
			{
				return pubsub.NackDiscard
			}
		}

	}

}
