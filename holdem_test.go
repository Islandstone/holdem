package holdem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func isValidCard(c Card) bool {
	return c >= 0 && c < (52*4)
}

func isValidDeck(deck []Card, t *testing.T) bool {
	count := make(map[Card]int)

	for _, v := range deck {
		count[v] += 1
	}

	for k, v := range count {
		if v != Decks {
			t.Error("Invalid card count:", k)
			return false
		}
	}

	return true
}

func TestInit(t *testing.T) {
	game := New()

	assert.NotEmpty(t, game.deck)
	assert.Equal(t, Decks*DeckSize, len(game.deck))

	valid := true
	for _, v := range game.deck {
		if !isValidCard(v) {

			assert.Fail(t, "Invalid card")
			t.Log(v)
			valid = false
			break
		}
	}

	assert.True(t, valid)
	assert.True(t, isValidDeck(game.deck, t))
}

func TestShuffle(t *testing.T) {
	game := New()

	game.shuffleDeck()
	assert.True(t, isValidDeck(game.deck, t))

	old := make([]Card, Decks*DeckSize)
	copy(old, game.deck)

	game.shuffleDeck()
	assert.True(t, isValidDeck(game.deck, t))

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
	t.SkipNow()
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
	t.SkipNow()
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
