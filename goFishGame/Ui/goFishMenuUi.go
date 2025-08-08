package Ui

import (
	"CombinedCardgames/goFishGame/goFishBackEnd"
	ui2 "CombinedCardgames/uiFunctions"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

func RunMenu(window *app.Window) (string, error) {
	var ops op.Ops

	//Declare clickable for each menu option
	var startGameButton, mainMenuButton, exitButton widget.Clickable

	background := ui2.LoadImage(backgroundImagePath)
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
				gameInstance, err := goFishBackEnd.NewGame()
				choice, err := RunGoFishUi(window, gameInstance)
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
					return ui2.DrawThreeButtons(gtx,
						//button variables
						&startGameButton, &mainMenuButton, &exitButton,
						//button text
						"Start Game", "Main Menu", "Exit",
						//button colors
						materialRed, materialBlue, materialYellow,
						//text color
						black)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
