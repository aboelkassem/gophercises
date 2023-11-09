// Package deck can be used to create decks of playing cards
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

// Suit is used to define the suit of a Card
type Suit int // instead of enums

// The list of Suits that can be assigned to a Card
const (
	SuitSpades Suit = iota // iota = auto incremental index from 0, must be in const
	SuitHearts             // no need to define type and iota again, it will understand we are in the same const grouping
	SuitDiamonds
	SuitClubs
	SuitJoker
)

func (s Suit) String() string { // custom ToString implementation that Implemented fmt.Stringer
	switch s {
	case SuitSpades:
		return "♠"
	case SuitHearts:
		return "♥"
	case SuitDiamonds:
		return "♦"
	case SuitClubs:
		return "♣"
	case SuitJoker:
		return "J"
	default:
		return "Unknown"
	}
}

// Value is used to define the value of a Card
type Value int

// The list of Values that can be assigned to a Card
const (
	_ Value = iota
	_
	ValueTwo
	ValueThree
	ValueFour
	ValueFive
	ValueSix
	ValueSeven
	ValueEight
	ValueNine
	ValueTen
	ValueJack
	ValueQueen
	ValueKing
	ValueAce
)

func (v Value) String() string {
	switch v {
	case ValueJack:
		return "J"
	case ValueQueen:
		return "Q"
	case ValueKing:
		return "K"
	case ValueAce:
		return "A"
	default:
		return strconv.Itoa(int(v))
	}
}

// Card holds a combination of a Suit and a Value
type Card struct {
	Suit  Suit
	Value Value
}

func (c Card) String() string {
	if c.Suit == SuitJoker {
		return "[Joker]"
	}
	return fmt.Sprintf("[ %v  %-2v]", c.Suit, c.Value)
}

// functional programming,
// meaning just define your logic with functions caring only with input and output, don't care of implementation details or let it for user
// must follow func signature (input, output)

type Option func([]Card) []Card

// ... is like params in c# and its optional
// takes from 1 to infinity number of parameters

// New creates a new deck with the specified Options
func New(opts ...Option) []Card {
	var deck []Card

	// loop for suites
	for suit := SuitSpades; suit <= SuitClubs; suit++ {
		for value := ValueTwo; value <= ValueAce; value++ {
			deck = append(deck, Card{
				Suit:  suit,
				Value: value,
			})
		}
	}

	// loop for options and apply them
	// functional programming
	for _, opt := range opts {
		deck = opt(deck)
	}

	return deck
}

// you can use functions in go like variables
// must follow
// type Option func([]Card) []Card
// func OptionShuffle(cards []Card) []Card {
// 	rand.Shuffle(len(cards), func(i, j int) {
// 		// cards[i] = cards[j]
// 		// cards[j] = cards[i]
// 		cards[i], cards[j] = cards[j], cards[i]
// 	})
// 	return cards
// }

// OptionShuffle shuffles a deck
func OptionShuffle() Option {
	return func(cards []Card) []Card {
		rand.Shuffle(len(cards), func(i, j int) {
			// cards[i] = cards[j]
			// cards[j] = cards[i]
			cards[i], cards[j] = cards[j], cards[i]
		})
		return cards
	}
}

// functional programming
// will be callable like
//
//	deck.OptionSort(func(i, j deck.Card) bool {
//		return i.Suit > j.Suit
//	}),
// can return
// Option(func(cards []Card) []Card {
// })
// or direct
// func(cards []Card) []Card {
// }

// OptionSort can sort a deck based on the sorting function fn
func OptionSort(fn func(i, j Card) bool) Option {
	return func(cards []Card) []Card {
		sort.Slice(cards, func(i, j int) bool {
			return fn(cards[i], cards[j])
		})
		return cards
	}
}

// SortDefault provides the default sorting logic for a deck
func SortDefault(i, j Card) bool {
	return i.Suit > j.Suit ||
		i.Suit == j.Suit && i.Value < j.Value
}

// OptionAddJokers adds n arbitary Jokers to the end of a deck
func OptionAddJokers(n int) Option {
	return Option(func(cards []Card) []Card {
		for i := 1; i <= n; i++ {
			cards = append(cards, Card{
				Suit: SuitJoker,
			})
		}
		return cards
	})
}

// OptionExclude uses fn to know which cards to exludes from a deck
func OptionExclude(fn func(Card) bool) Option {
	return Option(func(cards []Card) []Card {
		var newCards []Card
		for _, c := range cards {
			if fn(c) {
				continue
			}
			newCards = append(newCards, c)
		}
		return newCards
	})
}

// OptionCompose composes a bigger deck by adding other decks to a deck
func OptionCompose(decks ...[]Card) Option {
	return func(cards []Card) []Card {
		for _, deck := range decks {
			cards = append(cards, deck...) // deck... = will take this list one by one
		}
		return cards
	}
}
