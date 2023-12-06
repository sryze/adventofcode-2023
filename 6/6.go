package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func readArray(s string) ([]int, error) {
	arr := []int{}
	r := bytes.NewReader([]byte(s))
	for {
		var x int
		n, err := fmt.Fscanf(r, "%d", &x)
		if err != nil && err != io.EOF {
			return arr, err
		}
		if n == 0 {
			break
		}
		arr = append(arr, x)
	}
	return arr, nil
}

func calculateDistance(holdTime int, raceTime int) int {
	if holdTime == 0 {
		return 0
	}
	// dist := 0
	// speed := holdTime
	// for t := holdTime; t < raceTime; t++ {
	// 	dist += speed
	// }
	// return dist
	return (raceTime - holdTime) * holdTime
}

func part1(input io.Reader) {
	times := []int{}
	distances := []int{}

	scanner := bufio.NewScanner(input)
	var err error
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Time:") {
			times, err = readArray(strings.Replace(line, "Time:", "", 1))
			if err != nil {
				panic(err)
			}
		}
		if strings.HasPrefix(line, "Distance:") {
			distances, err = readArray(strings.Replace(line, "Distance:", "", 1))
			if err != nil {
				panic(err)
			}
		}
	}

	if len(times) != len(distances) {
		panic("time and distance array lengths don't match")
	}

	waysToWin := map[int]int{}
	for i := 0; i < len(times); i++ {
		time := times[i]
		prevMaxDistance := distances[i]
		for t := 0; t <= time; t++ {
			d := calculateDistance(t, time)
			if d > prevMaxDistance {
				waysToWin[i]++
			}
		}
	}

	result := 1
	for _, n := range waysToWin {
		result *= n
	}

	fmt.Println(result)
}

func part2(input io.Reader) {
	time := 0
	distance := 0

	scanner := bufio.NewScanner(input)
	var err error
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Time:") {
			time, err = strconv.Atoi(strings.ReplaceAll(strings.Replace(line, "Time:", "", 1), " ", ""))
			if err != nil {
				panic(err)
			}
		}
		if strings.HasPrefix(line, "Distance:") {
			distance, err = strconv.Atoi(strings.ReplaceAll(strings.Replace(line, "Distance:", "", 1), " ", ""))
			if err != nil {
				panic(err)
			}
		}
	}

	waysToWin := 0
	for t := 0; t <= time; t++ {
		d := calculateDistance(t, time)
		if d > distance {
			waysToWin++
		}
	}

	fmt.Println(waysToWin)
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
