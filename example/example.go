package main

import (
	"fmt"
	"os"
	h "holdem"
	"bufio"
)

var game h.Game
var stdin *bufio.Reader

func main() {
	game := h.New()

	game.SetPreRoundCallback(preRoundCallback)
	game.SetDisplayPlayerCardCallback(displayPlayerCardsCallback)
	game.SetBetCallback(betCallback)
	game.SetCommunityCallback(communityCallback)

	stdin = bufio.NewReader(os.Stdin)

	game.Play()

	fmt.Println("Done")
}

func preRoundCallback(g *h.Game, done chan bool) {
	fmt.Println("PreRound Callback")

	g.AddPlayer("A")
	g.AddPlayer("B")
	g.AddPlayer("C")

	done <- true
}

func displayPlayerCardsCallback(name string, cards []h.Card, done chan bool) {
	fmt.Printf("%s has cards %s\n", name, cards)

	done <- true
}

func betCallback(g *h.Game, name string) {
	out:
	for {
		fmt.Printf("%s, place your bet [r/c/f]: ", name)
		line, _, _ := stdin.ReadLine()

		switch line[0] {
		case 'r':
			g.Raise(name, 1)
			break out
		case 'c':
			g.Check(name)
			break out
		case 'f':
			g.Fold(name)
			break out
		}
	}
}

func communityCallback(state h.RoundStatus, cards []h.Card) {
	switch state {
	case h.Flop:
		fmt.Printf("Flop: %s\n", cards)
	case h.Turn:
		fmt.Printf("Turn: %s\n", cards)
	case h.River:
		fmt.Printf("River: %s\n", cards)
	}
}
