package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func part1(input io.Reader) {
	var directions string
	type nextNode struct {
		left  string
		right string
	}
	var nodes = map[string]nextNode{}

	scanner := bufio.NewScanner(input)
	if scanner.Scan() {
		directions = scanner.Text()
	}

	nodeRegexp := regexp.MustCompile("([0-9A-Z]+)\\s*=\\s*\\(([0-9A-Z]+),\\s*([0-9A-Z]+)\\)")

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			match := nodeRegexp.FindStringSubmatch(line)
			if len(match) == 4 {
				if _, ok := nodes[match[1]]; ok {
					panic("found duplicate source node")
				}
				nodes[match[1]] = nextNode{match[2], match[3]}
			} else {
				panic("string does not match node regexp")
			}
		}
	}

	current := "AAA"
	stepCount := 0
outer:
	for {
		for i := 0; i < len(directions); i++ {
			if current == "ZZZ" {
				break outer
			}
			next := nodes[current]
			if directions[i] == 'L' {
				current = next.left
			} else {
				current = next.right
			}
			stepCount++
		}
	}

	fmt.Println(stepCount)
}

func part2(input io.Reader) {
	var directions string
	type nextNode struct {
		left  string
		right string
	}
	var nodes = map[string]nextNode{}

	scanner := bufio.NewScanner(input)
	if scanner.Scan() {
		directions = scanner.Text()
	}

	nodeRegexp := regexp.MustCompile("([0-9A-Z]+)\\s*=\\s*\\(([0-9A-Z]+),\\s*([0-9A-Z]+)\\)")

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			match := nodeRegexp.FindStringSubmatch(line)
			if len(match) == 4 {
				if _, ok := nodes[match[1]]; ok {
					panic("found duplicate source node")
				}
				nodes[match[1]] = nextNode{match[2], match[3]}
			} else {
				panic("string does not match node regexp")
			}
		}
	}

	startNodes := []string{}
	for node := range nodes {
		if strings.HasSuffix(node, "A") {
			startNodes = append(startNodes, node)
		}
	}

	stepCounts := map[string]int{}
	for _, start := range startNodes {
		current := start
		stepCount := 0
	outer:
		for {
			for i := 0; i < len(directions); i++ {
				if strings.HasSuffix(current, "Z") {
					break outer
				}
				next := nodes[current]
				if directions[i] == 'L' {
					current = next.left
				} else {
					current = next.right
				}
				stepCount++
			}
		}
		stepCounts[start] = stepCount
		if stepCount%len(directions) != 0 {
			panic("step count is expected to be divisible by directions length")
		}
	}

	maxStepCount := 0
	for _, count := range stepCounts {
		if count > maxStepCount {
			maxStepCount = count
		}
	}

	commonStepCount := 0
	ok := false
	for !ok {
		commonStepCount += maxStepCount
		ok = true
		for _, count := range stepCounts {
			if commonStepCount%count != 0 {
				ok = false
				break
			}
		}
	}

	fmt.Println(commonStepCount)
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
