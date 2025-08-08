package ui

import (
	"CombinedCardgames/blackJackGame/game"
	ui2 "CombinedCardgames/uiFunctions"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

func RunMenu(window *app.Window, gameInstance *game.Game) (string, error) {

	var ops op.Ops

	//Declare clickable for each menu option
	var startGameButton, scoreButton, mainMenuButton, exitButton widget.Clickable

	background := ui2.LoadImage(backgroundPath)
	bgOp := paint.NewImageOp(background)

	// Declare imageWidget once outside the loop for correct usage
	var imageWidget widget.Image

	imageWidget = widget.Image{
		Src: bgOp,
		Fit: widget.Cover,
	}

	for {
		switch e := window.Event().(type) {

		case app.DestroyEvent:
			return "Error", e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Clear ops for the new frame
			ops.Reset()

			// Calls first game Ui and game input
			if startGameButton.Clicked(gtx) {
				gameInstance.PlayGame("1")
				choice, err := RunPhaseOne(window, gameInstance)
				if err != nil {
					return "Error", err
				}

				return choice, nil
			}
			//Calls Score UI and Game input
			if scoreButton.Clicked(gtx) {
				gameInstance.PlayGame("2")
				choice, err := RunScoreUi(window, gameInstance)
				if err != nil {
					return "Error", err
				}
				return choice, nil
			}
			//returns to Main Menu
			if mainMenuButton.Clicked(gtx) {
				return "mainMenu", nil
			}

			// Exits program
			if exitButton.Clicked(gtx) {
				return "exit", nil
			}

			//Layout Stack of background image and menu buttons
			layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return imageWidget.Layout(gtx)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return ui2.DrawFourButtons(gtx,
						//button variables
						&startGameButton, &scoreButton, &mainMenuButton, &exitButton,
						//button text
						"Start Game", "Score", "Main Menu", "Exit",
						//button colors
						green, red, blue, yellow,
						//text color
						white)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
