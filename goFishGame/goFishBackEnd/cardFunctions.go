package goFishBackEnd

import (
	"CombinedCardgames/goFishGame/logHandling"
	"fmt"
	"math/rand"
	"strings"
)

// newDeck creates and returns a new deck of 52 cards
func newDeck() []Card {
	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
	suits := []string{"C", "D", "H", "S"}
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

// Deal deals initial hands to players (now a method)
func (g *GoFish) Deal() {
	for i := 0; i < 7; i++ {
		if len(g.Deck) > 0 { // Ensure deck is not empty before dealing
			g.UserPlayer.Hand = append(g.UserPlayer.Hand, g.Deck[0])
			g.Deck = g.Deck[1:]
		} else {
			fmt.Println("Deck ran out during initial deal for user!")
			break
		}
		if len(g.Deck) > 0 { // Ensure deck is not empty before dealing
			g.ComputerPlayer.Hand = append(g.ComputerPlayer.Hand, g.Deck[0])
			g.Deck = g.Deck[1:]
		} else {
			fmt.Println("Deck ran out during initial deal for computer!")
			break
		}
	}
}

// GiveCards transfers all cards of the requested rank from one player to another
func (g *GoFish) GiveCards(rank string, from *Player, to *Player) int {
	count := 0
	var newHand []Card
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

// DrawCard draws the top card from the deck and adds it to the player's hand
func (g *GoFish) DrawCard(player *Player) {
	if len(g.Deck) == 0 {
		fmt.Printf("%s tries to draw, but the deck is empty!\n", player.Name)
		logHandling.AppendLog(player.Name + " tries to draw, but the deck is empty!\n")
		return
	}
	card := g.Deck[0]
	player.Hand = append(player.Hand, card)
	g.Deck = g.Deck[1:]
	fmt.Printf("%s drew a card: %s%s\n", player.Name, card.Rank, card.Suit)
	logHandling.AppendLog(player.Name + " drew a card: " + card.Rank + " " + card.Suit)
}
