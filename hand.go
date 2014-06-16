package holdem

import (
	"bytes"
	"fmt"
	"strings"
)

type HandClass uint
type HandValue uint32

const (
	Clubs int = iota
	Diamonds
	Hearts
	Spades

	spadeOffset   = 13 * uint32(Spades)
	clubOffset    = 13 * uint32(Clubs)
	diamondOffset = 13 * uint32(Diamonds)
	heartOffset   = 13 * uint32(Hearts)
)

const (
	HighCard HandClass = iota
	Pair
	TwoPair
	Trips
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
)

const (
	handTypeShift   = 24
	topCardShift    = 16
	secondCardShift = 12
	thirdCardShift  = 8
	fourthCardShift = 4
	fifthCardShift  = 0
)

const (
	straightFlushVal = uint32(StraightFlush) << handTypeShift
	straightVal      = uint32(Straight) << handTypeShift
	flushVal         = uint32(Flush) << handTypeShift
	fullHouseVal     = uint32(FullHouse) << handTypeShift
	fourOfAKindVal   = uint32(FourOfAKind) << handTypeShift
	tripsVal         = uint32(Trips) << handTypeShift
	twoPairVal       = uint32(TwoPair) << handTypeShift
	pairVal          = uint32(Pair) << handTypeShift
	highCardVal      = uint32(HighCard) << handTypeShift
)

const (
	topCardMask    uint32 = 0x000F0000
	secondCardMask uint32 = 0x0000F000
	fifthCardMask  uint32 = 0x0000000F
	cardMask       uint64 = 0x0F
)

const (
	numberOfCards = 52
	cardWidth     = 4
)

// Hand is a player's combination of cards that form a "poker hand".
type Hand uint64

// NewHand creates a hand from the player and the table cards.
func NewHand(player, table []Card) Hand {
	var hand Hand

	for _, card := range player {
		hand |= (1 << card)
	}
	for _, card := range table {
		hand |= (1 << card)
	}

	return hand
}

// NewHandCards creates a hand from the cards given.
func NewHandCards(cards ...[]Card) Hand {
	var hand Hand

	for _, cardArr := range cards {
		for _, card := range cardArr {
			hand |= (1 << card)
		}
	}

	return hand
}

// NewHandStr creates a hand from the cards given.
func NewHandStr(h string) Hand {
	var hand Hand

	cards := strings.Split(h, " ")
	for _, cardstr := range cards {
		hand |= (1 << NewCardStr(cardstr))
	}

	return hand
}

// Display computes the value of the hand and creates an output
// string to display to a user with as well.
func (h Hand) Display() (HandValue, string) {
	nCards := countBits(uint64(h))

	sc := uint32(h>>clubOffset) & 0x1FFF
	sd := uint32(h>>diamondOffset) & 0x1FFF
	sh := uint32(h>>heartOffset) & 0x1FFF
	ss := uint32(h>>spadeOffset) & 0x1FFF

	val := h.ValueCards(nCards)
	var str string
	cls := val.Class()

	switch cls {
	case HighCard, Pair, TwoPair, Trips, Straight, FullHouse, FourOfAKind:
		str = val.String()
	case Flush, StraightFlush:
		if cls == Flush {
			str = "Flush"
		} else {
			str = "Straight Flush"
		}
		var suitRune rune
		switch {
		case nBitsTable[ss] >= 5:
			suitRune = '\u2660'
		case nBitsTable[sc] >= 5:
			suitRune = '\u2663'
		case nBitsTable[sd] >= 5:
			suitRune = '\u2666'
		case nBitsTable[sh] >= 5:
			suitRune = '\u2665'
		}
		str = fmt.Sprintf("%s (%c) with %v high", str, suitRune, val.TopCard())
	}

	return val, str
}

// String changes a HandValue into a readable representation.
func (h HandValue) String() string {
	b := &bytes.Buffer{}

	cls := h.Class()
	switch cls {
	case HighCard:
		fmt.Fprintf(b, "High card: %-v", h.TopCard())
	case Pair:
		fmt.Fprintf(b, "One pair: %-v", h.TopCard())
	case TwoPair:
		fmt.Fprintf(b, "Two pair: %-v's and %-v's with %-v kicker", h.TopCard(), h.SecondCard(), h.ThirdCard())
	case Trips:
		fmt.Fprintf(b, "Three of a kind: %-v's", h.TopCard())
	case Straight:
		fmt.Fprintf(b, "Straight with %-v high", h.TopCard())
	case Flush:
		b.WriteString("Flush")
	case FullHouse:
		fmt.Fprintf(b, "Full house: %-v's and %-v's", h.TopCard(), h.SecondCard())
	case FourOfAKind:
		fmt.Fprintf(b, "Four of a kind: %-v", h.TopCard(), h.SecondCard())
	case StraightFlush:
		b.WriteString("Straight Flush")
	}

	return b.String()
}

// Class determines what class of hand this is.
func (h HandValue) Class() HandClass {
	return HandClass(h >> handTypeShift)
}

// TopCard determines the top card in the hand.
func (h HandValue) TopCard() Card {
	return Card((uint64(h) >> topCardShift) & cardMask)
}

// SecondCard determines the second card in the hand.
func (h HandValue) SecondCard() Card {
	return Card((uint64(h) >> secondCardShift) & cardMask)
}

