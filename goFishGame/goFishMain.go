package goFishGame

import (
	"gioui.org/app"

	"CombinedCardgames/goFishGame/Ui"            //import UI for go fish
	"CombinedCardgames/goFishGame/goFishBackEnd" //import backend for GoFish

	"log"
	"os"
)

func GoFishMain(window *app.Window) {
	// Initializing game logic first
	gameInstance, err := goFishBackEnd.NewGame()
	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	// Starting the Gio GUI in a separate goroutine
	go func() {
		//main application window
		//window := new(app.Window)

		// Pass the window and the initialized game state to your UI's RunGoFishUi function
		if err := Ui.RunGoFishUi(window, gameInstance); err != nil { // Pass currentGame here
			log.Fatal(err)
		}
		//exit
		os.Exit(0)
	}()

	//Gio application main loop
	app.Main()

}
