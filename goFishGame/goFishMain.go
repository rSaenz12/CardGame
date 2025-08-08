package goFishGame

import (
	"fmt"
	"gioui.org/app"

	"CombinedCardgames/goFishGame/Ui" //import UI for go fish

	"os"
)

func GoFishMain(window *app.Window) {

	for {
		choice, err := Ui.RunMenu(window)

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
