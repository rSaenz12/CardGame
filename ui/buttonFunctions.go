package ui

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
)

// function that creates 2 buttons
func DrawTwoButtons(gtx layout.Context, button1Click, button2Click *widget.Clickable, button1Text, button2Text string) layout.Dimensions {
	targetY := int(float32(gtx.Constraints.Max.Y) * 0.9)
	//boxHeight := gtx.Constraints.Max.Y / 4

	offset := op.Offset(image.Pt(0, targetY)).Push(gtx.Ops)
	defer offset.Pop()

	th := material.NewTheme()

	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			hitButton := material.Button(th, button1Click, button1Text)
			hitButton.Background = green
			hitButton.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255}
			return hitButton.Layout(gtx)
		}),
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			standButton := material.Button(th, button2Click, button2Text)
			standButton.Background = red
			standButton.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255}
			return standButton.Layout(gtx)
		}),
	)
}

// function that creates 3 buttons
func DrawThreeButtons(gtx layout.Context, button1Click, button2Click, button3Click *widget.Clickable, button1Text, button2Text, button3Text string) layout.Dimensions {
	targetY := int(float32(gtx.Constraints.Max.Y) * 0.9)
	//boxHeight := gtx.Constraints.Max.Y / 4

	offset := op.Offset(image.Pt(0, targetY)).Push(gtx.Ops)
	defer offset.Pop()

	th := material.NewTheme()

	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			startButton := material.Button(th, button1Click, button1Text)
			startButton.Background = green
			startButton.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255}
			return startButton.Layout(gtx)
		}),
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			scoreButton := material.Button(th, button2Click, button2Text)
			scoreButton.Background = red
			scoreButton.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255}
			return scoreButton.Layout(gtx)
		}),
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			exitButton := material.Button(th, button3Click, button3Text)
			exitButton.Background = blue
			exitButton.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255}
			return exitButton.Layout(gtx)
		}),
	)
}
