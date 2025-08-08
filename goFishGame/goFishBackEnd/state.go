package goFishBackEnd

import (
	"CombinedCardgames/goFishGame/logHandling"

	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Card represents a playing card with a rank and suit
type Card struct {
	Rank string
	Suit string
}

// Player represents a player in the game, holding a hand and completed books
type Player struct {
	Name          string
	Hand          []Card
	Books         []string
	NumberOfBooks int // Corrected to be exported for external access
}

// GoFish encapsulates the entire game state and logic.
type GoFish struct {
	Deck           []Card
	UserPlayer     Player
	ComputerPlayer Player
	CurrentTurn    string // "user" or "computer" - useful for UI
}

// CheckBooks identifies and removes sets of four cards of the same rank (books) from a player's hand
func (g *GoFish) CheckBooks(player *Player) {
	rankCount := make(map[string]int)
	for _, card := range player.Hand {
		rankCount[card.Rank]++
	}
	var newHand []Card
	var bookedRanks []string // Keep track of ranks that formed a book

	for rank, count := range rankCount {
		if count == 4 {
			player.Books = append(player.Books, rank)
			player.NumberOfBooks++ // Increment book count
			fmt.Printf("%s completed a book of %s!\n", player.Name, rank)
			logHandling.AppendLog(player.Name + " completed a book of " + rank)
			bookedRanks = append(bookedRanks, rank)
		}
	}

	// Rebuild hand without the booked cards
	for _, card := range player.Hand {
		found := false
		for _, bookRank := range bookedRanks {
			if card.Rank == bookRank {
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

// PlayerTurn handles the human player's turn, where they ask for a specific rank
func (g *GoFish) PlayerTurn(chosenRank string) { // No longer returns bool for "quit"
	fmt.Printf("You asked for: %s\n", chosenRank) // Log what was asked
	logHandling.AppendLog("You asked for " + chosenRank + "!\n")

	cardsGivenCount := g.GiveCards(chosenRank, &g.ComputerPlayer, &g.UserPlayer)
	if cardsGivenCount == 0 {
		fmt.Println("Go Fish!")
		logHandling.AppendLog("Go Fish!\n")
		g.DrawCard(&g.UserPlayer) // User draws a card
	} else {
		fmt.Printf("You got %d cards!\n", cardsGivenCount)
		logHandling.AppendLog("You got " + strconv.Itoa(cardsGivenCount) + " cards!\n")
	}
	g.CheckBooks(&g.UserPlayer) // Check books for the user
}

// ComputerTurn handles the computer player's turn
func (g *GoFish) ComputerTurn() {
	if len(g.ComputerPlayer.Hand) == 0 {
		g.DrawCard(&g.ComputerPlayer)
		fmt.Println("Computer had no cards, drew one.")
		logHandling.AppendLog("Computer had no cards, drew one.\n")
		g.CheckBooks(&g.ComputerPlayer) // Check if drawing formed a book
		return
	}

	// Computer randomly selects a rank from its hand to ask for
	rank := g.ComputerPlayer.Hand[rand.Intn(len(g.ComputerPlayer.Hand))].Rank
	fmt.Printf("Computer asks for: %s\n", rank)
	logHandling.AppendLog("Computer asks for: " + rank)

	cardsGivenCount := g.GiveCards(rank, &g.UserPlayer, &g.ComputerPlayer)
	if cardsGivenCount == 0 {
		fmt.Println("Go Fish!")
		logHandling.AppendLog("Go Fish!\n")
		g.DrawCard(&g.ComputerPlayer)
	} else {
		fmt.Printf("Computer got %d cards!\n", cardsGivenCount)
		logHandling.AppendLog("Computer got " + strconv.Itoa(cardsGivenCount) + " cards!\n")
	}
	g.CheckBooks(&g.ComputerPlayer)
}

// NewGame initializes and returns a new GoFish game instance.
func NewGame() (*GoFish, error) {
	// Initialize the GoFish struct
	game := &GoFish{
		UserPlayer:     Player{Name: "You"},
		ComputerPlayer: Player{Name: "Computer"},
		CurrentTurn:    "user", // User starts first
	}

	rand.Seed(time.Now().UnixNano())

	// Create and shuffle the deck, then assign it to the game instance
	game.Deck = newDeck()
	shuffle(game.Deck) // shuffle operates on a slice, so pass game.Deck directly

	// Deal cards using the game's method
	game.Deal()

	fmt.Println("Go Fish game started!")

	return game, nil // Return the initialized game instance
}
