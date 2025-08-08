package Ui

import (
	"CombinedCardgames/goFishGame/goFishBackEnd"
	"CombinedCardgames/goFishGame/logHandling"

	ui2 "CombinedCardgames/uiFunctions"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"image"
	"image/color"
	"strconv"
	"strings"
	"time"
)

func RunGoFishUi(window *app.Window, game *goFishBackEnd.GoFish) (string, error) {
	var ops op.Ops

	menuButton := new(widget.Clickable)

	// Load the background image once outside the loop for efficiency.
	background := ui2.LoadImage(backgroundImagePath)
	bgOp := paint.NewImageOp(background)

	// Declare imageWidget once outside the loop for correct usage
	var imageWidget widget.Image
	// Assign to the already declared imageWidget
	imageWidget = widget.Image{
		Src: bgOp,
		Fit: widget.Cover,
	}

	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	numButtons := len(ranks)
	clicks := make([]*widget.Clickable, numButtons)

	th := material.NewTheme()
	for i := range clicks {
		clicks[i] = new(widget.Clickable)
	}

	for {
		switch e := window.Event().(type) {

		case app.DestroyEvent:
			return "Error", e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			//Variables used for placement of elements, declared within loop due to gtx
			boxWidth := int(float32(gtx.Constraints.Max.X) / 2)
			boxHeight := int(float32(gtx.Constraints.Max.Y) / 3)
			textPosX := int(float32(gtx.Constraints.Max.X)*.35) - (boxWidth / 8)
			textPosY := int(float32(gtx.Constraints.Max.Y)*.45) - (boxHeight / 9)

			//Update card images for display if the game state changes
			var userCards []image.Image
			var userHand []string
			var computerCards []image.Image
			var computerHand []string

			//creating array of strings for loading users card image2
			for _, card := range game.UserPlayer.Hand {
				cardString := card.Rank + card.Suit
				userHand = append(userHand, cardString)
			}
			//creating array of strings for loading back of card images
			for range game.ComputerPlayer.Hand {
				computerHand = append(computerHand, faceDownCardPath)
			}

			// Clear ops for the new frame
			ops.Reset()

			//Update card images for display if the game state changes
			userCards = make([]image.Image, len(game.UserPlayer.Hand))
			userCards = ui2.GetCardImage(userHand)
			userCardHeight := int(float32(gtx.Constraints.Max.Y) * .6)

			computerCards = make([]image.Image, len(game.ComputerPlayer.Hand))
			computerCards = ui2.GetCardImage(computerHand)
			computerCardHeight := int(float32(gtx.Constraints.Max.Y) * .001)

			//scale for objects,
			scale := float32(0.25) // 25% size

			if menuButton.Clicked(gtx) {
				return "", nil
			}

			// handling button clicks and game turns
			if game.CurrentTurn == "user" {
				for i, click := range clicks {
					if click.Clicked(gtx) {
						clickedRank := ranks[i]
						logHandling.AppendLog("User clicked to ask for rank: " + clickedRank)

						// Execute the player's turn in the backend game logic
						game.PlayerTurn(clickedRank)

						// After the player's action, switch turn to computer
						game.CurrentTurn = "computer"
						logHandling.AppendLog("Turn switched to:" + game.CurrentTurn)

						// redraw frame, showing the updated game state
						window.Invalidate()

						break // Only process one button click per frame

					}
				}
			}

			//Handling computers turn
			if game.CurrentTurn == "computer" {
				logHandling.AppendLog("Computer's turn initiated.")
				// You might want to add a short delay here for a better user experience,
				// so the computer's move doesn't happen instantly.
				time.Sleep(time.Second) // Requires "time" package import
				game.ComputerTurn()
				logHandling.AppendLog("Computer's turn completed.")
				game.CurrentTurn = "user" // Switch turn back to the user

				window.Invalidate() // Redraw after the computer's move
			}

			//Checking game end condition
			// This checks if the game has run out of cards in the deck and all players' hands
			if len(game.Deck) == 0 && len(game.UserPlayer.Hand) == 0 && len(game.ComputerPlayer.Hand) == 0 {
				logHandling.AppendLog("\nGame Over! Determining winner...")
				// Determine winner based on collected books
				if game.UserPlayer.NumberOfBooks > game.ComputerPlayer.NumberOfBooks {
					logHandling.AppendLog("You win!")
					choice, err := runEndScreen(window, "You win!")
					if err != nil {
						return "Error", err
					}
					return choice, nil
				} else if game.UserPlayer.NumberOfBooks < game.ComputerPlayer.NumberOfBooks {
					logHandling.AppendLog("Computer wins!")
					choice, err := runEndScreen(window, "Computer wins!")
					if err != nil {
						return "Error", err
					}
					return choice, nil
				} else {
					logHandling.AppendLog("It's a tie!")
					choice, err := runEndScreen(window, "It's a tie!")
					if err != nil {
						return "Error", err
					}
					return choice, nil
				}
			}
			// End game turn / input handling

			layout.Stack{}.Layout(gtx,
				//background
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return imageWidget.Layout(gtx)
				}),
				//current turn text
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						//position on screen, should be about center, and raised above other text
						Top:  unit.Dp(textPosY - 25),
						Left: unit.Dp(textPosX),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						//creating the text
						h2 := material.H2(th, "Computer Books: "+strconv.Itoa(game.ComputerPlayer.NumberOfBooks))
						h2.Color = white
						return h2.Layout(gtx)
					})
				}),
				//current turn text
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						//position on screen, should be about center, and lowered below other text
						Top:  unit.Dp(textPosY + 25), // Adjust these values to move the score
						Left: unit.Dp(textPosX),      // from the top-left corner
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						//creating the text
						h2 := material.H2(th, "User Books: "+strconv.Itoa(game.UserPlayer.NumberOfBooks))
						h2.Color = white
						return h2.Layout(gtx)
					})
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					yOffset := gtx.Constraints.Max.Y / 10

					return layout.Flex{
						Axis:      layout.Horizontal,
						Spacing:   layout.SpaceStart,
						Alignment: layout.Middle,
					}.Layout(gtx,
						//Users cards layout
						ui2.PrintCards(userCards, yOffset, scale, userCardHeight),

						//computers cards layout
						ui2.PrintCards(computerCards, yOffset, scale, computerCardHeight),
					)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return ui2.DrawDynamicButtons(gtx,
						//theme
						th,
						//array of clickable widgets
						clicks,
						//button texts (ranks of cards)
						ranks,
						//button colors
						materialRed, materialBlue,
						//text color
						black)
				}),
				//menu button
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						//position on screen, should be top right corner
						Top:  unit.Dp(10),                        // Adjust these values to move the score
						Left: unit.Dp(gtx.Constraints.Max.X / 2), // Using the X axis /2 to position in the middle
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						//button background and ext
						hitButton := material.Button(th, menuButton, "Menu")
						hitButton.Background = materialYellow
						hitButton.Color = white
						return hitButton.Layout(gtx)
					})

				}),
				//cosole output
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					// Position the console output (e.g., top-left corner)
					return layout.Inset{
						Top:  unit.Dp(10), // 10 Dp from the top
						Left: unit.Dp(10), // 10 Dp from the left
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						logHandling.LogMutex.Lock() // Acquire mutex before reading gameLog
						logText := strings.Join(logHandling.GameLog, "\n")
						logHandling.LogMutex.Unlock() // Release mutex after reading

						lbl := material.Label(th, unit.Sp(14), logText)         // Font size 14sp
						lbl.Color = color.NRGBA{R: 200, G: 200, B: 200, A: 255} // Light gray text
						return lbl.Layout(gtx)
					})
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
