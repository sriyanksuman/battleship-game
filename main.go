package main

import (
	"battleship-game/service"
	"fmt"
	"os"
)

func main() {
	inputFilePath := "input.json"

	if len(os.Args) > 1 {
		inputFilePath = os.Args[1]
	}

	game := service.NewGameService()

	err := game.InitGame(inputFilePath)
	if err != nil {
		fmt.Printf("Error initializing game: %v\n", err)
		os.Exit(1)
	}

	game.PlayGame()
}
