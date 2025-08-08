package baccaratGame

import (
	"CombinedCardgames/baccaratGame/game"
	"CombinedCardgames/baccaratGame/ui"
	"gioui.org/app"
	"log"
)

func BaccaratMain(window *app.Window) {
	currentGame, err := game.NewGame()
	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	for {
		// Show the Baccarat menu and wait for the user's choice
		choice := ui.RunBaccaratMenu(window, currentGame)

		switch choice {
		case "start":
			// If the user chooses to start, run the main game UI
			if err := ui.RunBaccaratUI(window, currentGame); err != nil {
				log.Fatal(err)
			}
		case "main_menu":
			// If the user chooses to return to the main menu, exit this function
			return
		case "exit":
			// If the user chooses to exit, terminate the application
			return
		}
	}
}
