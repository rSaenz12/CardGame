package main

import (
	"fmt"
	"math/rand"
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
	gameLoop(&deck, &player, &computer)
}
