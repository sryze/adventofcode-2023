package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func parseNumberList(s string) []int {
	xs := []int{}
	for _, x := range strings.Split(s, " ") {
		x := strings.TrimSpace(x)
		if x != "" {
			xVal, err := strconv.Atoi(x)
			if err != nil {
				panic(err)
			}
			xs = append(xs, xVal)
		}
	}
	return xs
}

func part1(input io.Reader) {
	sum := 0

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		numberLists := strings.Split(strings.Split(line, ":")[1], "|")
		winningNumbers := parseNumberList(numberLists[0])
		winningNumbersMap := make(map[int]bool)
		for _, x := range winningNumbers {
			winningNumbersMap[x] = true
		}
		haveNumbers := parseNumberList(numberLists[1])

		winCount := 0
		for _, x := range haveNumbers {
			if winningNumbersMap[x] {
				winCount++
			}
		}
		score := 0
		if winCount > 0 {
			score = 1 << (winCount - 1)
			sum += score
		}
	}

	fmt.Println(sum)
}

func part2(input io.Reader) {
	type cardData struct {
		numMatches int
		numCopies  int
	}
	cards := map[int]*cardData{}
	cardNumbers := []int{}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		card := strings.Split(line, ":")
		cardNum, err := strconv.Atoi(strings.TrimSpace(card[0][5:]))
		if err != nil {
			panic(err)
		}

		cardNumbers = append(cardNumbers, cardNum)

		numberLists := strings.Split(card[1], "|")
		winningNumbers := parseNumberList(numberLists[0])
		winningNumbersMap := make(map[int]bool)
		for _, x := range winningNumbers {
			winningNumbersMap[x] = true
		}
		haveNumbers := parseNumberList(numberLists[1])

		matchCount := 0
		for _, x := range haveNumbers {
			if winningNumbersMap[x] {
				matchCount++
			}
		}

		cards[cardNum] = &cardData{matchCount, 0}
	}

	sort.Ints(cardNumbers)

	for _, cardNum := range cardNumbers {
		card := cards[cardNum]
		for i := 1; i <= card.numMatches; i++ {
			cardToCopy := cards[cardNum+i]
			if cardToCopy != nil {
				cardToCopy.numCopies += (card.numCopies + 1)
			}
		}
	}

	sum := 0
	for _, card := range cards {
		sum += 1 + card.numCopies
	}

	fmt.Println(sum)
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
