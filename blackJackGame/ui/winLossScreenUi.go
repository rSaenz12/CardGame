package ui

import (
	"CombinedCardgames/blackJackGame/game"
	ui2 "CombinedCardgames/uiFunctions"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"image"

	"gioui.org/widget"
	"gioui.org/widget/material"
)

func RunWinOrLossScreen(window *app.Window, gameInstance *game.Game) error {

	var ops op.Ops

	//default message for the screen
	var finalMessage = "You lose!"

	//Declaring clickables for UI
	var menuButton, playAgainButton widget.Clickable

	//Hands for this UI image
	var userCards []image.Image // Slice to hold the card images
	var userHand []string
	var dealerCards []image.Image
	var dealerHand []string
	//creating array of strings for loading users card image2
	for _, card := range gameInstance.UserHand {
		cardString := card.Rank + card.Suit
		userHand = append(userHand, cardString)
	}
	//creating array of strings for loading back of card images
	for _, card := range gameInstance.DealerHand {
		cardString := card.Rank + card.Suit
		dealerHand = append(dealerHand, cardString)
	}

	//scale for objects, blown up to 50% for looks
	scale := float32(0.5)

	// Load the background image once outside the loop for efficiency.
	background := ui2.LoadImage(backgroundPath)
	bgOp := paint.NewImageOp(background)

	// Declare imageWidget once outside the loop for correct usage
	var imageWidget widget.Image

	imageWidget = widget.Image{
		Src: bgOp,
		Fit: widget.Cover,
	}

	//creating theme to be used
	th := material.NewTheme()

	//Checks for conditions to change finalMessage
	if gameInstance.CheckBlackJack() {
		finalMessage = "Black Jack! 3:2 Pay Out!"
	} else if gameInstance.CheckUserWin() {
		finalMessage = "You Win!"
	} else if gameInstance.CheckTieGame() {
		finalMessage = "You Tie!"
	}

	for {
		switch e := window.Event().(type) {

		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Clear ops for the new frame
			ops.Reset()

			//Update card images for display if the game state changes
			userCards = make([]image.Image, len(gameInstance.UserHand))
			userCards = ui2.GetCardImage(userHand)
			dealerCards = make([]image.Image, len(gameInstance.DealerHand))
			dealerCards = ui2.GetCardImage(dealerHand)

			//Calls the phaseOne UI, running PlayGame input of 1 to return to correct state
			if playAgainButton.Clicked(gtx) {
				gameInstance.PlayGame("1")
				err := RunPhaseOne(window, gameInstance)
				if err != nil {
					return err
				}
			}
			//Calls Menu UI, no need for a game input
			if menuButton.Clicked(gtx) {
				err := RunMenu(window, gameInstance)
				if err != nil {
					return err
				}
			}
			//Layout Stack of background image, cards, and buttons
			layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return imageWidget.Layout(gtx)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					yOffset := gtx.Constraints.Max.Y / 10

					return layout.Flex{
						Axis:      layout.Vertical,
						Alignment: layout.Middle,
						Spacing:   layout.SpaceStart,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.H2(th, finalMessage).Layout(gtx)
						}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:      layout.Horizontal,
								Spacing:   layout.SpaceStart,
								Alignment: layout.Middle,
							}.Layout(gtx,
								// Layout layer responsible for rendering the user's cards.
								layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions { // <-- Corrected
									// apply a vertical offset for the row of userCards
									offset := op.Offset(image.Pt(0, yOffset)).Push(gtx.Ops)
									defer offset.Pop()

									// Loop through userCards and draw each with spacing
									cardWidth := 25 // approx width after scaling
									spacing := 0
									cardHeight := int(float32(gtx.Constraints.Max.Y) * 0.6)
									totalWidth := len(userCards)*(cardWidth+spacing) - spacing
									startingOffset := (gtx.Constraints.Max.X - totalWidth) / 2

									//loop to "print" (paint) the cards
									for i, card := range userCards {
										x := startingOffset + i*(cardWidth+spacing)
										cardOffset := op.Offset(image.Pt(x, cardHeight)).Push(gtx.Ops)
										scaleOp := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Push(gtx.Ops)

										//prints card
										paint.NewImageOp(card).Add(gtx.Ops)
										paint.PaintOp{}.Add(gtx.Ops)

										scaleOp.Pop()
										cardOffset.Pop()
									}
									return layout.Dimensions{Size: image.Pt(len(userCards)*(cardWidth+spacing), cardWidth)}
								}), // End of layout.Rigid
								// Layout layer responsible for rendering the dealer's cards.
								layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions { // <-- Corrected
									// vertical offset for the cards to be in correct places
									offset := op.Offset(image.Pt(0, yOffset)).Push(gtx.Ops)
									defer offset.Pop()

									// Loop through userCards and draw each with spacing
									cardWidth := 25 // approx width after scaling
									spacing := 0
									cardHeight := int(float32(gtx.Constraints.Max.Y) * .001)
									totalWidth := len(userCards)*(cardWidth+spacing) - spacing
									startingOffset := (gtx.Constraints.Max.X - totalWidth) / 2

									//loop to "print" (paint) the cards
									for i, card := range dealerCards {
										x := startingOffset + i*(cardWidth+spacing)
										cardOffset := op.Offset(image.Pt(x, cardHeight)).Push(gtx.Ops)
										scaleOp := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Push(gtx.Ops)

										//prints card
										paint.NewImageOp(card).Add(gtx.Ops)
										paint.PaintOp{}.Add(gtx.Ops) // <-- This actually draws the image

										scaleOp.Pop()
										cardOffset.Pop()
									}
									return layout.Dimensions{Size: image.Pt(len(dealerCards)*(cardWidth+spacing), cardWidth)}
								}), //End of layer.Rigid
							)
						}),
					)
				}), layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return ui2.DrawTwoButtons(gtx,
						//button variables
						&playAgainButton, &menuButton,
						//button text
						"Play Again", "Menu",
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
