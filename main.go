package main

import (
	"battleship-game/service"
	"fmt"
	"os"
)

func main() {
	game := service.NewGameService()
	if len(os.Args) == 1 {
		// Interactive-style sample aligning to required APIs
		if err := game.InitGameSize(6); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := game.AddShip("SH1", 2, 1, 1, 4, 4); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := game.Start(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	inputFilePath := os.Args[1]
	if err := game.InitGame(inputFilePath); err != nil {
		fmt.Printf("Error initializing game: %v\n", err)
		os.Exit(1)
	}
	game.PlayGame()
}
