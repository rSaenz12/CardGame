package main

import (
	"CombinedCardgames/signals"
	"CombinedCardgames/ui"
	"gioui.org/app"

	"log"
	"os"
)

func main() {
	// Starting the Gio GUI in a separate goroutine
	go func() {
		//main application window
		window := new(app.Window)

		// Pass the window and the initialized game state to your UI's Run function
		if err := ui.RunMainMenu(window); err != nil { // Pass currentGame here
			log.Fatal(err)
		}
		//exit
		sig := <-signals.MenuSignal
		if sig {
			if err := ui.RunMainMenu(window); err != nil { // Pass currentGame here
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}()

	//Gio application's main loop
	app.Main()
}
