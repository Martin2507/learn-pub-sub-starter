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

	_, _, pubsubErr := pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTransient)

	if pubsubErr != nil {
		log.Fatalf("Could not declare and bind: %v", pubsubErr)
		return
	}

	state := gamelogic.NewGameState(username)

	if state == nil {
		log.Fatalf("Could not get a new state of the game")
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
				_, moveError := state.CommandMove(userInput)
				if moveError != nil {
					fmt.Println(moveError)
					continue
				}
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
