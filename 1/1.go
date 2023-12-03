package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func part1(input io.Reader) {
	sum := 0

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		firstDigitIndex := strings.IndexFunc(line, isDigit)
		lastDigitIndex := strings.LastIndexFunc(line, isDigit)
		sum += (int(line[firstDigitIndex])-int('0'))*10 + (int(line[lastDigitIndex]) - int('0'))
	}

	fmt.Printf("%d\n", sum)
}

func part2(input io.Reader) {
	sum := 0

	digitWords := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		firstDigitIndex := strings.IndexFunc(line, isDigit)
		firstDigit := 0
		if firstDigitIndex >= 0 {
			firstDigit = int(line[firstDigitIndex]) - int('0')
		}
		for i, w := range digitWords {
			if index := strings.Index(line, w); index >= 0 && index < firstDigitIndex {
				firstDigitIndex = index
				firstDigit = i + 1
			}
		}
		lastDigitIndex := strings.LastIndexFunc(line, isDigit)
		lastDigit := 0
		if lastDigitIndex >= 0 {
			lastDigit = int(line[lastDigitIndex]) - int('0')
		}
		for i, w := range digitWords {
			if index := strings.LastIndex(line, w); index >= 0 && index > lastDigitIndex {
				lastDigitIndex = index
				lastDigit = i + 1
			}
		}
		sum += firstDigit*10 + lastDigit
	}

	fmt.Printf("%d\n", sum)
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
