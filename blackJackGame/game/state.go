package game

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var filename = "gameScore.txt"

type Game struct {
	UserHand    []Card
	DealerHand  []Card
	CurrentShoe []Card

	UserInput    string
	UserPoints   int
	Wins         int
	Losses       int
	CardCounter  int
	GameEnded    bool //true when game ends
	RevealDealer bool //reveals dealers hidden card when its dealers turn
	UserWin      bool // true = user, false = dealer
	TieGame      bool //incase game is tied
	BlackJack    bool //incase of blackjack
}

type Card struct {
	Rank string
	Suit string
}

// saves points,wins,losses to file
func saveStats(filename string, userPoints, wins, losses int) error {
	// Format the values as a string like "12,5,7"
	content := fmt.Sprintf("%d,%d,%d", userPoints, wins, losses)

	// Write to the file
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil

}

// NewGame initializes a new Game instance.
func NewGame() (*Game, error) {
	game := &Game{}

	//default points/wins/losses
	defaultContent := "0,0,0"

	//checks for file "gameScore.txt", if it doesnt exist then it is created
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.WriteFile(filename, []byte(defaultContent), 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to create file: %w", err)
			}
			data = []byte(defaultContent)
		} else {
			return nil, fmt.Errorf("error reading file: %w", err)
		}
	}

	//checks file for 3 fields
	parts := strings.Split(strings.TrimSpace(string(data)), ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid file format. Expected format: x,y,z")
	}

	//grabs point total
	game.UserPoints, err = strconv.Atoi(parts[0]) // Fixed index
	if err != nil {
		return nil, fmt.Errorf("error parsing userPoints: %w", err)
	}
	//grabs wins total
	game.Wins, err = strconv.Atoi(parts[1]) // Fixed index
	if err != nil {
		return nil, fmt.Errorf("error parsing wins: %w", err)
	}
	//grabs losses total
	game.Losses, err = strconv.Atoi(parts[2]) // Fixed index
	if err != nil {
		return nil, fmt.Errorf("error parsing losses: %w", err)
	}

	//creates a shoe for use in game
	game.CurrentShoe = generateShoe()
	game.shuffleShoe()

	return game, nil
}

// PlayGame initial state of game, has 2 switch cases, 1st for play game, score, menu, 2nd for hit, stand, double down
func (g *Game) PlayGame(action string) {

	switch action {
	case "1": //starts game

		g.dealCards()
		g.GameEnded = false
		if g.CardCounter > 195 {
			g.shuffleShoe()
		}

		if checkSum(g.UserHand) == 21 && checkSum(g.DealerHand) == 21 {
			//if user and dealer are both dealt 21 right away, its a tie
			err := saveStats(filename, g.UserPoints, g.Wins, g.Losses)
			if err != nil {
				return
			}

			g.TieGame = true
			g.GameEnded = true
			g.PlayGame("")
		} else if checkSum(g.UserHand) == 21 && checkSum(g.DealerHand) != 21 {
			//when user is dealt 21 right away,it is an instant win
			g.Wins += 1
			g.UserPoints += 150
			err := saveStats(filename, g.UserPoints, g.Wins, g.Losses)
			if err != nil {
				return
			}

			g.printCards()
			g.BlackJack = true
			g.UserWin = true //user wins
			g.GameEnded = true
			g.PlayGame("")
		}

	case "2": //score board
		return

	case "3": //exits program
		err := saveStats(filename, g.UserPoints, g.Wins, g.Losses)
		if err != nil {
			return
		}
		os.Exit(0)
	}

}

func (g *Game) PhaseOne(userGameInput string) {

	//Hit
	if userGameInput == "1" {

		g.UserHand = append(g.UserHand, (g.CurrentShoe)[0])
		g.CardCounter += 1

		g.CurrentShoe = (g.CurrentShoe)[1:] // Remove first
		g.printCards()
		g.HitUser("")
	}

	//Stand
	if userGameInput == "2" {
		g.dealerTurn(false)
	}

	//Double down
	if userGameInput == "3" {

		g.UserHand = append(g.UserHand, (g.CurrentShoe)[0])
		g.CardCounter += 1

		g.CurrentShoe = (g.CurrentShoe)[1:] // Remove first

		g.printCards()
		total := checkSum(g.UserHand)

		if total > 21 {
			//busted on double down, double the points lost
			g.Losses += 1
			g.UserPoints -= 100
			err := saveStats(filename, g.UserPoints, g.Wins, g.Losses)
			if err != nil {
				return
			}

			g.GameEnded = true

			g.PlayGame("")
		}
		g.dealerTurn(true)
	}
}

