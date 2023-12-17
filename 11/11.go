package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

const (
	tileSpace  = '.'
	tileGalaxy = '#'
)

func readImage(input io.Reader) [][]byte {
	image := [][]byte{}
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		row := []byte(line)
		image = append(image, row)
	}
	return image
}

func rowHasGalaxies(image [][]byte, row int) bool {
	return slices.Contains(image[row], tileGalaxy)
}

func colHasGalaxies(image [][]byte, col int) bool {
	for i := 0; i < len(image); i++ {
		if image[i][col] == '#' {
			return true
		}
	}
	return false
}

func adjustImageForExpansion(sourceImage [][]byte) [][]byte {
	image := [][]byte{}
	for _, row := range sourceImage {
		newRow := make([]byte, len(row))
		copy(newRow, row)
		image = append(image, newRow)
		if !slices.Contains(row, '#') {
			newRow := make([]byte, len(row))
			copy(newRow, row)
			image = append(image, newRow)
		}
	}
	for j := 0; j < len(image[0]); j++ {
		if !colHasGalaxies(image, j) {
			for i := 0; i < len(image); i++ {
				row := image[i]
				newRow := make([]byte, len(row)+1)
				copy(newRow, row[:j+1])
				newRow[j+1] = tileSpace
				copy(newRow[j+1:], row[j:])
				image[i] = newRow
			}
			j++
		}
	}
	return image
}

func distance(from coordinate, to coordinate) int {
	return (from.row-to.row)*(from.row-to.row) + (from.col-to.col)*(from.col-to.col)
}

type coordinate struct {
	row, col int
}

func findShortestPath(image [][]byte, from coordinate, to coordinate) ([]coordinate, error) {
	var path []coordinate

	if from == to {
		return []coordinate{to}, nil
	}
	if from.row < 0 || from.row >= len(image) || from.col < 0 || from.col >= len(image[from.row]) {
		return path, nil
	}

	neighbors := []coordinate{
		{from.row - 1, from.col},
		{from.row + 1, from.col},
		{from.row, from.col - 1},
		{from.row, from.col + 1},
	}
	closestDistance := -1
	closestNeighborIndex := -1
	for i, neighbor := range neighbors {
		dist := distance(neighbor, to)
		if closestDistance == -1 || dist < closestDistance {
			closestDistance = dist
			closestNeighborIndex = i
		}
	}
	if closestNeighborIndex != -1 {
		var err error
		path, err = findShortestPath(image, neighbors[closestNeighborIndex], to)
		if err != nil {
			return path, err
		}
		path = append([]coordinate{from}, path...)
	}
	return path, nil
}

func printPath(image [][]byte, path []coordinate) {
	fmt.Println(path)
	for i := 0; i < len(image); i++ {
		for j := 0; j < len(image[i]); j++ {
			coord := coordinate{i, j}
			if coord == path[0] {
				fmt.Print("A")
			} else if coord == path[len(path)-1] {
				fmt.Print("B")
			} else if slices.Contains(path, coord) && image[i][j] != tileGalaxy {
				fmt.Print("*")
			} else {
				fmt.Printf("%c", image[i][j])
			}
		}
		fmt.Println("")
	}
}

func part1(input io.Reader) {
	image := adjustImageForExpansion(readImage(input))

	galaxies := []coordinate{}
	for i := 0; i < len(image); i++ {
		for j := 0; j < len(image[i]); j++ {
			if image[i][j] == tileGalaxy {
				galaxies = append(galaxies, coordinate{i, j})
			}
		}
	}

	// for _, row := range image {
	// 	fmt.Println(string(row))
	// }

	// fmt.Println(galaxies)

	sum := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			path, err := findShortestPath(image, galaxies[i], galaxies[j])
			if err != nil {
				panic(err)
			}
			pathLen := len(path) - 1
			// fmt.Printf("Shortest path from %v to %v is %d steps:\n",
			// 	galaxies[i], galaxies[j], pathLen)
			// printPath(image, path)
			sum += pathLen
		}
	}

	fmt.Println(sum)
}

func part2(input io.Reader) {
	image := readImage(input)
	emptyRows := map[int]bool{}
	emptyCols := map[int]bool{}

	galaxies := []coordinate{}
	for i := 0; i < len(image); i++ {
		for j := 0; j < len(image[i]); j++ {
			if image[i][j] == tileGalaxy {
				galaxies = append(galaxies, coordinate{i, j})
			}
		}
	}

	sum := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			path, err := findShortestPath(image, galaxies[i], galaxies[j])
			if err != nil {
				panic(err)
			}
			pathLen := 0
			for coordIdx, coord := range path {
				if coordIdx == 0 {
					continue
				}
				prevCoord := path[coordIdx-1]
				if coord.row != prevCoord.row {
					isEmptyRow, ok := emptyRows[coord.row]
					if !ok {
						isEmptyRow = !rowHasGalaxies(image, coord.row)
						emptyRows[coord.row] = isEmptyRow
					}
					if isEmptyRow {
						pathLen += 1000000
					} else {
						pathLen += 1
					}
				} else {
					isEmptyCol, ok := emptyCols[coord.col]
					if !ok {
						isEmptyCol = !colHasGalaxies(image, coord.col)
						emptyCols[coord.col] = isEmptyCol
					}
					if isEmptyCol {
						pathLen += 1000000
					} else {
						pathLen += 1
					}
				}
			}
			sum += pathLen
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
