package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Rabbit Connection String assigned to a variable
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	// Attempt to connect to the rabbitmq instance
	connection, err := amqp.Dial(rabbitConnString)

	// Return and display an error message in case the conection was unsuccesfull
	if err != nil {
		log.Fatalf("Unable to establish the connection")
		return
	}

	// Close the connection to the rabbitmq after action has been carried out
	defer connection.Close()

	fmt.Println("A connection was succesfully established")
	fmt.Println("Starting Peril client...")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Unable to get username from the user")
		return
	}

	// _, _, pubsubErr := pubsub.DeclareAndBind(
	// 	connection,
	// 	routing.ExchangePerilDirect,
	// 	routing.PauseKey+"."+username,
	// 	routing.PauseKey,
	// 	pubsub.SimpleQueueTransient)

	// if pubsubErr != nil {
	// 	log.Fatalf("Could not declare and bind: %v", pubsubErr)
	// 	return
	// }

	state := gamelogic.NewGameState(username)

	if state == nil {
		log.Fatalf("Could not get a new state of the game")
		return
	}

	pausePubSubError := pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(state))

	if pausePubSubError != nil {
		log.Fatalf("Could not declare and bind: %v", pausePubSubError)
		return
	}

	newChannel, err := connection.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
		return
	}

	movePubSubError := pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+username,
		routing.ArmyMovesPrefix+".*",
		pubsub.SimpleQueueTransient,
		handlerMove(state, newChannel))

	if movePubSubError != nil {
		log.Fatalf("Could not declare and bind: %v", movePubSubError)
		return
	}

	warPubSubError := pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		pubsub.SimpleQueueDurable,
		handlerConsumeWarMessages(state, newChannel))

	if warPubSubError != nil {
		log.Fatalf("Could not declare and bind: %v", warPubSubError)
		return
	}

	publishPubSubError := pubsub.SubscribeGob(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
		handlerLogs())

	if publishPubSubError != nil {
		log.Fatalf("Could not declare and bind: %v", publishPubSubError)
		return
	}

	for {
		userInput := gamelogic.GetInput()

		if len(userInput) == 0 {
			continue
		}

		switch userInput[0] {

		case "spawn":
			{
				spawnError := state.CommandSpawn(userInput)
				if spawnError != nil {
					fmt.Println(spawnError)
					continue
				}
			}

		case "move":
			{
				move, moveError := state.CommandMove(userInput)
				if moveError != nil {
					fmt.Println(moveError)
					continue
				}

				publishError := pubsub.PublishJSON(
					newChannel,
					routing.ExchangePerilTopic,
					routing.ArmyMovesPrefix+"."+username,
					move)

				if publishError != nil {
					fmt.Println(publishError)
					continue
				}

				fmt.Println("You have made your move")
			}

		case "status":
			{
				state.CommandStatus()
			}

		case "help":
			{
				gamelogic.PrintClientHelp()
			}

		case "spam":
			{
				fmt.Println("Spamming not allowed yet!")
			}

		case "quit":
			{
				gamelogic.PrintQuit()
				return
			}

		default:
			{
				fmt.Println("Unknown command")
			}
		}
	}
}
