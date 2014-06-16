package holdem

import (
	"fmt"
	"io"
	"strings"
)

// Card represents a single card.
type Card uint8

// NewCard takes a value and a suit to create a card. The value is modified
// to benefit calculation, and so calling .Value() will not return this value.
func NewCard(value, suit int) Card {
	return Card((value - 2) + (suit * 13))
}

// NewCardStr makes a card from a string.
func NewCardStr(str string) Card {
	str = strings.ToLower(str)
	var value, suit int

	index := 0
	switch str[index] {
	case '2', '3', '4', '5', '6', '7', '8', '9':
		value = int(str[index] - '2')
	case '1':
		value = 8
		if str[index+1] != '0' {
			panic("Malformed 10")
		}
		index++
	case 'j':
		value = 9
	case 'q':
		value = 10
	case 'k':
		value = 11
	case 'a':
		value = 12
	}
	index++

	switch rune(str[index]) {
	case 'c', '\u2663':
		suit = Clubs
	case 'd', '\u2666':
		suit = Diamonds
	case 'h', '\u2665':
		suit = Hearts
	case 's', '\u2660':
		suit = Spades
	}

	return Card(value + (suit * 13))
}

// Format implements Formatter.
func (c Card) Format(f fmt.State, kind rune) {
	value := c.Value()
	var valueStr string

	if f.Flag('-') {
		switch value {
		case 8:
			valueStr = "10"
		case 9:
			valueStr = "J"
		case 10:
			valueStr = "Q"
		case 11:
			valueStr = "K"
		case 12:
			valueStr = "A"
		default:
			valueStr = string(value + '2')
		}

		io.WriteString(f, valueStr)
	} else {
		io.WriteString(f, c.String())
	}
}

// String implements Stringer.
func (c Card) String() string {
	value := c.Value()
	suit := c.Suit()
	var valueStr string
	var suitRune rune

	switch value {
	case 8:
		valueStr = "10"
	case 9:
		valueStr = "J"
	case 10:
		valueStr = "Q"
	case 11:
		valueStr = "K"
	case 12:
		valueStr = "A"
	default:
		valueStr = string(value + '2')
	}

	switch suit {
	case Clubs:
		suitRune = '\u2663'
	case Diamonds:
		suitRune = '\u2666'
	case Hearts:
		suitRune = '\u2665'
	case Spades:
		suitRune = '\u2660'
	}
	return fmt.Sprintf("%s%c", valueStr, suitRune)
}

// Value returns the transformed value, not the value that was put into the card
func (c Card) Value() int {
	return int(c) % 13
}

// Suit returns the suit of the card.
func (c Card) Suit() int {
	return int(c) / 13
}
