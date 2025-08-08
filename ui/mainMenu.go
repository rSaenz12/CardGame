package ui

import (
	"CombinedCardgames/uiFunctions"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"log"

	"CombinedCardgames/baccaratGame"
	"CombinedCardgames/blackJackGame"
	"CombinedCardgames/goFishGame"
	"CombinedCardgames/signals"
)

func RunMainMenu(window *app.Window) error {

	var ops op.Ops

	//Declare clickable for each menu option
	var playBlackJackButton, playGoFishButton, playBaccaratButton, exitButton widget.Clickable
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
			log.Printf("Main Menu: Window destroyed: %v", e.Err)
			// Signal the main loop to exit the application
			select {
			case signals.ExitAppSignal <- true:
			default:
			}
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
			//calls baccarat ui and game
			if playBaccaratButton.Clicked(gtx) {
				baccaratGame.BaccaratMain(window)
			}

			// Exits program
			if exitButton.Clicked(gtx) {
				log.Println("Main Menu: Exit button clicked. Sending exit signal.")
				select {
				case signals.ExitAppSignal <- true: // Send the global exit signal
				default:
					log.Println("Main Menu: ExitAppSignal channel blocked.")
				}
				return nil
			}

			//Layout Stack of background image and menu buttons
			layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return imageWidget.Layout(gtx)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return uiFunctions.DrawFourButtons(gtx,
						&playBlackJackButton, //buton1
						&playGoFishButton,
						&playBaccaratButton, //button3
						&exitButton,         //button3
						"BlackJack",         //button1 text
						"GoFish",
						"Baccarat", //button2 text
						"Exit",     //button3 text
						retroRed,   //button1 color
						retroBlue,  //button2 color
						retroYellow,
						retroGreen, //button3 color
						white)      //text color
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
