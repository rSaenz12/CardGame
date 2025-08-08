package ui

import (
	"CombinedCardgames/blackJackGame/game"
	ui2 "CombinedCardgames/uiFunctions"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"image"
)

func RunHitUserUi(window *app.Window, gameInstance *game.Game) (string, error) {

	var ops op.Ops

	//Declaring clickables for UI
	var hitButton, standButton widget.Clickable

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

	//when true (upon win or loss), next screen is the winLossScreenUi.go
	switchToWinOrLossScreen := false

	scale := float32(0.25) // 25% size

	// Load the background image once outside the loop for efficiency.
	background := ui2.LoadImage(backgroundPath)
	bgOp := paint.NewImageOp(background)

	// Declare imageWidget once outside the loop for correct usage
	var imageWidget widget.Image

	// Assign to the already declared imageWidget
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

			//Update card images for display if the game state changes
			userCards = make([]image.Image, len(gameInstance.UserHand))
			userCards = ui2.GetCardImage(userHand)
			dealerCards = make([]image.Image, len(gameInstance.DealerHand))
			dealerCards = ui2.GetCardImage(dealerHand)

			//dealer does not reveal their second card until its dealers turn
			if !gameInstance.CheckRevealDealer() {
				if len(dealerCards) != 0 {
					dealerCards = dealerCards[:1]
					dealerCards = append(dealerCards, ui2.LoadImage("deckImages/1B.png"))
				}
			}

			// Hit calls next UI and game input
			if hitButton.Clicked(gtx) {
				gameInstance.HitUser("1")
			}

			// Stand calls next game input
			//***NOTE*** There is no UI being ran from here, by game logic the game should end-
			//causing the winLossScreenUi.go to be called
			if standButton.Clicked(gtx) {
				gameInstance.HitUser("2")
			}

			//Layout Stack of background image, cards, and buttons
			layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return imageWidget.Layout(gtx)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					yOffset := gtx.Constraints.Max.Y / 10

					return layout.Flex{
						Axis:      layout.Horizontal,
						Spacing:   layout.SpaceStart,
						Alignment: layout.Middle,
					}.Layout(gtx,
						// Layout layer responsible for rendering the user's cards.
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							// apply a vertical offset for the row of cards
							offset := op.Offset(image.Pt(0, yOffset)).Push(gtx.Ops)
							defer offset.Pop()

							// Loop through cards and draw each with spacing
							cardWidth := 25 // approx width after scaling
							spacing := 0
							cardHeight := int(float32(gtx.Constraints.Max.Y) * 0.6)
							startingOffset := int(float32(gtx.Constraints.Max.X) / 2)

							for i, card := range userCards {
								cardOffset := op.Offset(image.Pt(startingOffset+(i*(cardWidth+spacing)), cardHeight)).Push(gtx.Ops)
								scaleOp := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Push(gtx.Ops)

								//prints card
								paint.NewImageOp(card).Add(gtx.Ops)
								paint.PaintOp{}.Add(gtx.Ops)

								scaleOp.Pop()
								cardOffset.Pop()
							}
							return layout.Dimensions{Size: image.Pt(len(userCards)*(cardWidth+spacing), cardWidth)}
						}), //End of layout.Rigid
						// Layout layer responsible for rendering the dealer's cards.
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							// vertical offset for the cards to be in correct places
							offset := op.Offset(image.Pt(0, yOffset)).Push(gtx.Ops)
							defer offset.Pop()

							// Loop through cards and draw each with spacing
							cardWidth := 25 // approx width after scaling
							spacing := 0
							cardHeight := int(float32(gtx.Constraints.Max.Y) * .001)
							startingOffset := int(float32(gtx.Constraints.Max.X) / 2)

							//loop to "print" (paint) the cards
							for i, card := range dealerCards {
								cardOffset := op.Offset(image.Pt(startingOffset+(i*(cardWidth+spacing)), cardHeight)).Push(gtx.Ops)
								scaleOp := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Push(gtx.Ops)

								//prints card
								paint.NewImageOp(card).Add(gtx.Ops)
								paint.PaintOp{}.Add(gtx.Ops)

								scaleOp.Pop()
								cardOffset.Pop()
							}
							return layout.Dimensions{Size: image.Pt(len(dealerCards)*(cardWidth+spacing), cardWidth)}
						}), //End of layout.Rigid
					)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return ui2.DrawTwoButtons(gtx,
						//button variables
						&hitButton, &standButton,
						//button text
						"Hit", "Stand",
						//button colors
						green, red,
						// text color
						white)
				}),
			)
			if gameInstance.CheckGameEnded() {
				switchToWinOrLossScreen = true
			}

			e.Frame(gtx.Ops)

			if hitButton.Clicked(gtx) || standButton.Clicked(gtx) {
				window.Invalidate() //Forces repaint on next loop
			}
			if switchToWinOrLossScreen {
				choice, err := RunWinOrLossScreen(window, gameInstance)
				if err != nil {
					return "Error", err
				}
				return choice, nil
			}
		}
	}
}
