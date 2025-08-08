package uiFunctions

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"log"
)

// DrawTwoButtons function that creates 2 buttons
func DrawTwoButtons(gtx layout.Context, button1Click, button2Click *widget.Clickable, button1Text, button2Text string, buttonColor1, buttonColor2, textColor color.NRGBA) layout.Dimensions {
	targetY := int(float32(gtx.Constraints.Max.Y) * 0.9)

	offset := op.Offset(image.Pt(0, targetY)).Push(gtx.Ops)
	defer offset.Pop()

	th := material.NewTheme()

	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			hitButton := material.Button(th, button1Click, button1Text)
			hitButton.Background = buttonColor1
			hitButton.Color = textColor
			return hitButton.Layout(gtx)
		}),
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			standButton := material.Button(th, button2Click, button2Text)
			standButton.Background = buttonColor2
			standButton.Color = textColor
			return standButton.Layout(gtx)
		}),
	)
}

// DrawThreeButtons function that creates 3 buttons
func DrawThreeButtons(gtx layout.Context, button1Click, button2Click, button3Click *widget.Clickable, button1Text, button2Text, button3Text string, buttonColor1, buttonColor2, buttonColor3, textColor color.NRGBA) layout.Dimensions {
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
			startButton.Background = buttonColor1
			startButton.Color = textColor
			return startButton.Layout(gtx)
		}),
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			scoreButton := material.Button(th, button2Click, button2Text)
			scoreButton.Background = buttonColor2
			scoreButton.Color = textColor
			return scoreButton.Layout(gtx)
		}),
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			exitButton := material.Button(th, button3Click, button3Text)
			exitButton.Background = buttonColor3
			exitButton.Color = textColor
			return exitButton.Layout(gtx)
		}),
	)
}

// DrawFourButtons function that creates 4 buttons
func DrawFourButtons(gtx layout.Context, button1Click, button2Click, button3Click, button4Click *widget.Clickable, button1Text, button2Text, button3Text, button4Text string, buttonColor1, buttonColor2, buttonColor3, buttonColor4, textColor color.NRGBA) layout.Dimensions {
	targetY := int(float32(gtx.Constraints.Max.Y) * 0.9)

	offset := op.Offset(image.Pt(0, targetY)).Push(gtx.Ops)
	defer offset.Pop()

	th := material.NewTheme()

	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
			button1 := material.Button(th, button1Click, button1Text)
			button1.Background = buttonColor1
			button1.Color = textColor
			return button1.Layout(gtx)
		}),
		layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
			button2 := material.Button(th, button2Click, button2Text)
			button2.Background = buttonColor2
			button2.Color = textColor
			return button2.Layout(gtx)
		}),
		layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
			button3 := material.Button(th, button3Click, button3Text)
			button3.Background = buttonColor3
			button3.Color = textColor
			return button3.Layout(gtx)
		}),
		layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
			button4 := material.Button(th, button4Click, button4Text)
			button4.Background = buttonColor4
			button4.Color = textColor
			return button4.Layout(gtx)
		}),
	)
}

func DrawDynamicButtons(gtx layout.Context, theme *material.Theme, clicks []*widget.Clickable, texts []string, buttonColor1, buttonColor2, textColor color.NRGBA) layout.Dimensions {

	if len(clicks) != len(texts) {
		log.Printf("Warning: Mismatch in number of clickables (%d) and texts (%d)", len(clicks), len(texts))
		return layout.Dimensions{}
	}
	if len(clicks) == 0 {
		return layout.Dimensions{}
	}

	targetY := int(float32(gtx.Constraints.Max.Y) * 0.9)

	offset := op.Offset(image.Pt(0, targetY)).Push(gtx.Ops)
	defer offset.Pop()

	flexRatio := 1.0 / float32(len(clicks))

	//create element
	children := make([]layout.FlexChild, len(clicks))

	//using loop to generate buttons
	for i := range clicks {
		buttonClick := clicks[i]
		buttonText := texts[i]
		buttonColor := buttonColor1
		switch i % 2 {
		case 0:
			buttonColor = buttonColor2
		case 1:
			buttonColor = buttonColor1
		}

		// Each element in 'children' is a layout.Flexed or layout.Rigid
		// which internally holds the function that draws the button.
		children[i] = layout.Flexed(flexRatio, func(gtx layout.Context) layout.Dimensions {
			button := material.Button(theme, buttonClick, buttonText)
			button.Background = buttonColor
			button.Color = textColor
			return button.Layout(gtx)
		})
	}

	// returning childrens slice
	return layout.Flex{
		Axis: layout.Horizontal,
		// Spacing: layout.SpaceAround, // Optional
	}.Layout(gtx, children...) // <-- Corrected: Pass the slice of FlexChild
}
