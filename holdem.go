package holdem

import (
	"math/rand"
)

const (
	Undefined = 0
	Spades    = 1
	Diamonds  = 2
	Hearts    = 3
	Clubs     = 4

	Folded = 1 // No longer in the round
	Active = 2 // Participating in the round (checking, raised, all in)
)

var players map[string]Player

type Callback func(done chan bool)

type Game struct {
	deck      []Card
	community []Card

	currentBetter *Player
	players       []Player // Around the table in this round
	frozen        bool

	// blind uint32 // Might not require blinds in IRC gameplay

	preRoundCallback  Callback
	postFlopCallback  Callback
	postTurnCallback  Callback
	postRiverCallback Callback
}

type Player struct {
	// TODO: ID field for db?
	// TODO: Replace player name with a generic pointer with user data instead
	Name    string
	Bet     uint32
	Status  uint32
	Balance uint32

	Hand []Card
	// AllIn bool
}

type Card struct {
	Suit  int
	Value int
}

func New() Game {
	g := Game{}

	g.createNewDeck()
	players = nil

	return g
}

const (
	DECKS     = 6
	DECK_SIZE = 52
)

func (g *Game) SetPreRoundCallback(c func(done chan bool)) {
	g.preRoundCallback = c
}

func (g *Game) AddPlayer(name string) {
	if g.frozen {
		return
	}

	if _, exists := players[name]; exists {
		// TODO: error value?
		// err = errors.New("Player already exists")
		return
	}

	p := newPlayer(name)

	if players == nil {
		players = make(map[string]Player)
	}

	players[name] = p
	g.players = append(g.players, p)
}

func newPlayer(name string) Player {
	return Player{name, 0, 0, 100, nil} // TODO: Configurable initial balance
}

func (g *Game) Play() {
	g.newRound()    // Initiate the round
	g.dealPreFlop() // 2 cards to each player

	//g.DoBets()

	g.dealFlop() // Deal 3 community cards
	//g.DoBets()
	g.dealTurn() // 4th community card
	//g.DoBets()
	g.dealRiver() // 5th community card
	//g.DoBets()

	//g.Showdown()

	g.finishRound()
}

func (g *Game) shuffle() {
	newDeck := make([]Card, DECKS*DECK_SIZE)
	p := rand.Perm(DECKS * DECK_SIZE)

	for i, k := range p {
		newDeck[i] = g.deck[k]
	}

	g.deck = newDeck
}

func (g *Game) createNewDeck() {
	g.deck = make([]Card, DECKS*DECK_SIZE, DECKS*DECK_SIZE)

	suits := []int{Spades, Diamonds, Hearts, Clubs}
	values := []int{'A', 'K', 'Q', 'J', 10, 9, 8, 7, 6, 5, 4, 3, 2}

	index := 0

	for deck_count := 0; deck_count < DECKS; deck_count += 1 {
		for _, suit := range suits {
			for _, value := range values {
				g.deck[index] = Card{Suit: suit, Value: value}
				index += 1
			}
		}
	}
}

func (g *Game) newRound() {
	if g.preRoundCallback != nil {
		done := make(chan bool)

		go func() {
			g.frozen = false
			g.preRoundCallback(done)
		}()

		g.shuffle()
		<-done
		g.frozen = true
	}
}

func (g *Game) dealCard() (c Card) {
	c = g.deck[0]
	g.deck = g.deck[1:]
	return
}

func (g *Game) dealCards(n int) (c []Card) {
	c = g.deck[:n]
	g.deck = g.deck[n:]
	return
}

func (g *Game) endOfBets() bool {
	// TODO
	// Bets should end if
	// - there's only one player left (player wins)
	// - all players have been asked at least once
	return true
}

func (g *Game) doBets() {
	if !g.endOfBets() {
		g.doBets()
	}
}

func (g *Game) PlaceBet(player string, bet uint32) {
}

func (g *Game) BetTimeout(player string) {
}

func (g *Game) dealPreFlop() {
	for _, p := range g.players {
		p.Hand = append(p.Hand, g.dealCards(2)...)
	}

	// g.PostFlopCallback()
}

func (g *Game) dealFlop() {
	g.community = append(g.community, g.dealCards(3)...)

	// TODO: g.PostFlopCallback()
}

func (g *Game) dealTurn() {
	g.community = append(g.community, g.dealCard())

	// TODO: g.PostTurnCallback()
}

func (g *Game) dealRiver() {
	g.community = append(g.community, g.dealCard())

	// TODO: g.PostRiverCallback()
}

func (g *Game) finishRound() {
	// TODO: g.updateBalanceCallback(player, amount)
	// TODO: g.endOfRoundCallback()
}
