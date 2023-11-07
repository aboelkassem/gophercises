package main

import (
	"fmt"
	"github/aboelkassem/gophercises/blackjack/deck"
	"log"
	"strings"
)

type PlayType int

const (
	PlayHit PlayType = iota
	PlayStand
)

func main() {
	cards := deck.New(deck.OptionShuffle())

	players := []player{
		&humanPlayer{},
		&dealerPlayer{},
	}

	deal := func(p player) {
		var c deck.Card
		c, cards = cards[0], cards[1:]
		p.Deal(c)
	}
	// Round is Label for dealing with looping (continue, break)
Game:
	for {
		// 1- Every player is dealt 2 cards
		for i := 1; i <= 2; i++ {
			for _, p := range players {
				deal(p)
			}
		}

		// 2. The player's turn
		// 3. The dealer's turn
	Round:
		for r := 1; ; r++ {
			fmt.Println("Round", r)

			// print players cards
			for _, p := range players {
				// // if player interface is type of humanPlayer
				// if _, ok := p.(*humanPlayer); ok {
				// 	fmt.Print("Player: ")
				// } else {
				// 	fmt.Print("Dealer: ")
				// }
				fmt.Printf("%s: %v, Score: %d\n", p.Name(), p.Hand(), calcScore(p.Hand()))
			}

			// 4. Determining the winner
			// TODO: Fix to handling multiple players
			for i := 0; i < len(players); i++ {
				p := players[i]
				s := calcScore(p.Hand())
				if s == 21 {
					fmt.Printf("%s won!\n", p.Name())
					break Game
				}

				if s > 21 {
					otherPlayer := players[(i+1)%len(players)]
					fmt.Printf("%s disqualified!\n", otherPlayer.Name())
					break Game
				}
			}

			for i := 0; i < len(players); {
				p := players[i]

				// if player interface is type of humanPlayer
				if _, ok := p.(*humanPlayer); ok {
					fmt.Printf("Player %d's turn\n", i+1)
				} else {
					fmt.Printf("Dealer's turn\n")
					if r == 1 {
						continue
					}
				}

				input, err := p.Play()

				if err != nil {
					log.Fatal(err)
					continue
				}

				fmt.Printf("Score: %d\n", calcScore(p.Hand()))

				switch input {
				case PlayHit:
					deal(p)
					continue Round
				case PlayStand:
					i++
				}
			}
		}
	}
}

func calcScore(hand []deck.Card) int {
	score := 0
	var aces int
	for _, c := range hand {
		switch { // switch true, all cases must return bool
		case deck.ValueTwo <= c.Value && c.Value <= deck.ValueTen:
			score += int(c.Value)
		case deck.ValueJack <= c.Value && c.Value <= deck.ValueKing:
			score += 10
		case c.Value == deck.ValueAce:
			aces++
			score++
			score += 10 // TODO: Remove THIS
			// fallthrough // continue to next and don't break
		}
	}

	// TODO: Decide how to calculate aces (0 or 10)

	return score
}

type player interface {
	Name() string
	Hand() []deck.Card
	Deal(deck.Card)
	Play() (PlayType, error)
}

type basePlayer struct {
	hand []deck.Card
}

// implemented Hand()
func (p *basePlayer) Hand() []deck.Card {
	return p.hand
}

// implemented Deal()
func (p *basePlayer) Deal(c deck.Card) {
	p.hand = append(p.hand, c)
}

type humanPlayer struct {
	basePlayer
}

// implemented Play()
func (p *humanPlayer) Play() (PlayType, error) {
	for {
		fmt.Println("Hit or Stand h/s?")

		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			return 0, err
		}

		input = strings.ToLower(strings.TrimSpace(input))

		if input == "" {
			continue
		}

		switch input[0:1] {
		case "h":
			return PlayHit, nil
		case "s":
			return PlayStand, nil
		default:
			fmt.Println("Invalid input")
			continue // continue to next and don't break
		}
	}
}

// implemented Name()
func (p *humanPlayer) Name() string {
	return "Player"
}

type dealerPlayer struct {
	basePlayer
}

// implemented Play()
func (p *dealerPlayer) Play() (PlayType, error) {
	if calcScore(p.hand) <= 16 {
		return PlayHit, nil
	} else {
		return PlayStand, nil
	}
}

// implemented Name()
func (p *dealerPlayer) Name() string {
	return "Dealer"
}