// HitUser prompts user after using if statements to determine if they lost
func (g *Game) HitUser(userInput string) {
	//users total
	var total = checkSum(g.UserHand)

	if total > 21 {
		//user busts on the draw
		g.Losses += 1
		g.UserPoints -= 50
		err := saveStats(filename, g.UserPoints, g.Wins, g.Losses)
		if err != nil {
			return
		}

		g.GameEnded = true

		g.PlayGame("")
	} else if total < 21 {
		//since user hasnt busted, user can choose to hit or stand

		if userInput == "1" {
			//if user choses to hit again, gives them a card and recalls HitUser()
			g.UserHand = append(g.UserHand, (g.CurrentShoe)[0])
			g.CardCounter += 1

			g.CurrentShoe = (g.CurrentShoe)[1:] // Remove first

			g.printCards()

			g.HitUser("")
		}
		if userInput == "2" {
			//if user chooses to stand, it becomes dealers turn
			g.dealerTurn(false)
		}
	} else {
		//if user has 21, auto dealers turn
		g.dealerTurn(false)
	}
}

// dealers turn, user has hit as many as they wanted/could, or chose to stand
// if user double downed, the points double
func (g *Game) dealerTurn(doubleDown bool) {
	g.RevealDealer = true
	g.printCards()

	//dealer and users total
	var total = checkSum(g.DealerHand)
	var userTotal = checkSum(g.UserHand)

	var win = 100
	var loss = 50

	//doubledown doubles win or loss point acquisition
	if doubleDown {
		win *= 2
		loss *= 2
	}

	if total > 21 {
		//if dealer goes over 21, dealer busts and loses
		g.Wins += 1
		g.UserPoints += win

		err := saveStats(filename, g.UserPoints, g.Wins, g.Losses)
		if err != nil {
			return
		}

		g.UserWin = true
		g.GameEnded = true

		g.PlayGame("")
	} else if total >= 17 && total > userTotal {
		//dealer only required to draw up to 17,
		//since its over 17, dealer stands and wins the game
		g.Losses += 1
		g.UserPoints -= loss

		err := saveStats(filename, g.UserPoints, g.Wins, g.Losses)
		if err != nil {
			return
		}

		g.GameEnded = true

		g.PlayGame("")
	} else if total >= 17 && total < userTotal {
		//Dealer cant draw at or over 17
		//So if dealer has less than user, user wins
		g.Wins += 1
		g.UserPoints += win

		err := saveStats(filename, g.UserPoints, g.Wins, g.Losses)
		if err != nil {
			return
		}

		g.UserWin = true
		g.GameEnded = true

		g.PlayGame("")
	} else if total <= 17 {
		//dealer must draw if they have less than 17
		//dealer draws here and recursivley calls the same function
		g.DealerHand = append(g.DealerHand, (g.CurrentShoe)[0])
		g.CardCounter += 1

		//removes front card
		g.CurrentShoe = (g.CurrentShoe)[1:] // Remove first

		g.dealerTurn(doubleDown)

	} else if total == userTotal /*&& total >= 17*/ {
		// if they are tied, and its after the other if statements

		err := saveStats(filename, g.UserPoints, g.Wins, g.Losses)
		if err != nil {
			return
		}

		g.TieGame = true
		g.GameEnded = true

		g.PlayGame("")
	}

}

// checks the sum of a hand
func checkSum(currentHand []Card) int {
	var total = 0
	var aceCount = 0

	// Calculate total hand value, tracking number of Aces
	for _, card := range currentHand {
		// Direct lookup using card.Rank, no string parsing needed
		val := CardValueMap[card.Rank]
		total += val
		if card.Rank == "A" { // Check for Ace using its rank
			aceCount++
		}
	}

	// Adjusts Aces from 11 to 1 as needed
	for total > 21 && aceCount > 0 {
		total -= 10
		aceCount--
	}

	return total
}

// CheckGameEnded Checks if the game ended, resets the flag if true
func (g *Game) CheckGameEnded() bool {
	if g.GameEnded {
		g.GameEnded = false
		return true
	}
	return false
}

// CheckRevealDealer Checks if the dealers card can be revealed, resets the flag if true
func (g *Game) CheckRevealDealer() bool {
	if g.RevealDealer {
		g.RevealDealer = false
		return true
	}
	return false
}

// CheckUserWin Checks if the user won, resets the flag if true
func (g *Game) CheckUserWin() bool {
	if g.UserWin {
		g.UserWin = false
		return true
	}
	return false
}

// CheckTieGame Checks if the game ended in a tie, resets the flag if true
func (g *Game) CheckTieGame() bool {
	if g.TieGame {
		g.TieGame = false
		return true
	}
	return false
}

// CheckBlackJack Checks if the user got a black jack, resets the flag if true
func (g *Game) CheckBlackJack() bool {
	if g.BlackJack {
		g.BlackJack = false
		return true
	}
	return false
}
