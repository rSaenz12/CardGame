package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Card represents a playing card with a rank and suit
type Card struct {
	Rank string
	Suit string
}

// Player represents a player in the game, holding a hand and completed books
type Player struct {
	Name  string
	Hand  []Card
	Books []string
}

// newDeck creates and returns a new deck of 52 cards
func newDeck() []Card {
	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	suits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	var deck []Card
	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{Rank: rank, Suit: suit})
		}
	}
	return deck
}

// shuffle randomizes the order of cards in the deck
func shuffle(deck []Card) {
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

// deal distributes 7 cards to each player from the deck
func deal(deck *[]Card, player1 *Player, player2 *Player) {
	for i := 0; i < 7; i++ {
		player1.Hand = append(player1.Hand, (*deck)[0])
		*deck = (*deck)[1:]
		player2.Hand = append(player2.Hand, (*deck)[0])
		*deck = (*deck)[1:]
	}
}

// giveCards transfers all cards of the requested rank from one player to another
func giveCards(rank string, from *Player, to *Player) int {
	count := 0
	newHand := []Card{}
	for _, card := range from.Hand {
		if strings.EqualFold(card.Rank, rank) {
			to.Hand = append(to.Hand, card)
			count++
		} else {
			newHand = append(newHand, card)
		}
	}
	from.Hand = newHand
	return count
}

// drawCard draws the top card from the deck and adds it to the player's hand
func drawCard(deck *[]Card, player *Player) {
	if len(*deck) == 0 {
		return
	}
	player.Hand = append(player.Hand, (*deck)[0])
	*deck = (*deck)[1:]
}

// checkBooks checks if the player has any books (4 of the same rank) and records them
func checkBooks(player *Player) {
	rankCount := make(map[string]int)
	for _, card := range player.Hand {
		rankCount[card.Rank]++
	}
	newHand := []Card{}
	for rank, count := range rankCount {
		if count == 4 {
			player.Books = append(player.Books, rank)
			fmt.Printf("%s completed a book of %s!\n", player.Name, rank)
		}
	}
	for _, card := range player.Hand {
		found := false
		for _, book := range player.Books {
			if card.Rank == book {
				found = true
				break
			}
		}
		if !found {
			newHand = append(newHand, card)
		}
	}
	player.Hand = newHand
}

// playerTurn handles the human player's turn, including the option to quit
func playerTurn(deck *[]Card, player *Player, opponent *Player) bool {
	fmt.Println("Your hand:")
	for _, card := range player.Hand {
		fmt.Printf("%s ", card.Rank)
	}
	fmt.Println()
	fmt.Print("Ask for a rank (or type 'quit' to exit): ")
	var rank string
	fmt.Scanln(&rank)
	if strings.EqualFold(rank, "quit") {
		return true // signal to quit
	}
	if giveCards(rank, opponent, player) == 0 {
		fmt.Println("Go Fish!")
		drawCard(deck, player)
	} else {
		fmt.Println("You got cards!")
	}
	checkBooks(player)
	return false
}

// computerTurn handles the computer's turn, randomly asking for a rank in its hand
func computerTurn(deck *[]Card, computer *Player, opponent *Player) {
	if len(computer.Hand) == 0 {
		drawCard(deck, computer)
		return
	}
	rank := computer.Hand[rand.Intn(len(computer.Hand))].Rank
	fmt.Printf("Computer asks for: %s\n", rank)
	if giveCards(rank, opponent, computer) == 0 {
		fmt.Println("Go Fish!")
		drawCard(deck, computer)
	} else {
		fmt.Println("Computer got cards!")
	}
	checkBooks(computer)
}

// gameLoop runs the main game, alternating turns until the game ends or the player quits
func gameLoop(deck *[]Card, player *Player, computer *Player) {
	for len(*deck) > 0 || len(player.Hand) > 0 || len(computer.Hand) > 0 {
		if len(player.Hand) > 0 {
			fmt.Println("\nYour turn!")
			if playerTurn(deck, player, computer) {
				fmt.Println("You quit the game.")
				return
			}
		}
		if len(*deck) == 0 && len(player.Hand) == 0 && len(computer.Hand) == 0 {
			break
		}
		if len(computer.Hand) > 0 {
			fmt.Println("\nComputer's turn!")
			computerTurn(deck, computer, player)
		}
	}
	fmt.Println("\nGame over!")
	fmt.Printf("Your books: %v\n", player.Books)
	fmt.Printf("Computer's books: %v\n", computer.Books)
	if len(player.Books) > len(computer.Books) {
		fmt.Println("You win!")
	} else if len(player.Books) < len(computer.Books) {
		fmt.Println("Computer wins!")
	} else {
		fmt.Println("It's a tie!")
	}
}

// main initializes the game and starts the game loop
func main() {
	rand.Seed(time.Now().UnixNano())
	deck := newDeck()
	shuffle(deck)
	player := Player{Name: "You"}
	computer := Player{Name: "Computer"}
	deal(&deck, &player, &computer)
	fmt.Println("Go Fish game started!")

	// Starting the Gio GUI in a separate goroutine
	go func() {
		//main application window
		window := new(app.Window)

		// Pass the window and the initialized game state to your UI's Run function
		if err := Run(window, &player); err != nil { // Pass currentGame here
			log.Fatal(err)
		}
		os.Exit(0) // Exit the program gracefully when the window closes
	}()

	//Gio application's main loop
	app.Main()

	gameLoop(&deck, &player, &computer)

}

func Run(window *app.Window, player *Player) error {

	var ops op.Ops

	background := LoadImage("tableTop.png")
	bgOp := paint.NewImageOp(background)

	var userCards []image.Image
	var userHand []string

	for _, card := range player.Hand {
		userHand = append(userHand, card.Rank)
	}

	// Declare imageWidget once outside the loop for correct usage
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

			userCards = make([]image.Image, len(player.Hand))
			userCards = GetCardImage(userHand)

			//scale for objects,
			scale := float32(0.25) // 25% size

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
						}),
					)
				}),
			)
			e.Frame(gtx.Ops)
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
