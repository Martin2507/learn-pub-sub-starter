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

	// Attempt to connect to the RabbitMQ instance
	connection, err := amqp.Dial(rabbitConnString)

	// Return and display an error message in case the connection was unsuccessful
	if err != nil {
		log.Fatalf("Unable to establish the connection")
		return
	}

	// Close the connection to RabbitMQ after the server exits
	defer connection.Close()

	// Create a new channel on the successful connection
	newChannel, err := connection.Channel()

	// Return and display an error message in case there was an issue creating a new channel
	if err != nil {
		log.Fatalf("Unable to create a new channel")
		return
	}

	// Close the channel when the server exits
	defer newChannel.Close()

	fmt.Println("A connection was successfully established")

	fmt.Println("Starting Peril server...")

	_, _, pubsubErr := pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable)

	if pubsubErr != nil {
		log.Fatalf("Could not declare and bind: %v", pubsubErr)
		return
	}

	// Print available server commands to the console
	gamelogic.PrintServerHelp()

	// Start the REPL loop to process user input
	for {
		// Wait for the user to enter a command
		userInput := gamelogic.GetInput()

		// Skip empty input and wait for the next command
		if len(userInput) == 0 {
			continue
		}

		// Handle the first word of the input as the command
		switch userInput[0] {

		case "pause":
			// Notify the console and publish a pause message to the broker
			fmt.Println("The game is paused")

			pubsubError := pubsub.PublishJSON(
				newChannel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)

			// Log and exit if publishing the pause message fails
			if pubsubError != nil {
				log.Fatalf(pubsubError.Error())
				return
			}

		case "resume":
			// Notify the console and publish a resume message to the broker
			fmt.Println("The game is resuming")

			pubsubError := pubsub.PublishJSON(
				newChannel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)

			// Log and exit if publishing the resume message fails
			if pubsubError != nil {
				log.Fatalf(pubsubError.Error())
				return
			}

		case "quit":
			// Notify the console and exit the REPL loop
			fmt.Println("Exiting game")
			return

		default:
			// Notify the console that the command is not recognized
			fmt.Println("Unknown command")
		}
	}
}
