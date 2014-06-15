package holdem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func isValidCard(c Card) bool {
	if c.Suit != Clubs &&
		c.Suit != Diamonds &&
		c.Suit != Hearts &&
		c.Suit != Spades {
		println("Invalid suit", c.Suit)
		return false
	}

	if c.Value > 0 && c.Value <= 10 {
		return true
	}

	switch c.Value {
	case 'A', 'K', 'Q', 'J':
		return true
	}

	println("Invalid value", c.Value)
	return false
}

func isValidDeck(deck []Card) bool {
	count := make(map[Card]int)

	for _, v := range deck {
		count[v] += 1
	}

	for k, v := range count {
		if v != DECKS {
			println("Invalid card count:", k.Suit, k.Value, v)
			return false
		}
	}

	return true
}

func TestInit(t *testing.T) {
	game := New()

	assert.NotEmpty(t, game.deck)
	assert.Equal(t, DECKS*DECK_SIZE, len(game.deck))

	valid := true
	for _, v := range game.deck {
		if !isValidCard(v) {

			assert.Fail(t, "Invalid card")
			println(v.Suit, v.Value)
			valid = false
			break
		}
	}

	assert.True(t, valid)
	assert.True(t, isValidDeck(game.deck))
}

func TestShuffle(t *testing.T) {
	game := New()

	game.shuffleDeck()
	assert.True(t, isValidDeck(game.deck))

	old := make([]Card, DECKS*DECK_SIZE)
	copy(old, game.deck)

	game.shuffleDeck()
	assert.True(t, isValidDeck(game.deck))

	equal := true
	for i, v := range old {
		if v != game.deck[i] {
			equal = false
			break
		}
	}

	assert.False(t, equal)
}

func TestDeal(t *testing.T) {

}

func TestPreRound(t *testing.T) {
	game := New()

	game.SetPreRoundCallback(func(g *Game, done chan bool) {
		g.AddPlayer("A")
		g.AddPlayer("B")

		// done <- true
	})

	game.Play()

	assert.Equal(t, 2, len(game.players))
}

func TestDeals(t *testing.T) {
	game := New()

	game.SetPreRoundCallback(func(g *Game, done chan bool) {
		g.AddPlayer("A")
		g.AddPlayer("B")

		done <- true
	})

	game.dealPreFlop()

	for _, p := range game.players {
		assert.True(t, len(p.Hand) == 2)
	}

	assert.True(t, len(game.community) == 0)
	game.dealFlop()
	assert.True(t, len(game.community) == 3)
	game.dealTurn()
	assert.True(t, len(game.community) == 4)
	game.dealRiver()
	assert.True(t, len(game.community) == 5)
}

func TestBettingPlayerCanBet(t *testing.T) {
	return // Test disabled

	game := New()

	game.SetPreRoundCallback(func(g *Game, done chan bool) {
		g.AddPlayer("A")
		g.AddPlayer("B")

		// done <- true
	})

	game.dealPreFlop()

	game.doBets()

	assert.True(t, game.currentBetter.Name == "A")

	game.Check("A")
}
