package ui

import (
	"blackJack/game"
	"fmt"
	"gioui.org/app"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"image/color"
)

var (
	red   = color.NRGBA{200, 50, 50, 255}
	green = color.NRGBA{0, 150, 0, 255}
	blue  = color.NRGBA{0, 70, 200, 255} // Deeper blue
)

func RunMenu(window *app.Window, gameInstance *game.Game) error {

	var ops op.Ops

	//Declare clickable for each menu option
	var startGameButton, scoreButton, exitButton widget.Clickable

	background := LoadImage("blackJackTable.png")
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
				fmt.Println("Start Game clicked")

				gameInstance.PlayGame("1")

				err := RunPhaseOne(window, gameInstance)
				if err != nil {
					return err
				}
			}
			//Calls Score UI and Game input
			if scoreButton.Clicked(gtx) {
				fmt.Println("Score Clicked")

				gameInstance.PlayGame("2")

				err := RunScorUi(window, gameInstance)
				if err != nil {
					return err
				}

			}
			// Exits program
			if exitButton.Clicked(gtx) {
				fmt.Println("Exit Clicked!")

				gameInstance.PlayGame("3")

			}

			//Layout Stack of background image and menu buttons
			layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return imageWidget.Layout(gtx)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return DrawThreeButtons(gtx, &startGameButton, &scoreButton, &exitButton, "Start Game", "Score", "Exit")
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
