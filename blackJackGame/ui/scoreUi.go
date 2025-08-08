package ui

import (
	"CombinedCardgames/blackJackGame/game"
	ui2 "CombinedCardgames/uiFunctions"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"image/color"
	"strconv"

	"gioui.org/widget"
	"gioui.org/widget/material"
)

func RunScorUi(window *app.Window, gameInstance *game.Game) error {

	var ops op.Ops

	//Declaring clickables for UI
	var menuButton, exitButton widget.Clickable

	// Load the background image once outside the loop for efficiency.
	background := ui2.LoadImage(backgroundPath)
	bgOp := paint.NewImageOp(background)

	// Declare imageWidget once outside the loop for correct usage
	var imageWidget widget.Image

	imageWidget = widget.Image{
		Src: bgOp,
		Fit: widget.Cover, // layout.Cover is assumed to be resolved now
	}
	//creating theme to be used
	th := material.NewTheme()

	for {
		switch e := window.Event().(type) {

		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			//Variables used for placement of elements, declared within loop due to gtx
			boxWidth := int(float32(gtx.Constraints.Max.X) / 2)
			boxHeight := int(float32(gtx.Constraints.Max.Y) / 3)

			textPosX := int(float32(gtx.Constraints.Max.X)*.45) - (boxWidth / 9)
			textPosY := int(float32(gtx.Constraints.Max.Y) * .1)

			// Clear ops for the new frame
			ops.Reset()

			//Sends user to menuUi.go
			//No game input needed
			if menuButton.Clicked(gtx) {
				err := RunMenu(window, gameInstance)
				if err != nil {
					return err
				}

			}
			//uses game function to exit program
			if exitButton.Clicked(gtx) {
				gameInstance.PlayGame("3")
			}
			//Layout Stack of background image,Box,Text, and buttons
			layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return imageWidget.Layout(gtx)
				}),
				//creates the box that the text is displayed over
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					offsetX := (gtx.Constraints.Max.X / 2) - (boxWidth / 2)
					offsetY := (gtx.Constraints.Max.Y / 4) - (boxHeight / 2)

					defer op.Offset(image.Point{
						X: offsetX,
						Y: offsetY,
					}).Push(gtx.Ops).Pop()

					// White background box
					rect := clip.Rect{
						Max: image.Pt(boxWidth, boxHeight),
					}.Push(gtx.Ops)
					paint.Fill(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}) // white
					rect.Pop()
					return layout.Dimensions{Size: image.Pt(boxWidth, boxHeight)}
				}),
				//"Score"
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					// --- Option 1: Position with Inset (recommended for simple offsets) ---
					return layout.Inset{
						Top:  unit.Dp(textPosY), // Adjust these values to move the score
						Left: unit.Dp(textPosX), // from the top-left corner
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.H2(th, "Score").Layout(gtx)
					})

				}),
				//Wins,Losses, points
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:  unit.Dp(int(float32(textPosX) * .4)),
						Left: unit.Dp(textPosX),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:      layout.Vertical,
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Body1(th, "Wins: "+strconv.Itoa(gameInstance.Wins)).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Body1(th, "Loses: "+strconv.Itoa(gameInstance.Losses)).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Body1(th, "Points: "+strconv.Itoa(gameInstance.UserPoints)).Layout(gtx)
							}),
						)
					})
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return ui2.DrawTwoButtons(gtx,
						//button variables
						&menuButton, &exitButton,
						//button text
						"Menu", "Exit",
						//button colors
						green, red,
						//text color
						white)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
