package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

type card rune

type hand struct {
	cards [5]card
	typ   int
	bid   int
}

func (h hand) String() string {
	cardsStr := ""
	for _, c := range h.cards {
		cardsStr += string(c)
	}
	typStr := ""
	if h.typ >= 0 && h.typ < len(typeNames) {
		typStr = typeNames[h.typ]
	}
	return fmt.Sprintf("%s:%d(%s)", cardsStr, h.bid, typStr)
}

const (
	HIGH_CARD       = iota
	ONE_PAIR        = iota
	TWO_PAIR        = iota
	THREE_OF_A_KIND = iota
	FULL_HOUSE      = iota
	FOUR_OF_A_KIND  = iota
	FIVE_OF_A_KIND  = iota
)

var typeNames = map[int]string{
	HIGH_CARD:       "High card",
	ONE_PAIR:        "One pair",
	TWO_PAIR:        "Two pair",
	THREE_OF_A_KIND: "Three of a Kind",
	FULL_HOUSE:      "Full house",
	FOUR_OF_A_KIND:  "Four of a kind",
	FIVE_OF_A_KIND:  "Five of a kind",
}

var kindOrder = map[card]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

var kindOrder2 = map[card]int{
	'J': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'Q': 11,
	'K': 12,
	'A': 13,
}

func calculateType(cards [5]card) int {
	counts := map[card]int{}
	for _, k := range cards {
		counts[k]++
	}
	groupCounts := [6]int{}
	for _, count := range counts {
		groupCounts[count]++
	}
	if groupCounts[5] == 1 {
		return FIVE_OF_A_KIND
	} else if groupCounts[4] == 1 {
		return FOUR_OF_A_KIND
	} else if groupCounts[3] == 1 {
		if groupCounts[2] == 1 {
			return FULL_HOUSE
		} else {
			return THREE_OF_A_KIND
		}
	} else if groupCounts[2] == 2 {
		return TWO_PAIR
	} else if groupCounts[2] == 1 {
		return ONE_PAIR
	} else {
		return HIGH_CARD
	}
}

func calculateType2(cards [5]card) int {
	counts := map[card]int{}
	maxCount := -1
	var maxCountCard card = -1
	for _, c := range cards {
		counts[c]++
		if counts[c] > maxCount && c != 'J' {
			maxCount = counts[c]
			maxCountCard = c
		}
	}
	fmt.Println(counts)
	for _, c := range cards {
		if c == 'J' {
			counts[maxCountCard]++
			counts['J']--
		}
	}
	fmt.Println(counts)
	groupCounts := [6]int{}
	for _, count := range counts {
		groupCounts[count]++
	}
	if groupCounts[5] == 1 {
		return FIVE_OF_A_KIND
	} else if groupCounts[4] == 1 {
		return FOUR_OF_A_KIND
	} else if groupCounts[3] == 1 {
		if groupCounts[2] == 1 {
			return FULL_HOUSE
		} else {
			return THREE_OF_A_KIND
		}
	} else if groupCounts[2] == 2 {
		return TWO_PAIR
	} else if groupCounts[2] == 1 {
		return ONE_PAIR
	} else {
		return HIGH_CARD
	}
}

func compareHands(h1 hand, h2 hand) bool {
	if h1.typ < h2.typ {
		return true
	} else if h1.typ > h2.typ {
		return false
	} else {
		for i := 0; i < len(h1.cards); i++ {
			if kindOrder[h1.cards[i]] < kindOrder[h2.cards[i]] {
				return true
			} else if kindOrder[h1.cards[i]] > kindOrder[h2.cards[i]] {
				return false
			}
		}
	}
	panic("shouldn't happen")
}

func compareHands2(h1 hand, h2 hand) bool {
	if h1.typ < h2.typ {
		return true
	} else if h1.typ > h2.typ {
		return false
	} else {
		for i := 0; i < len(h1.cards); i++ {
			if kindOrder2[h1.cards[i]] < kindOrder2[h2.cards[i]] {
				return true
			} else if kindOrder2[h1.cards[i]] > kindOrder2[h2.cards[i]] {
				return false
			}
		}
	}
	panic("shouldn't happen")
}

func part1(input io.Reader) {
	hands := []hand{}

	scanner := bufio.NewScanner(input)
	var err error
	for scanner.Scan() {
		line := scanner.Text()
		var handStr string
		var bid int
		_, err = fmt.Sscanf(line, "%s %d", &handStr, &bid)
		if err != nil {
			panic(err)
		}
		cards := [5]card{}
		for i := 0; i < 5; i++ {
			cards[i] = card(handStr[i])
		}
		hands = append(hands, hand{cards, 0, bid})
	}

	for i := range hands {
		hands[i].typ = calculateType(hands[i].cards)
	}

	fmt.Println(hands)

	sort.Slice(hands, func(i, j int) bool {
		return compareHands(hands[i], hands[j])
	})

	fmt.Println(hands)

	totalWinnings := 0
	for i, h := range hands {
		totalWinnings += h.bid * (i + 1)
	}

	fmt.Println(totalWinnings)
}

func part2(input io.Reader) {
	hands := []hand{}

	scanner := bufio.NewScanner(input)
	var err error
	for scanner.Scan() {
		line := scanner.Text()
		var handStr string
		var bid int
		_, err = fmt.Sscanf(line, "%s %d", &handStr, &bid)
		if err != nil {
			panic(err)
		}
		cards := [5]card{}
		for i := 0; i < 5; i++ {
			cards[i] = card(handStr[i])
		}
		hands = append(hands, hand{cards, 0, bid})
	}

	for i := range hands {
		hands[i].typ = calculateType2(hands[i].cards)
	}

	fmt.Println(hands)

	sort.Slice(hands, func(i, j int) bool {
		return compareHands2(hands[i], hands[j])
	})

	fmt.Println(hands)

	totalWinnings := 0
	for i, h := range hands {
		totalWinnings += h.bid * (i + 1)
	}

	fmt.Println(totalWinnings)
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	part1(file)
	file.Seek(0, 0)
	part2(file)
}
