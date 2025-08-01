package main

import (
	"log"
	"os"

	"CardGame/game" // Import game package
	"CardGame/ui"   // Import ui package
	"gioui.org/app" // Required for app.Window and app.Main
)

func main() {
	// Initializing game logic first
	currentGame, err := game.NewGame()
	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	// Starting the Gio GUI in a separate goroutine
	go func() {
		//main application window
		window := new(app.Window)

		// Pass the window and the initialized game state to your UI's Run function
		if err := ui.RunMenu(window, currentGame); err != nil { // Pass currentGame here
			log.Fatal(err)
		}
		os.Exit(0) // Exit the program gracefully when the window closes
	}()

	//Gio application's main loop
	app.Main()
}
