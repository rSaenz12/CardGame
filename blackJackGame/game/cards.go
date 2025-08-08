package game

import (
	"fmt"
	"math/rand"
	"time"
)

var CardValueMap = map[string]int{
	"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
	"T": 10, "J": 10, "Q": 10, "K": 10,
	"A": 11, // Ace starts as 11, adjusted later if needed
}

// shuffle decks after assembly
// shuffle again after 195 cards are used, announce to user that its being reshuffled
func (g *Game) dealCards() {

	//emptying user and dealers hands
	g.UserHand = []Card{}
	g.DealerHand = []Card{}

	// dealing user and dealer
	for i := 0; i < 4; i++ {
		//user gets cards 1, 3
		if i%2 == 1 {
			//adds card to hand
			g.UserHand = append(g.UserHand, (g.CurrentShoe)[0])
			g.CardCounter += 1

			//removes top card
			g.CurrentShoe = (g.CurrentShoe)[1:]
		}
		//dealer gets cards 2,4
		if i%2 == 0 {
			//adds card to hand
			g.DealerHand = append(g.DealerHand, (g.CurrentShoe)[0])
			g.CardCounter += 1

			//removes top card
			g.CurrentShoe = (g.CurrentShoe)[1:]
		}
	}

}

// prints cards
func (g *Game) printCards() {
	// Print dealer's hand
	fmt.Print("Dealer: ")

	//keeps dealers second card hidden until dealers turn
	if g.CheckRevealDealer() {
		for i, card := range g.DealerHand {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(card)
		}
	} else {
		if len(g.DealerHand) > 0 {
			fmt.Print(g.DealerHand[0].Rank + ", hidden card")
		} else {
			fmt.Print("No dealer cards dealt yet.") // Or some other message
		}
	}
	fmt.Println()

	// Print user hand
	fmt.Print("User: ")
	if len(g.UserHand) > 0 {
		for i, card := range g.UserHand {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(card)
		}
	} else {
		fmt.Print("No user cards dealt yet.")
	}
	fmt.Println()
	fmt.Println()
}

// 13 cards added to suit, 4 per deck
func generateSuit() []string {
	return []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "K", "Q", "A"}
}

// creating a deck using the suits created, 4 suits per deck
func generateDeck() []Card {
	unSuited := generateSuit()
	suits := []string{"C", "D", "H", "S"}
	var deck []Card

	//outer loop: will go 2,3,4..., A, inner goes C,D,H,S
	//loops together go 2C,2D,2H,2S,3C,3D...,AH,AS
	//creates a fully suited deck
	for _, suit := range suits {
		for _, rank := range unSuited {
			deck = append(deck, Card{Rank: rank, Suit: suit})
		}
	}
	return deck
}

// use 5 decks to make  shoe
func generateShoe() []Card {
	deck := generateDeck()
	var shoe []Card

	//appends 5 decks together
	for i := 0; i < 5; i++ {
		shoe = append(shoe, deck...)
	}
	return shoe
}

func (g *Game) shuffleShoe() {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range g.CurrentShoe {
		j := r.Intn(i + 1)
		g.CurrentShoe[i], g.CurrentShoe[j] = g.CurrentShoe[j], g.CurrentShoe[i]
	}
}
