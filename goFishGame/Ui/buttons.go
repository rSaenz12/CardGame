package Ui

/*
func DrawDynamicButtons(gtx layout.Context, theme *material.Theme, clicks []*widget.Clickable, texts []string) layout.Dimensions {

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

	// 1. Create a slice to hold all the FlexChild elements
	children := make([]layout.FlexChild, len(clicks))

	// 2. Populate the children slice within the loop
	for i := range clicks {
		buttonClick := clicks[i]
		buttonText := texts[i]
		buttonColor := materialRed
		switch i % 2 {
		case 0:
			buttonColor = materialBlue
		case 1:
			buttonColor = materialRed
		}

		// Each element in 'children' is a layout.Flexed or layout.Rigid
		// which internally holds the function that draws the button.
		children[i] = layout.Flexed(flexRatio, func(gtx layout.Context) layout.Dimensions {
			button := material.Button(theme, buttonClick, buttonText)
			button.Background = buttonColor
			button.Color = black
			return button.Layout(gtx)
		})
	}

	// 3. Pass the 'children' slice (using the variadic ... operator) directly to Layout
	return layout.Flex{
		Axis: layout.Horizontal,
		// Spacing: layout.SpaceAround, // Optional
	}.Layout(gtx, children...) // <-- Corrected: Pass the slice of FlexChild
}

// DrawTwoButtons function that creates 2 buttons
func DrawTwoButtons(gtx layout.Context, button1Click, button2Click *widget.Clickable, button1Text, button2Text string, buttonColor1, buttonColor2 color.NRGBA) layout.Dimensions {
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
			hitButton.Background = buttonColor1
			hitButton.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255}
			return hitButton.Layout(gtx)
		}),
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			standButton := material.Button(th, button2Click, button2Text)
			standButton.Background = buttonColor2
			standButton.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255}
			return standButton.Layout(gtx)
		}),
	)
}
*/
