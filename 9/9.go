package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
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

func isAllZeros(seq []int) bool {
	for _, x := range seq {
		if x != 0 {
			return false
		}
	}
	return true
}

func part1(input io.Reader) {
	sequences := [][]int{}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		seq, err := readArray(line)
		if err != nil {
			panic(err)
		}
		sequences = append(sequences, seq)
	}

	extrapolatedSequences := [][]int{}
	for _, seq := range sequences {
		newSequences := [][]int{}
		curSeq := seq
		for {
			newSequences = append(newSequences, curSeq)
			if isAllZeros(curSeq) {
				break
			}
			newSeq := make([]int, len(curSeq)-1)
			for i := 0; i < len(curSeq)-1; i++ {
				diff := curSeq[i+1] - curSeq[i]
				newSeq[i] = diff
			}
			curSeq = newSeq
		}

		newSequences[len(newSequences)-1] = append(newSequences[len(newSequences)-1], 0)

		for i := len(newSequences) - 2; i >= 0; i-- {
			seq := newSequences[i]
			nextSeq := newSequences[i+1]
			diff := nextSeq[len(nextSeq)-1]
			newSequences[i] = append(seq, seq[len(seq)-1]+diff)
		}

		extrapolatedSequences = append(extrapolatedSequences, newSequences[0])
	}

	sum := 0
	for _, seq := range extrapolatedSequences {
		sum += seq[len(seq)-1]
	}

	fmt.Println(sum)
}

func part2(input io.Reader) {
	sequences := [][]int{}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		seq, err := readArray(line)
		if err != nil {
			panic(err)
		}
		sequences = append(sequences, seq)
	}

	extrapolatedSequences := [][]int{}
	for _, seq := range sequences {
		newSequences := [][]int{}
		curSeq := seq
		for {
			newSequences = append(newSequences, curSeq)
			if isAllZeros(curSeq) {
				break
			}
			newSeq := make([]int, len(curSeq)-1)
			for i := 0; i < len(curSeq)-1; i++ {
				diff := curSeq[i+1] - curSeq[i]
				newSeq[i] = diff
			}
			curSeq = newSeq
		}

		newSequences[len(newSequences)-1] = append(newSequences[len(newSequences)-1], 0)

		for i := len(newSequences) - 2; i >= 0; i-- {
			seq := newSequences[i]
			nextSeq := newSequences[i+1]
			diff := nextSeq[0]
			newSequences[i] = append([]int{seq[0] - diff}, seq...)
		}

		extrapolatedSequences = append(extrapolatedSequences, newSequences[0])
	}

	sum := 0
	for _, seq := range extrapolatedSequences {
		sum += seq[0]
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