// ThirdCard determines the third card in the hand.
func (h HandValue) ThirdCard() Card {
	return Card((uint64(h) >> thirdCardShift) & cardMask)
}

// ValueCards computes the value of the hand. Determines automatically
// how many cards exist.
func (h Hand) Value() HandValue {
	return h.ValueCards(countBits(uint64(h)))
}

// ValueCards computes the value of the hand of size nCards.
func (h Hand) ValueCards(nCards int) HandValue {
	var val, fourMask, threeMask, twoMask uint32

	sc := uint32(h>>clubOffset) & 0x1FFF
	sd := uint32(h>>diamondOffset) & 0x1FFF
	sh := uint32(h>>heartOffset) & 0x1FFF
	ss := uint32(h>>spadeOffset) & 0x1FFF

	values := sc | sd | sh | ss
	nValues := nBitsTable[values]
	nDups := uint32(nCards) - uint32(nValues)

	if nValues >= 5 {
		switch {
		case nBitsTable[ss] >= 5:
			if v := straightTable[ss]; v != 0 {
				val = straightFlushVal + (uint32(v) << topCardShift)
			} else {
				val = flushVal + topFiveCardsTable[ss]
			}
		case nBitsTable[sc] >= 5:
			if v := straightTable[sc]; v != 0 {
				val = straightFlushVal + (uint32(v) << topCardShift)
			} else {
				val = flushVal + topFiveCardsTable[sc]
			}
		case nBitsTable[sd] >= 5:
			if v := straightTable[sd]; v != 0 {
				val = straightFlushVal + (uint32(v) << topCardShift)
			} else {
				val = flushVal + topFiveCardsTable[sd]
			}
		case nBitsTable[sh] >= 5:
			if v := straightTable[sh]; v != 0 {
				val = straightFlushVal + (uint32(v) << topCardShift)
			} else {
				val = flushVal + topFiveCardsTable[sh]
			}
		default:
			st := straightTable[values]
			if st != 0 {
				val = straightVal + uint32(st<<topCardShift)
			}
		}

		if val != 0 && nDups < 3 {
			return HandValue(val)
		}
	}

	switch nDups {
	case 0:
		return HandValue(highCardVal + topFiveCardsTable[values])
	case 1:
		var t, kickers uint32

		twoMask = values ^ (sc ^ sd ^ sh ^ ss)

		val = pairVal + uint32(topCardTable[twoMask]<<topCardShift)
		t = values ^ twoMask
		kickers = (topFiveCardsTable[t] >> cardWidth) & ^fifthCardMask
		val += kickers
		return HandValue(val)
	case 2:
		twoMask = uint32(values) ^ uint32(sc^sd^sh^ss)
		if twoMask != 0 {
			var t uint32 = values ^ twoMask
			return HandValue(twoPairVal +
				(topFiveCardsTable[twoMask] & (topCardMask | secondCardMask)) +
				(uint32(topCardTable[t]) << thirdCardShift))
		} else {
			var t, second uint32
			threeMask = ((sc & sd) | (sh & ss)) & ((sc & sh) | (sd & ss))
			val = tripsVal + (uint32(topCardTable[threeMask]) << topCardShift)

			t = values ^ threeMask
			second = uint32(topCardTable[t])
			val += (second << secondCardShift)
			t ^= (1 << second)
			val += uint32(topCardTable[t] << thirdCardShift)
			return HandValue(val)
		}
	default:
		fourMask = sh & sd & sc & ss
		if fourMask != 0 {
			tc := uint32(topCardTable[fourMask])
			val = fourOfAKindVal +
				(tc << topCardShift) +
				(uint32(topCardTable[values^(1<<tc)]) << secondCardShift)
			return HandValue(val)
		}

		twoMask = values ^ (sc ^ sd ^ sh ^ ss)
		if uint32(nBitsTable[twoMask]) != nDups {
			var tc, t uint32
			threeMask = ((sc & sd) | (sh & ss)) & ((sc & sh) | (sd & ss))
			val = fullHouseVal
			tc = uint32(topCardTable[threeMask])
			val += (tc << topCardShift)
			t = (twoMask | threeMask) ^ (1 << tc)
			val += (uint32(topCardTable[t]) << secondCardShift)
			return HandValue(val)
		}

		if val != 0 {
			return HandValue(val)
		} else {
			var top, second uint32

			val = twoPairVal
			top = uint32(topCardTable[twoMask])
			val += (top << topCardShift)
			second = uint32(topCardTable[twoMask^(1<<top)])
			val += (second << secondCardShift)
			val += uint32(topCardTable[values^(1<<top)^(1<<second)]) << thirdCardShift
			return HandValue(val)
		}
	}

	return HandValue(val)
}

func countBits(bits uint64) int {
	return int(bitCounts[int(bits&0x00000000000000FF)] +
		bitCounts[int((bits&0x000000000000FF00)>>8)] +
		bitCounts[int((bits&0x0000000000FF0000)>>16)] +
		bitCounts[int((bits&0x00000000FF000000)>>24)] +
		bitCounts[int((bits&0x000000FF00000000)>>32)] +
		bitCounts[int((bits&0x0000FF0000000000)>>40)] +
		bitCounts[int((bits&0x00FF000000000000)>>48)] +
		bitCounts[int((bits&0xFF00000000000000)>>56)])
}
