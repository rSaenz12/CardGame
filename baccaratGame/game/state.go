// This file manages the game state across turns and delivers structures to save score, and the progression of the deck of cards. 
package game

import (
	"log"
	"strconv"
)

// Game holds the state for the Baccarat game.
type Game struct {
	Shoe          []Card
	PlayerHand    []Card
	BankerHand    []Card
	UserBet       string // "player", "banker", "tie"
	UserPoints    int
	LastResult    string // "Player Wins", "Banker Wins", "Tie"
	CurrentPhase  string // "betting", "result"
	CardsRevealed bool
}

// NewGame creates a new Baccarat game instance.
func NewGame() (*Game, error) {
	game := &Game{
		UserPoints:    1000, // Starting points
		CurrentPhase:  "betting",
		CardsRevealed: false,
	}
	game.Shoe = newShoe(6) // Baccarat is usually played with 6 decks
	shuffle(game.Shoe)
	// Pre-deal 
	game.PlayerHand = make([]Card, 2)
	game.BankerHand = make([]Card, 2)
	return game, nil
}

// PlaceBet sets the user's bet for the round.
func (g *Game) PlaceBet(betType string) {
	if g.CurrentPhase == "betting" {
		g.UserBet = betType
	}
}

// DealHand plays the round
func (g *Game) DealHand() {
	if g.UserBet == "" {
		g.LastResult = "Please place a bet first!"
		return
	}
	g.CurrentPhase = "result"
	g.CardsRevealed = true // Reveal the cards

	// Reset hands and deal new cards
	g.PlayerHand = []Card{}
	g.BankerHand = []Card{}

	// Initial deal
	g.PlayerHand = append(g.PlayerHand, g.drawCard())
	g.BankerHand = append(g.BankerHand, g.drawCard())
	g.PlayerHand = append(g.PlayerHand, g.drawCard())
	g.BankerHand = append(g.BankerHand, g.drawCard())

	playerScore := getHandValue(g.PlayerHand)
	bankerScore := getHandValue(g.BankerHand)

	// Check for natural win
	if playerScore >= 8 || bankerScore >= 8 {
		g.determineWinner()
		return
	}

	// Player's third card rule
	drewThirdPlayer := false
	var playerThirdCardValue int
	if playerScore <= 5 {
		thirdCard := g.drawCard()
		g.PlayerHand = append(g.PlayerHand, thirdCard)
		playerThirdCardValue = getCardValue(thirdCard)
		drewThirdPlayer = true
	}

	// Banker's third card rule
	if drewThirdPlayer {
		if (bankerScore <= 2) ||
			(bankerScore == 3 && playerThirdCardValue != 8) ||
			(bankerScore == 4 && playerThirdCardValue >= 2 && playerThirdCardValue <= 7) ||
			(bankerScore == 5 && playerThirdCardValue >= 4 && playerThirdCardValue <= 7) ||
			(bankerScore == 6 && playerThirdCardValue >= 6 && playerThirdCardValue <= 7) {
			g.BankerHand = append(g.BankerHand, g.drawCard())
		}
	} else if bankerScore <= 5 {
		g.BankerHand = append(g.BankerHand, g.drawCard())
	}

	g.determineWinner()
}

// determineWinner calculates the result and updates points.
func (g *Game) determineWinner() {
	playerScore := getHandValue(g.PlayerHand)
	bankerScore := getHandValue(g.BankerHand)

	winner := ""
	if playerScore > bankerScore {
		winner = "player"
		g.LastResult = "Player Wins! " + strconv.Itoa(playerScore) + " to " + strconv.Itoa(bankerScore)
	} else if bankerScore > playerScore {
		winner = "banker"
		g.LastResult = "Banker Wins! " + strconv.Itoa(bankerScore) + " to " + strconv.Itoa(playerScore)
	} else {
		winner = "tie"
		g.LastResult = "Tie! Both have " + strconv.Itoa(playerScore)
	}

	// Update points
	betAmount := 100 // Fixed bet for now
	if g.UserBet == winner {
		if winner == "banker" {
			g.UserPoints += int(float64(betAmount) * 0.95) // Banker pays 1:1 with 5% commission
		} else if winner == "tie" {
			g.UserPoints += betAmount * 8 // Tie pays 8:1
		} else {
			g.UserPoints += betAmount // Player pays 1:1
		}
	} else {
		g.UserPoints -= betAmount
	}
}

// NewRound resets the game for the next round of betting.
func (g *Game) NewRound() {
	if len(g.Shoe) < 20 { // Reshuffle when deck is low
		g.Shoe = newShoe(6)
		shuffle(g.Shoe)
	}
	g.CurrentPhase = "betting"
	g.UserBet = ""
	g.LastResult = ""
	g.CardsRevealed = false
	// Pre-deal two face-down cards to each hand for initial display
	g.PlayerHand = make([]Card, 2)
	g.BankerHand = make([]Card, 2)
}

// getCardValue returns the Baccarat value of a card.
func getCardValue(card Card) int {
	switch card.Rank {
	case "A":
		return 1
	case "2", "3", "4", "5", "6", "7", "8", "9":
		val, err := strconv.Atoi(card.Rank)
		if err != nil {
			log.Printf("Error converting card rank to int: %v", err)
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
}
