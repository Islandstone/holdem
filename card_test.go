package holdem

import (
	"fmt"
	"testing"
)

func TestCard_New(t *testing.T) {
	t.Parallel()

	if c := NewCard(3, Spades); c != 40 {
		t.Error("Expected the string to be parsed to a proper card value.")
	}

	if c := NewCardStr("3s"); c != 40 {
		t.Error("Expected the string to be parsed to a proper card value.")
	}
}

func TestCard_String(t *testing.T) {
	t.Parallel()

	cards := []struct {
		Value  int
		Suit   int
		Expect string
	}{
		{2, Clubs, "2♣"},
		{3, Clubs, "3♣"},
		{4, Clubs, "4♣"},
		{5, Clubs, "5♣"},
		{6, Clubs, "6♣"},
		{7, Clubs, "7♣"},
		{8, Clubs, "8♣"},
		{9, Clubs, "9♣"},
		{10, Clubs, "10♣"},
		{11, Clubs, "J♣"},
		{12, Clubs, "Q♣"},
		{13, Clubs, "K♣"},
		{14, Clubs, "A♣"},
		{2, Diamonds, "2♦"},
		{3, Diamonds, "3♦"},
		{4, Diamonds, "4♦"},
		{5, Diamonds, "5♦"},
		{6, Diamonds, "6♦"},
		{7, Diamonds, "7♦"},
		{8, Diamonds, "8♦"},
		{9, Diamonds, "9♦"},
		{10, Diamonds, "10♦"},
		{11, Diamonds, "J♦"},
		{12, Diamonds, "Q♦"},
		{13, Diamonds, "K♦"},
		{14, Diamonds, "A♦"},
		{2, Hearts, "2♥"},
		{3, Hearts, "3♥"},
		{4, Hearts, "4♥"},
		{5, Hearts, "5♥"},
		{6, Hearts, "6♥"},
		{7, Hearts, "7♥"},
		{8, Hearts, "8♥"},
		{9, Hearts, "9♥"},
		{10, Hearts, "10♥"},
		{11, Hearts, "J♥"},
		{12, Hearts, "Q♥"},
		{13, Hearts, "K♥"},
		{14, Hearts, "A♥"},
		{2, Spades, "2♠"},
		{3, Spades, "3♠"},
		{4, Spades, "4♠"},
		{5, Spades, "5♠"},
		{6, Spades, "6♠"},
		{7, Spades, "7♠"},
		{8, Spades, "8♠"},
		{9, Spades, "9♠"},
		{10, Spades, "10♠"},
		{11, Spades, "J♠"},
		{12, Spades, "Q♠"},
		{13, Spades, "K♠"},
		{14, Spades, "A♠"},
	}

	for _, c := range cards {
		if got := NewCard(c.Value, c.Suit).String(); got != c.Expect {
			t.Errorf("Expected: %s, got: %s", c.Expect, got)
		}
	}
}

func TestCard_Format(t *testing.T) {
	t.Parallel()

	c := NewCard(10, Spades)
	if exp, got := "10♠ 10", fmt.Sprintf("%v %-v", c, c); exp != got {
		t.Errorf("Expected: %s, got: %s", exp, got)
	}
}

func TestCard_Value(t *testing.T) {
	t.Parallel()
	if 8 != NewCard(10, Spades).Value() {
		t.Error("Expected it to have a value of 8.")
	}
}

func TestCard_Suit(t *testing.T) {
	t.Parallel()
	if Spades != NewCard(2, Spades).Suit() {
		t.Error("Expected it to be the spades suit.")
	}
}
