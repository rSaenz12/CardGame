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

// drawCard removes and returns the top card from the shoe.
func (g *Game) drawCard() Card {
	card := g.Shoe[0]
	g.Shoe = g.Shoe[1:]
	return card
}

/*
// getCardValue returns the Baccarat value of a card.
func getCardValue(card Card) int {
	switch card.Rank {
	case "A":
		return 1
	case "2", "3", "4", "5", "6", "7", "8", "9":
		val := 0
		_, err := fmt.Sscanf(card.Rank, "%d", &val)
		if err != nil {
			return 0
		}
		return val
	case "T", "J", "Q", "K":
		return 0
	default:
		return 0
	}
}

// getHandValue returns the total Baccarat value of a hand.
func getHandValue(hand []Card) int {
	value := 0
	for _, card := range hand {
		value += getCardValue(card)
	}
	return value % 10
}*/
