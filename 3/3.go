package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func isDigit(c byte) bool {
	return unicode.IsDigit(rune(c))
}

func isSymbol(c byte) bool {
	return c != '.' && !isDigit(c)
}

func part1(input io.Reader) {
	var schematic []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		schematic = append(schematic, line)
	}

	var partNumbers []string
	for rowIndex, row := range schematic {
		numIndex := 0
		for {
			offset := strings.IndexFunc(row[numIndex:], unicode.IsDigit)
			if offset < 0 {
				break
			}

			numIndex += offset
			numStart := numIndex
			num := ""
			for numIndex < len(row) && isDigit(row[numIndex]) {
				num += string(row[numIndex])
				numIndex++
			}

		numCharsLoop:
			for i := numStart; i < numStart+len(num); i++ {
				for j := max(0, rowIndex-1); j < len(schematic) && j <= rowIndex+1 && j < len(schematic); j++ {
					for k := max(0, i-1); k <= i+1 && k < len(row); k++ {
						if isSymbol(schematic[j][k]) {
							partNumbers = append(partNumbers, num)
							break numCharsLoop
						}
					}
				}
			}
		}
	}

	sum := 0
	for _, partNum := range partNumbers {
		partNumInt, _ := strconv.Atoi(partNum)
		sum += partNumInt
	}

	fmt.Println(sum)
}

func part2(input io.Reader) {
	var schematic []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		schematic = append(schematic, line)
	}

	type gearPos struct {
		row, col int
	}

	var gearMap = make(map[gearPos][]int)
	for rowIndex, row := range schematic {
		numIndex := 0
		for {
			offset := strings.IndexFunc(row[numIndex:], unicode.IsDigit)
			if offset < 0 {
				break
			}

			numIndex += offset
			numStart := numIndex
			num := ""
			for numIndex < len(row) && isDigit(row[numIndex]) {
				num += string(row[numIndex])
				numIndex++
			}

		numCharsLoop:
			for i := numStart; i < numStart+len(num); i++ {
				for j := max(0, rowIndex-1); j < len(schematic) && j <= rowIndex+1 && j < len(schematic); j++ {
					for k := max(0, i-1); k <= i+1 && k < len(row); k++ {
						c := schematic[j][k]
						if c == '*' {
							numInt, _ := strconv.Atoi(num)
							gearMap[gearPos{j, k}] = append(gearMap[gearPos{j, k}], numInt)
							break numCharsLoop
						}
					}
				}
			}
		}
	}

	sum := 0
	for _, nums := range gearMap {
		if len(nums) == 2 {
			sum += nums[0] * nums[1]
		}
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
