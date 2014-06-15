package holdem

import (
	"fmt"
	"math/rand"
	"time"
)

type Suit int
type Value int
type RoundStatus int
type PlayerStatus int

const (
	Undefined Suit = iota
	Spades
	Diamonds
	Hearts
	Clubs

	Flop RoundStatus = iota
	Turn
	River

	Folded PlayerStatus = iota // No longer in the round
	Active                     // Participating in the round (checking, raised, all in)
)

var playerDB map[string]Player

type Callback func(game *Game, done chan bool)

type Game struct {
	deck       []Card
	community  []Card
	pot        uint32
	currentBet uint32

	currentBetter *Player
	// currentBetCompleted chan bool
	players []*Player // Around the table in this round
	frozen  bool

	// blind uint32 // Might not require blinds in IRC gameplay

	preRoundCallback  func(*Game, chan bool)
	communityCallback func(RoundStatus, []Card)
	/*
		postFlopCallback  Callback
		postTurnCallback  Callback
		postRiverCallback Callback
	*/

	displayPlayerCardCallback func(string, []Card, chan bool)
	betCallback               func(*Game, string)
}

type Player struct {
	// TODO: ID field for db?
	// TODO: Replace player name with a generic pointer with user data instead
	Name    string
	Bet     uint32
	Status  PlayerStatus
	Balance uint32

	Hand []Card
	// AllIn bool
}

type Card struct {
	Suit
	Value
}

func (s Suit) String() string {
	switch s {
	case Spades:
		return "\u2660"
	case Hearts:
		return "\u2665"
	case Diamonds:
		return "\u2666"
	case Clubs:
		return "\u2663"
	default:
		println("INVALID SUIT IN String()")
		return ""
	}
}

func (v Value) String() string {
	if v <= 10 {
		return fmt.Sprintf("%d", v)
	}

	return string(v)
}

func (c Card) String() string {
	return fmt.Sprintf("|%s%s|", c.Suit, c.Value)
}

func New() Game {
	g := Game{}

	rand.Seed(time.Now().UnixNano())

	g.createNewDeck()
	g.players = make([]*Player, 0, 2)

	return g
}

const (
	DECKS     = 6
	DECK_SIZE = 52
)

func (g *Game) SetPreRoundCallback(c func(*Game, chan bool)) {
	g.preRoundCallback = c
}

func (g *Game) SetDisplayPlayerCardCallback(c func(string, []Card, chan bool)) {
	g.displayPlayerCardCallback = c
}

func (g *Game) SetBetCallback(c func(*Game, string)) {
	g.betCallback = c
}

func (g *Game) SetCommunityCallback(c func(RoundStatus, []Card)) {
	g.communityCallback = c
}

func (g *Game) AddPlayer(name string) {
	/*
		if g.frozen {
			return
		}
	*/

	if _, exists := playerDB[name]; exists {
		// TODO: error value?
		// err = errors.New("Player already exists")
		return
	}

	player := newPlayer(name)

	if playerDB == nil {
		playerDB = make(map[string]Player)
	}

	playerDB[name] = player

	g.players = append(g.players, &player)
	// println("Appended", name, "len(g.players) ==", len(g.players))
}

func (g *Game) JoinTable(name string) {
}

func (g *Game) LeaveTable(name string) {
}

func newPlayer(name string) Player {
	return Player{name, 0, 0, 100, nil} // TODO: Configurable initial balance
}

func (g *Game) Play() {
	g.newRound()    // Initiate the round
	g.dealPreFlop() // 2 cards to each player

	g.doBets()

	g.dealFlop() // Deal 3 community cards
	g.doBets()
	g.dealTurn() // 4th community card
	g.doBets()
	g.dealRiver() // 5th community card
	g.doBets()

	/*
		//g.Showdown()

		g.finishRound()
	*/
}

