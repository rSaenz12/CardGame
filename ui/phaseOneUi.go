package ui

import (
	"blackJack/game"
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/widget"
	"log"
	"os"

	"image"
	"image/png"

	"gioui.org/op"
	"gioui.org/op/paint"
)

func RunPhaseOne(window *app.Window, gameInstance *game.Game) error {

	var ops op.Ops

	//Declaring clickables for UI
	var hitButton, standButton, doubleDownButton widget.Clickable

	//Hands for this UI image
	var userCards []image.Image // Slice to hold the card images
	var dealerCards []image.Image

	//when true (upon win or loss), next screen is the winLossScreenUi.go
	switchToWinOrLossScreen := false

	//scale for objects,
	scale := float32(0.25) // 25% size

	// Load the background image once outside the loop for efficiency.
	background := LoadImage("blackJackTable.png")
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
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Clear ops for the new frame
			ops.Reset()

			//Update card images for display if the game state changes
			userCards = make([]image.Image, len(gameInstance.UserHand))
			userCards = GetCardImage(gameInstance.UserHand)
			dealerCards = make([]image.Image, len(gameInstance.DealerHand))
			dealerCards = GetCardImage(gameInstance.DealerHand)

			//dealer does not reveal their second card until its dealers turn
			if !gameInstance.CheckRevealDealer() {
				if len(dealerCards) != 0 {
					dealerCards = dealerCards[:1]
					dealerCards = append(dealerCards, LoadImage("deckImages/1B.png"))
				}
			}

			// Hit calls next UI and game input
			if hitButton.Clicked(gtx) {
				fmt.Println("Hit")

				gameInstance.PhaseOne("1")
				err := RunHitUserUi(window, gameInstance)
				if err != nil {
					return err
				}

			}

			// Stand calls next game input
			//***NOTE*** There is no UI being ran from here, by game logic the game should end-
			//causing the winLossScreenUi.go to be called
			if standButton.Clicked(gtx) {
				fmt.Println("Stand")
				gameInstance.PhaseOne("2")
			}

			// Doubles down with game input
			//***NOTE*** There is no UI being ran from here, by game logic the game should end-
			//causing the winLossScreenUi.go to be called
			if doubleDownButton.Clicked(gtx) {
				fmt.Println("Double Down")
				gameInstance.PhaseOne("3")
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
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							// vertical offset for the cards to be in correct places
							offset := op.Offset(image.Pt(0, yOffset)).Push(gtx.Ops)
							defer offset.Pop()

							// Loop through cards and draw each with spacing
							cardWidth := 25
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
					return DrawThreeButtons(gtx, &hitButton, &standButton, &doubleDownButton, "Hit", "Stand", "Double Down")
				}),
			)
			if gameInstance.CheckGameEnded() {
				switchToWinOrLossScreen = true
			}

			e.Frame(gtx.Ops)

			if hitButton.Clicked(gtx) || standButton.Clicked(gtx) || doubleDownButton.Clicked(gtx) {
				//force repaint of ui
				window.Invalidate()
			}
			if switchToWinOrLossScreen {
				fmt.Println("SOMEONE WON")
				err := RunWinOrLossScreen(window, gameInstance)
				if err != nil {
					return err
				}
				return err
			}
		}
	}
}

// LoadImage loads the images using the path
func LoadImage(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("ERROR: Failed to open image %s: %v\n", path, err)
	}
	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("ERROR: Failed to decode image %s: %v\n", path, err)
	}
	// Check if closing the file results in an error
	if err := f.Close(); err != nil {
		log.Fatalf("ERROR: Failed to close image file %s: %v\n", path, err)
	}
	return img
}

// GetCardImage grabs the path,loads images, adds to an array of images
func GetCardImage(currentHand []string) []image.Image {
	var images []image.Image

	//loops through the current hand calling each card, adding them as slices of images to the array
	for _, card := range currentHand {
		path := "deckImages/" + card + ".png"

		f, err := os.Open(path)
		if err != nil {
			log.Fatalf("ERROR: Failed to open image %s: %v\n", path, err)
		}

		img, err := png.Decode(f)
		if err != nil {
			log.Fatalf("ERROR: Failed to decode image %s: %v\n", path, err)
		}

		// Check if closing the file results in an error
		if err := f.Close(); err != nil {
			log.Fatalf("ERROR: Failed to close image file %s: %v\n", path, err)
		}

		images = append(images, img)
	}
	return images
}
