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

func RunMenu(window *app.Window, gameInstance *game.Game) error {

	var ops op.Ops

	//Declare clickable for each menu option
	var startGameButton, scoreButton, exitButton widget.Clickable

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
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Clear ops for the new frame
			ops.Reset()

			// Calls first game Ui and game input
			if startGameButton.Clicked(gtx) {
				gameInstance.PlayGame("1")
				err := RunPhaseOne(window, gameInstance)
				if err != nil {
					return err
				}
			}
			//Calls Score UI and Game input
			if scoreButton.Clicked(gtx) {
				gameInstance.PlayGame("2")
				err := RunScorUi(window, gameInstance)
				if err != nil {
					return err
				}
			}
			// Exits program
			if exitButton.Clicked(gtx) {
				gameInstance.PlayGame("3")
			}

			//Layout Stack of background image and menu buttons
			layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return imageWidget.Layout(gtx)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return ui2.DrawThreeButtons(gtx,
						//button variables
						&startGameButton, &scoreButton, &exitButton,
						//button text
						"Start Game", "Score", "Exit",
						//button colors
						green, red, blue,
						//text color
						white)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
