//Game internal menu UI

package ui

import (
	ui2 "CombinedCardgames/uiFunctions"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"image/color"
)
//This fuction presents the menu for the game and returns a string message based on button clicks. 
func RunBaccaratMenu(window *app.Window) string {
	var ops op.Ops

	var startGameButton, mainMenuButton, exitButton widget.Clickable

	background := ui2.LoadImage(backgroundPath)
	bgOp := paint.NewImageOp(background)

	var imageWidget = widget.Image{Src: bgOp, Fit: widget.Cover}

	for {
		e := window.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return "exit"
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if startGameButton.Clicked(gtx) {
				return "start"
			}
			if mainMenuButton.Clicked(gtx) {
				return "main_menu"
			}
			if exitButton.Clicked(gtx) {
				return "exit"
			}

			layout.Stack{}.Layout(gtx,
				layout.Expanded(imageWidget.Layout),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return ui2.DrawThreeButtons(gtx,
						&startGameButton, &mainMenuButton, &exitButton,
						"Start Game", "Main Menu", "Exit",
						color.NRGBA{G: 150, A: 255}, color.NRGBA{B: 150, A: 255}, color.NRGBA{R: 150, A: 255},
						color.NRGBA{R: 255, G: 255, B: 255, A: 255})
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
