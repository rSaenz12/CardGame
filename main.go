package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var filename = "gameScore.txt"

func main() {

	defaultContent := "0,0,0"

	// Try to read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, so create it with default content
			err := os.WriteFile(filename, []byte(defaultContent), 0644)
			if err != nil {
				fmt.Println("Failed to create file:", err)
				return
			}
			fmt.Println("File created with default content:", defaultContent)
			data = []byte(defaultContent) // Use this for parsing below
		} else {
			fmt.Println("Error reading file:", err)
			return
		}
	}

	// Convert the file content into 3 integers
	parts := strings.Split(strings.TrimSpace(string(data)), ",")
	if len(parts) != 3 {
		fmt.Println("Invalid file format. Expected format: x,y,z")
		return
	}

	// Parse the integers
	var userPoints, wins, losses int
	userPoints, err = strconv.Atoi(parts[0])
	if err != nil {
		fmt.Println("Error parsing first value:", err)
		return
	}
	wins, err = strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("Error parsing second value:", err)
		return
	}
	losses, err = strconv.Atoi(parts[2])
	if err != nil {
		fmt.Println("Error parsing third value:", err)
		return
	}

	var cardCounter = 0
	shoe := generateShoe()
	shuffleShoe(shoe)

	//var userPoints int = 0

	playGame(&shoe, &cardCounter, &userPoints, &wins, &losses)
	saveStats(filename, userPoints, wins, losses)

}
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

func playGame(shoe *[]string, cardCounter *int, userPoints *int, wins *int, losses *int) {
	//NOTE INSTANT BLACK JACK MUST BE A WIN RIGTH AWAY

	var userInput = ""
	//for userInput != "1" || userInput != "2" {
	for {

		fmt.Println()

		fmt.Println("Type 1 to Deal a game, 2 to view score, 3 to exit")
		_, err := fmt.Scan(&userInput)
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return // Exit or handle the error appropriately
		}

		if userInput == "1" {
			var dealer []string
			var user []string
			dealer, user = dealCards(shoe, cardCounter)
			if *cardCounter > 195 {
				shuffleShoe(*shoe)
			}

			if checkSum(user) == 21 {
				fmt.Println("BlackJack, You win")
				*wins += 1
				*userPoints += 100
				printWinsLosses(wins, losses)
				printPoints(userPoints)
				saveStats(filename, *userPoints, *wins, *losses)

				playGame(shoe, cardCounter, userPoints, wins, losses)
			}

			var userGameInput = ""
			fmt.Println("Would you like to hit or stand ?")
			fmt.Println("Type 1 to hit, 2 to stand")
			_, err := fmt.Scan(&userGameInput)
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				return // Exit or handle the error appropriately
			}

			if userGameInput == "1" {
				hitUser(shoe, &user, &dealer, cardCounter, userPoints, wins, losses)
			}
			if userGameInput == "2" {
				dealerTurn(shoe, &dealer, user, cardCounter, userPoints, wins, losses)
			}

		}
		if userInput == "2" {
			for userInput != "1" {
				printWinsLosses(wins, losses)
				printPoints(userPoints)

				fmt.Println("Type 1 to return to menu")
				_, err := fmt.Scan(&userInput)
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					return // Exit or handle the error appropriately
				}
			}
		}
		if userInput == "3" {
			saveStats(filename, *userPoints, *wins, *losses)
			os.Exit(0)
		}

	}
}
func dealCards(shoe *[]string, cardCounter *int) ([]string, []string) {
	var dealer []string
	var user []string

	for i := 0; i < 4; i++ {
		if i%2 == 0 {
			user = append(user, (*shoe)[i])
			*cardCounter += 1

			*shoe = (*shoe)[1:] // Remove first

		}
		if i%2 == 1 {
			dealer = append(dealer, (*shoe)[i])
			*cardCounter += 1

			*shoe = (*shoe)[1:] // Remove first

		}
	}

	printCards(dealer, user, false)
	var userTotal = checkSum(user)
	fmt.Println("Your total is " + strconv.Itoa(userTotal))
	fmt.Println()

	return dealer, user
}
func printCards(dealer []string, user []string, revealDealer bool) {
	// Print dealer's hand
	fmt.Print("Dealer: ")
	if revealDealer {
		for i, card := range dealer {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(card)
		}
	} else {
		fmt.Print(dealer[0] + ", hidden card")
	}
	fmt.Println()

	// Print user's hand
	fmt.Print("User: ")
	for i, card := range user {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(card)
	}
	fmt.Println()
	fmt.Println()
}
func printPoints(userPoints *int) {
	fmt.Println("Points: " + strconv.Itoa(*userPoints))
	fmt.Println()
}
func printWinsLosses(wins *int, losses *int) {
	fmt.Println("Wins: " + strconv.Itoa(*wins) + " Losses: " + strconv.Itoa(*losses))
}

