package holdem

import "testing"

func TestHand_Value(t *testing.T) {
	t.Parallel()
}

func TestHand_New(t *testing.T) {
	t.Parallel()

	player := []Card{
		NewCard(2, Clubs),
		NewCard(3, Diamonds),
	}
	table := []Card{
		NewCard(13, Hearts),
		NewCard(5, Spades),
		NewCard(8, Diamonds),
		NewCard(11, Spades),
		NewCard(12, Diamonds),
	}

	hand1 := NewHand(player, table)
	if hand1 == 0 {
		t.Error("Expected hand to be set.")
	}

	hand2 := NewHandCards(player, table)
	if hand2 == 0 {
		t.Error("Expected hand to be set.")
	}

	hand3 := NewHandStr("2c 3d kh 5s 8d js qd")
	if hand3 == 0 {
		t.Error("Expected hand to be set.")
	}

	if hand1 != hand2 || hand2 != hand3 {
		t.Errorf("The hands did not match: %d %d %d", hand1, hand2, hand3)
	}
}

func TestHand_Display(t *testing.T) {
	t.Parallel()

	player := []Card{
		NewCard(14, Diamonds),
		NewCard(13, Diamonds),
	}
	table := []Card{
		NewCard(2, Diamonds),
		NewCard(13, Hearts),
		NewCard(12, Diamonds),
		NewCard(3, Hearts),
		NewCard(12, Clubs),
	}

	hand := NewHand(player, table)
	exp := "Two pair: K♣'s and Q♣'s with A♣ kicker"
	if val, str := hand.Display(); val == 0 || str != exp {
		t.Errorf(`Expected: "%s", got: "%s"`, exp, str)
	}

	hand = NewHandStr("ad ac as 3d 4h 10h 10d")
	exp = "Full house: A♣'s and 10♣'s"
	if val, str := hand.Display(); val == 0 || str != exp {
		t.Errorf(`Expected: "%s", got: "%s"`, exp, str)
	}
}

func TestHand_Values_5(t *testing.T) {
	t.Parallel()

	t.SkipNow()

	hcls := []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
	nCards := 52
	count := 0

	var i1, i2, i3, i4, i5 int
	var card1, n2, n3, n4 uint64
	var hnd Hand

	for i1 = nCards - 1; i1 >= 0; i1-- {
		card1 = cardMasksTable[i1]

		for i2 = i1 - 1; i2 >= 0; i2-- {
			n2 = card1 | cardMasksTable[i2]

			for i3 = i2 - 1; i3 >= 0; i3-- {
				n3 = n2 | cardMasksTable[i3]

				for i4 = i3 - 1; i4 >= 0; i4-- {
					n4 = n3 | cardMasksTable[i4]
					for i5 = i4 - 1; i5 >= 0; i5-- {
						hnd = Hand(n4 | cardMasksTable[i5])
						hcls[hnd.Value().Class()]++
						count++
					}
				}
			}
		}
	}

	if exp, got := 1302540, hcls[HighCard]; exp != got {
		t.Errorf("Expected: %v, got: %v", exp, got)
	}
	if exp, got := 1098240, hcls[Pair]; exp != got {
		t.Errorf("Expected: %v, got: %v", exp, got)
	}
	if exp, got := 123552, hcls[TwoPair]; exp != got {
		t.Errorf("Expected: %v, got: %v", exp, got)
	}
	if exp, got := 54912, hcls[Trips]; exp != got {
		t.Errorf("Expected: %v, got: %v", exp, got)
	}
	if exp, got := 10200, hcls[Straight]; exp != got {
		t.Errorf("Expected: %v, got: %v", exp, got)
	}
	if exp, got := 5108, hcls[Flush]; exp != got {
		t.Errorf("Expected: %v, got: %v", exp, got)
	}
	if exp, got := 3744, hcls[FullHouse]; exp != got {
		t.Errorf("Expected: %v, got: %v", exp, got)
	}
	if exp, got := 624, hcls[FourOfAKind]; exp != got {
		t.Errorf("Expected: %v, got: %v", exp, got)
	}
	if exp, got := 40, hcls[StraightFlush]; exp != got {
		t.Errorf("Expected: %v, got: %v", exp, got)
	}

	if exp, got := 2598960, count; exp != got {
		t.Errorf("Expected: %v, got: %v", exp, got)
	}
}

func TestHand_Values_7(t *testing.T) {
	t.Parallel()
	t.SkipNow()

	hcls := []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
	nCards := 52
	count := 0

	var i1, i2, i3, i4, i5, i6, i7 int
	var card1, n2, n3, n4, n5, n6 uint64
	var hnd Hand

	for i1 = nCards - 1; i1 >= 0; i1-- {
		card1 = cardMasksTable[i1]

		for i2 = i1 - 1; i2 >= 0; i2-- {
			n2 = card1 | cardMasksTable[i2]

			for i3 = i2 - 1; i3 >= 0; i3-- {
				n3 = n2 | cardMasksTable[i3]

				for i4 = i3 - 1; i4 >= 0; i4-- {
					n4 = n3 | cardMasksTable[i4]

					for i5 = i4 - 1; i5 >= 0; i5-- {
						n5 = n4 | cardMasksTable[i5]

						for i6 = i5 - 1; i6 >= 0; i6-- {
							n6 = n5 | cardMasksTable[i6]

							for i7 = i6 - 1; i7 >= 0; i7-- {
								hnd = Hand(n6 | cardMasksTable[i7])
								hcls[hnd.Value().Class()]++
								count++
							}
						}
					}
				}
			}
		}
	}

	if exp, got := 58627800, hcls[Pair]; exp != got {
		t.Errorf("Expected: %d, got: %d", exp, got)
	}
	if exp, got := 31433400, hcls[TwoPair]; exp != got {
		t.Errorf("Expected: %d, got: %d", exp, got)
	}
	if exp, got := 6461620, hcls[Trips]; exp != got {
		t.Errorf("Expected: %d, got: %d", exp, got)
	}
	if exp, got := 6180020, hcls[Straight]; exp != got {
		t.Errorf("Expected: %d, got: %d", exp, got)
	}
	if exp, got := 4047644, hcls[Flush]; exp != got {
		t.Errorf("Expected: %d, got: %d", exp, got)
	}
	if exp, got := 3473184, hcls[FullHouse]; exp != got {
		t.Errorf("Expected: %d, got: %d", exp, got)
	}
	if exp, got := 224848, hcls[FourOfAKind]; exp != got {
		t.Errorf("Expected: %d, got: %d", exp, got)
	}
	if exp, got := 41584, hcls[StraightFlush]; exp != got {
		t.Errorf("Expected: %d, got: %d", exp, got)
	}

	if exp, got := 133784560, count; exp != got {
		t.Errorf("Expected: %d, got: %d", exp, got)
	}
}

func BenchmarkHand(b *testing.B) {
	hand := NewHandStr("ad as 3d 5d 7h 10d 10c")
	for i := 0; i < b.N; i++ {
		hand.ValueCards(7)
	}
}
