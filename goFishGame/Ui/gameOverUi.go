package Ui

import (
	ui2 "CombinedCardgames/uiFunctions"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"os"
)

func runEndScreen(window *app.Window, finalMessage string) error {
	var ops op.Ops

	background := ui2.LoadImage(backgroundImagePath)
	bgOp := paint.NewImageOp(background)

	//Declaring clickables for UI
	var menuButton, exit widget.Clickable

	th := material.NewTheme()

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

			textPosX := int(float32(gtx.Constraints.Max.X) / 3)
			textPosY := int(float32(gtx.Constraints.Max.Y) * .4)

			//Calls Menu UI, no need for a game input
			if menuButton.Clicked(gtx) {
				//RunMainMenu(window)
			}
			//Exits the program completely
			if exit.Clicked(gtx) {
				os.Exit(0)
				return nil
			}

			//Layout Stack of background image, cards, and buttons
			layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return imageWidget.Layout(gtx)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:  unit.Dp(textPosY + 25), // Adjust these values to move the score
						Left: unit.Dp(textPosX),      // from the top-left corner
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

						h2 := material.H2(th, finalMessage)
						h2.Color = white
						return h2.Layout(gtx)
					})
				}), layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return ui2.DrawTwoButtons(gtx,
						//button variables
						&menuButton, &exit,
						//button texts
						"Exit to Menu", "Exit to Desktop",
						//button colors
						materialRed, materialBlue,
						//text color
						black)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
