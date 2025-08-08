package blackJackGame

import (
	"CombinedCardgames/blackJackGame/game" // Import game package
	"CombinedCardgames/blackJackGame/ui"   // Import ui package
	"fmt"
	"log"
	"os"

	"gioui.org/app" // Required for app.Window and app.Main
)

func BlackJackMain(window *app.Window) {
	// Initializing game logic first
	currentGame, err := game.NewGame()
	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	for {
		choice, err := ui.RunMenu(window, currentGame)

		switch choice {
		case "mainMenu":
			return

		case "exit":
			os.Exit(0)
			return

		case "Error":
			fmt.Println(err)
			return
		}
	}
}
