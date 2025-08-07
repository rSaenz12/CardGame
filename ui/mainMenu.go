package ui

import (
	"CombinedCardgames/uiFunctions"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"os"

	"CombinedCardgames/blackJackGame"
	"CombinedCardgames/goFishGame"
)

func RunMainMenu(window *app.Window) error {

	var ops op.Ops

	//Declare clickable for each menu option
	var playBlackJackButton, playGoFishButton, exitButton widget.Clickable

	background := uiFunctions.LoadImage("arcadeBackground.png")
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
			if playBlackJackButton.Clicked(gtx) {
				blackJackGame.BlackJackMain(window)
			}
			//Calls Score UI and Game input
			if playGoFishButton.Clicked(gtx) {
				goFishGame.GoFishMain(window)
			}
			// Exits program
			if exitButton.Clicked(gtx) {
				os.Exit(0)
			}

			//Layout Stack of background image and menu buttons
			layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return imageWidget.Layout(gtx)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return uiFunctions.DrawThreeButtons(gtx,
						&playBlackJackButton, //buton1
						&playGoFishButton,    //button2
						&exitButton,          //button3
						"BlackJack",          //button1 text
						"GoFish",             //button2 text
						"Exit",               //button3 text
						retroRed,             //button1 color
						retroBlue,            //button2 color
						retroYellow,          //button3 color
						white)                //text color
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
