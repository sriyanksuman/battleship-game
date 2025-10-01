package main

import (
	"battleship-game/service"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	game := service.NewGameService()

	// CLI REPL mode
	reader := bufio.NewScanner(os.Stdin)
	fmt.Println("Battleship CLI. Commands:")
	fmt.Println("  initGame N")
	fmt.Println("  addShip ID SIZE xA yA xB yB")
	fmt.Println("  viewBattleField")
	fmt.Println("  startGame")
	fmt.Println("  help | exit")
	for {
		fmt.Print(">> ")
		if !reader.Scan() {
			break
		}
		line := strings.TrimSpace(reader.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		cmd := strings.ToLower(parts[0])
		switch cmd {
		case "exit", "quit":
			return
		case "help":
			fmt.Println("Commands:")
			fmt.Println("  initGame N")
			fmt.Println("  addShip ID SIZE xA yA xB yB")
			fmt.Println("  viewBattleField")
			fmt.Println("  startGame")
			continue
		case "initgame":
			if len(parts) != 2 {
				fmt.Println("Usage: initGame N")
				continue
			}
			n, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Printf("Invalid N: %v\n", err)
				continue
			}
			if err := game.InitGameSize(n); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "addship":
			if len(parts) != 7 {
				fmt.Println("Usage: addShip ID SIZE xA yA xB yB")
				continue
			}
			id := parts[1]
			size, err1 := strconv.Atoi(parts[2])
			xA, err2 := strconv.Atoi(parts[3])
			yA, err3 := strconv.Atoi(parts[4])
			xB, err4 := strconv.Atoi(parts[5])
			yB, err5 := strconv.Atoi(parts[6])
			if err := firstErr(err1, err2, err3, err4, err5); err != nil {
				fmt.Printf("Invalid args: %v\n", err)
				continue
			}
			if err := game.AddShip(id, size, xA, yA, xB, yB); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "viewbattlefield":
			if err := game.ViewBattlefield(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "startgame":
			if err := game.Start(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		default:
			fmt.Println("Unknown command. Type 'help'.")
		}
	}
}

func firstErr(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}
