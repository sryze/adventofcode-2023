package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func part1(input io.Reader) {
	bagConfig := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	possibleGameIds := make([]int, 0)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		gameParts := strings.Split(line, ":")
		gameId, err := strconv.Atoi(strings.Split(gameParts[0], " ")[1])
		if err != nil {
			panic(fmt.Sprintf("error parsing game ID: %v", err))
		}

		ok := true
	setLoop:
		for _, set := range strings.Split(gameParts[1], ";") {
			set := strings.TrimSpace(set)
			for _, cubeGroup := range strings.Split(set, ",") {
				cubeGroup := strings.Split(strings.TrimSpace(cubeGroup), " ")
				count, err := strconv.Atoi(cubeGroup[0])
				if err != nil {
					panic(fmt.Sprintf("error parsing cube count: %v", err))
				}
				color := cubeGroup[1]
				if count > bagConfig[color] {
					ok = false
					break setLoop
				}
			}
		}
		if ok {
			possibleGameIds = append(possibleGameIds, gameId)
		}
	}

	sum := 0
	for _, gameId := range possibleGameIds {
		sum += gameId
	}

	fmt.Printf("%d\n", sum)
}

func part2(input io.Reader) {
	powers := make([]int, 0)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		gameParts := strings.Split(line, ":")
		// gameId, err := strconv.Atoi(strings.Split(gameParts[0], " ")[1])
		// if err != nil {
		// 	panic(fmt.Sprintf("error parsing game ID: %v", err))
		// }

		minConfig := make(map[string]int)
		for _, set := range strings.Split(gameParts[1], ";") {
			set := strings.TrimSpace(set)
			for _, cubeGroup := range strings.Split(set, ",") {
				cubeGroup := strings.Split(strings.TrimSpace(cubeGroup), " ")
				count, err := strconv.Atoi(cubeGroup[0])
				if err != nil {
					panic(fmt.Sprintf("error parsing cube count: %v", err))
				}
				color := cubeGroup[1]
				if count > minConfig[color] {
					minConfig[color] = count
				}
			}
		}
		power := minConfig["red"] * minConfig["green"] * minConfig["blue"]
		powers = append(powers, power)
	}

	sum := 0
	for _, power := range powers {
		sum += power
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
