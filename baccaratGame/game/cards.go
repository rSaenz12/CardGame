//This file defines the "baackend" functions of the cards. how they behavior and are manipulated. 

package game

import (
	"math/rand"
)

// Card represents a single playing card.
type Card struct {
	Rank string
	Suit string
}

// newShoe creates a collection of decks.
func newShoe(numDecks int) []Card {
	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
	suits := []string{"C", "D", "H", "S"}
	var shoe []Card
	for i := 0; i < numDecks; i++ {
		for _, suit := range suits {
			for _, rank := range ranks {
				shoe = append(shoe, Card{Rank: rank, Suit: suit})
			}
		}
	}
	return shoe
}

// shuffle randomizes the order of cards.
func shuffle(shoe []Card) {
	rand.Shuffle(len(shoe), func(i, j int) {
		shoe[i], shoe[j] = shoe[j], shoe[i]
	})
}

// drawCard removes and returns the top card from the stack called "shoe".
func (g *Game) drawCard() Card {
	card := g.Shoe[0]
	g.Shoe = g.Shoe[1:]
	return card
}