func hitUser(shoe *[]string, user *[]string, dealer *[]string, cardCounter *int, userPoints *int, wins *int, losses *int) {

	*user = append(*user, (*shoe)[0])
	*cardCounter += 1

	*shoe = (*shoe)[1:] // Remove first

	printCards(*dealer, *user, false)

	var total = checkSum(*user)
	if total > 21 {
		fmt.Println("Your total is " + strconv.Itoa(total))
		fmt.Println("You bust")
		*losses += 1
		*userPoints -= 50
		printWinsLosses(wins, losses)
		printPoints(userPoints)
		saveStats(filename, *userPoints, *wins, *losses)

		playGame(shoe, cardCounter, userPoints, wins, losses)
	}

	if total < 21 {
		fmt.Println("Your total is " + strconv.Itoa(total))
		var userInput = ""
		fmt.Println("Type 1 to hit, 2 to stand")
		for {
			_, err := fmt.Scan(&userInput)
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				return // Exit or handle the error appropriately
			}
			fmt.Println()
			if userInput == "1" {
				hitUser(shoe, user, dealer, cardCounter, userPoints, wins, losses)
			}
			if userInput == "2" {
				dealerTurn(shoe, dealer, *user, cardCounter, userPoints, wins, losses)
			}
		}
	}
	if total == 21 {
		fmt.Println("21!")
		fmt.Println("Dealers turn")
		fmt.Println()
		dealerTurn(shoe, dealer, *user, cardCounter, userPoints, wins, losses)

	}

}

func dealerTurn(shoe *[]string, dealer *[]string, user []string, cardCounter *int, userPoints *int, wins *int, losses *int) {

	printCards(*dealer, user, true)
	var total = checkSum(*dealer)
	var userTotal = checkSum(user)

	if total > 21 {
		fmt.Println("You win, Dealer busts with " + strconv.Itoa(total) + " !")
		*wins += 1
		*userPoints += 100
		printWinsLosses(wins, losses)
		printPoints(userPoints)
		saveStats(filename, *userPoints, *wins, *losses)

		playGame(shoe, cardCounter, userPoints, wins, losses)
	} else if total >= 17 && total > userTotal {
		fmt.Println("Dealer wins with " + strconv.Itoa(total) + " to beat User's " + strconv.Itoa(userTotal) + "!")
		*losses += 1
		*userPoints -= 50
		printWinsLosses(wins, losses)
		printPoints(userPoints)
		saveStats(filename, *userPoints, *wins, *losses)

		playGame(shoe, cardCounter, userPoints, wins, losses)
	} else if total < 17 {
		*dealer = append(*dealer, (*shoe)[0])
		*cardCounter += 1

		//removes front card
		*shoe = (*shoe)[1:] // Remove first
		fmt.Println()
		dealerTurn(shoe, dealer, user, cardCounter, userPoints, wins, losses)

	} else if total < userTotal {
		*dealer = append(*dealer, (*shoe)[0])
		*cardCounter += 1

		//removes front card
		*shoe = (*shoe)[1:] // Remove first
		fmt.Println()
		dealerTurn(shoe, dealer, user, cardCounter, userPoints, wins, losses)

	} else if total == userTotal && total > 17 {
		fmt.Println("It's a tie " + strconv.Itoa(total) + " to " + strconv.Itoa(userTotal) + " !")
		printWinsLosses(wins, losses)
		printPoints(userPoints)
		saveStats(filename, *userPoints, *wins, *losses)

		playGame(shoe, cardCounter, userPoints, wins, losses)
	}

}

// 13 cards added to suit, 4 per deck
func generateSuit() []string {
	return []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "King", "Queen", "Ace"}
}

// creating a deck using the suits created, 4 suits per deck
func generateDeck() []string {
	suits := generateSuit()
	var deck []string

	for i := 0; i < 4; i++ {
		deck = append(deck, suits...)
	}

	return deck
}

func generateShoe() []string {
	deck := generateDeck()
	var shoe []string
	for i := 0; i < 5; i++ {
		shoe = append(shoe, deck...)
	}
	return shoe
}

func shuffleShoe(shoe []string) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range shoe {
		j := r.Intn(i + 1)
		shoe[i], shoe[j] = shoe[j], shoe[i]
	}
}

//use 5 decks
//shuffle decks after assembly
//shuffle again after 195 cards are used, announce to user that its being reshuffled

func checkSum(currentHand []string) int {

	var total = 0
	var aceCount = 0
	values := map[string]int{
		"2": 2, "3": 3, "4": 4, "5": 5,
		"6": 6, "7": 7, "8": 8, "9": 9,
		"10": 10, "Jack": 10, "Queen": 10, "King": 10,
		"Ace": 11,
	}

	for _, card := range currentHand {
		val := values[card]
		total += val
		if card == "Ace" {
			aceCount++
		}
	}

	// Adjust Aces from 11 to 1 as needed
	for total > 21 && aceCount > 0 {
		total -= 10
		aceCount--
	}

	return total

}