func (g *Game) shuffleDeck() {
	newDeck := make([]Card, DECKS*DECK_SIZE)
	p := rand.Perm(DECKS * DECK_SIZE)

	for i, k := range p {
		newDeck[i] = g.deck[k]
	}

	g.deck = newDeck
}

func (g *Game) shufflePlayers() {
	if g.players != nil {
		g.players = append(g.players[1:], g.players[0])
	}

	/*
		newPlayers := make([]Player, len(g.players))
		p := rand.Perm(len(g.players))

		for i, k := range p {
			newPlayers[i] = g.players[k]
		}

		g.players = newPlayers
	*/

}

func (g *Game) createNewDeck() {
	g.deck = make([]Card, DECKS*DECK_SIZE, DECKS*DECK_SIZE)

	suits := []Suit{Spades, Diamonds, Hearts, Clubs}
	values := []Value{'A', 'K', 'Q', 'J', 10, 9, 8, 7, 6, 5, 4, 3, 2}

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

		g.frozen = false
		g.currentBet = 0

		g.shuffleDeck()

		go g.preRoundCallback(g, done) // Players register for a new round (.hit)
		// g.shufflePlayers()

		<-done
		// g.frozen = true
	}
}

func (g *Game) dealCard() (c Card) {
	c = g.deck[0]
	g.deck = g.deck[1:]
	return
}

func (g *Game) currentBetterDone() {
	// g.currentBetCompleted <- true
}

func (g *Game) dealCards(n int) (c []Card) {
	c = g.deck[:n]
	g.deck = g.deck[n:]
	return
}

func (g *Game) isEndOfBets() bool {
	for _, p := range g.players {
		if p.Status == Folded {
			continue
		}

		if p.Bet != g.currentBet {
			println(p.Name, "had bet of", p.Bet, "current is", g.currentBet)
			return false
		}
	}

	return true
}

func (g *Game) doBet() {
	for _, player := range g.players {
		if player.Status == Folded {
			continue
		}

		g.currentBetter = player
		g.betCallback(g, player.Name)
	}
}

func (g *Game) doBets() {
	g.doBet()

	for !g.isEndOfBets() {
		g.doBet()
	}
}

func (g *Game) Raise(player string, bet uint32) {
	g.currentBet += bet
	println(player, "raised bet to", g.currentBet)
	g.currentBetter.Bet = g.currentBet
	// go g.currentBetterDone()
}

func (g *Game) Check(player string) {
	println(player, "checked at", g.currentBet)
	g.currentBetter.Bet = g.currentBet

	// go g.currentBetterDone()
}

func (g *Game) Fold(player string) {
	println(player, "folded")
	g.currentBetter.Status = Folded
	// go g.currentBetterDone()
}

func (g *Game) BetTimeout(player string) {
	if player == g.currentBetter.Name {
		g.currentBetter.Status = Folded
		g.currentBetterDone()
	}
}

func (g *Game) dealPreFlop() {
	for i, p := range g.players {
		// p.Hand = append(p.Hand, g.dealCards(2)...)
		g.players[i].Hand = append(p.Hand, g.dealCards(2)...)
	}

	c := make(chan bool)
	for _, p := range g.players {
		//println(p.Name, p.Hand)
		go g.displayPlayerCardCallback(p.Name, p.Hand, c)
		<-c
	}
	// g.PostFlopCallback()
}

func (g *Game) dealFlop() {
	g.community = append(g.community, g.dealCards(3)...)

	// TODO: g.PostFlopCallback()

	g.communityCallback(Flop, g.community)
}

func (g *Game) dealTurn() {
	g.community = append(g.community, g.dealCard())

	// TODO: g.PostTurnCallback()
	g.communityCallback(Turn, g.community)
}

func (g *Game) dealRiver() {
	g.community = append(g.community, g.dealCard())

	// TODO: g.PostRiverCallback()
	g.communityCallback(River, g.community)
}

func (g *Game) finishRound() {
	// TODO: g.updateBalanceCallback(player, amount)
	// TODO: g.endOfRoundCallback()
}
