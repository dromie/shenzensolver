package solver

import (
	"fmt"
	"strings"
)

type CardSuit int

//go:generate stringer -type=CardSuit
const (
	DEFAULTSUIT CardSuit = iota
	RED
	GREEN
	BLACK
	FLOWER
)

var COLORDICT = map[string]CardSuit{
	"r": RED,
	"g": GREEN,
	"b": BLACK,
	"o": FLOWER,
}

var REVCOLORDICT = map[CardSuit]string{
	RED:    "r",
	GREEN:  "g",
	BLACK:  "b",
	FLOWER: "o",
}

type CardValue int

//go:generate stringer -type=CardValue
const (
	DEFAULTVALUE CardValue = iota
	ONE
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	BLOCK      = 254
	BLOCK_HOLD = 255
)

type Card struct {
	Suit  CardSuit
	Value CardValue
}

func (c Card) String() string {
	return fmt.Sprintf("%v of %v", c.Value, c.Suit)
}

func constructCard(card string) Card {
	if card == "" {
		return Card{Suit: DEFAULTSUIT, Value: BLOCK}
	}
	color := COLORDICT[strings.ToLower(string(card[0]))]
	var number CardValue
	if len(card) == 2 {
		number = CardValue(card[1] - '0')
	} else {
		number = BLOCK
	}
	return Card{Suit: color, Value: number}
}

func (c *Card) is_solution(other Card) bool {
	if c.Value == BLOCK {
		return false
	}
	return c.Value == 1 && other == Card{} || c.Suit == other.Suit && c.Value == other.Value+1
}

func (c *Card) can_be_put_over(other Card) bool {
	if other == (Card{}) {
		return true
	}
	if c.Value == BLOCK {
		return false
	}
	if c.Suit != other.Suit && c.Value == other.Value-1 {
		return true
	}
	return false
}

/*
func (c Card) can_be_put_over_column(col []Card) bool {
	if len(col) == 0 || col == nil {
		return true
	}
	return c.can_be_put_over(col[len(col)-1])
}
*/
