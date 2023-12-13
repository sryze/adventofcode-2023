package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

type position struct {
	x int
	y int
}

const (
	directionLeft  = 'L'
	directionRight = 'R'
	directionUp    = 'U'
	directionDown  = 'D'
)

const (
	sideLeft  = 'L'
	sideRight = 'R'
)

const (
	loopDirectionUnknown          = 0
	loopDirectionClockwise        = 1
	loopDirectionCounterClockwise = 2
)

func readTiles(input io.Reader) []string {
	tiles := []string{}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		tiles = append(tiles, line)
	}

	rowLen := 0
	for _, row := range tiles {
		if rowLen == 0 {
			rowLen = len(row)
		} else if rowLen != len(row) {
			panic("unexpected length")
		}
	}

	return tiles
}

func findAnimal(tiles []string) position {
	pos := position{-1, -1}
	for i := range tiles {
		for j := range tiles[i] {
			if tiles[i][j] == 'S' {
				pos = position{j, i}
				break
			}
		}
	}
	return pos
}

func findLoop(tiles []string, startX, startY int, path *[]position) int {
	n := findLoopNext(tiles, startX-1, startY, directionLeft, 0, map[position]bool{}, path)
	if n != -1 {
		*path = append(*path, position{startX - 1, startY})
		return n
	}
	n = findLoopNext(tiles, startX, startY+1, directionDown, 0, map[position]bool{}, path)
	if n != -1 {
		*path = append(*path, position{startX, startY + 1})
		return n
	}
	n = findLoopNext(tiles, startX+1, startY, directionRight, 0, map[position]bool{}, path)
	if n != -1 {
		*path = append(*path, position{startX + 1, startY})
		return n
	}
	n = findLoopNext(tiles, startX, startY-1, directionUp, 0, map[position]bool{}, path)
	if n != -1 {
		*path = append(*path, position{startX, startY - 1})
		return n
	}
	return -1
}

func findLoopNext(tiles []string, x, y int, direction int, numSteps int, visited map[position]bool, path *[]position) int {
	if visited[position{x, y}] {
		return -1
	}
	if y < 0 || y >= len(tiles) {
		return -1
	}
	if x < 0 || x >= len(tiles[y]) {
		return -1
	}

	tile := tiles[y][x]
	numSteps++
	visited[position{x, y}] = true

	switch tile {
	case 'S':
		return numSteps + 1
	case '.':
		return -1
	case '|':
		switch direction {
		case directionUp:
			y -= 1
		case directionDown:
			y += 1
		default:
			return -1 // can't go that way
		}
	case '-':
		switch direction {
		case directionLeft:
			x -= 1
		case directionRight:
			x += 1
		default:
			return -1 // can't go that way
		}
	case 'L':
		switch direction {
		case directionDown:
			x += 1
			direction = directionRight
		case directionLeft:
			y -= 1
			direction = directionUp
		default:
			return -1 // can't go that way
		}
	case 'J':
		switch direction {
		case directionDown:
			x -= 1
			direction = directionLeft
		case directionRight:
			y -= 1
			direction = directionUp
		default:
			return -1 // can't go that way
		}
	case '7':
		switch direction {
		case directionRight:
			y += 1
			direction = directionDown
		case directionUp:
			x -= 1
			direction = directionLeft
		default:
			return -1 // can't go that way
		}
	case 'F':
		switch direction {
		case directionLeft:
			y += 1
			direction = directionDown
		case directionUp:
			x += 1
			direction = directionRight
		default:
			return -1 // can't go that way
		}
	}

	result := findLoopNext(tiles, x, y, direction, numSteps, visited, path)
	if result != -1 {
		*path = append(*path, position{x, y})
	}
	return result
}

func part1(input io.Reader) {
	tiles := readTiles(input)

	animalPos := findAnimal(tiles)

	var loop []position
	numSteps := findLoop(tiles, animalPos.x, animalPos.y, &loop)
	fmt.Println(numSteps / 2)
}

func isLoopTile(x, y int, loop []position) bool {
	for _, p := range loop {
		if p.x == x && p.y == y {
			return true
		}
	}
	return false
}

func isLoopTilePos(p position, loop []position) bool {
	return isLoopTile(p.x, p.y, loop)
}

func getTile(tiles []string, p position) byte {
	return tiles[p.y][p.x]
}

func getTileChar(tiles []string, p position) byte {
	if p.y < 0 || p.y >= len(tiles) || p.x < 0 || p.x >= len(tiles[p.y]) {
		return ' '
	}
	return tiles[p.y][p.x]
}

/*

func canPassPipe(direction int, pipe byte, seenPipes map[byte]bool) bool {
	switch direction {
	case directionLeft, directionRight:
		switch pipe {
		case 'L':
			return !seenPipes['7'] && !seenPipes['F']
		case 'J':
			return !seenPipes['7'] && !seenPipes['F']
		case '7':
			return !seenPipes['L'] && !seenPipes['J']
		case 'F':
			return !seenPipes['L'] && !seenPipes['J']
		default:
			return pipe != '|'
		}
	case directionUp, directionDown:
		switch pipe {
		case 'L':
			return !seenPipes['J'] && !seenPipes['7']
		case 'J':
			return !seenPipes['L'] && !seenPipes['F']
		case '7':
			return !seenPipes['L'] && !seenPipes['F']
		case 'F':
			return !seenPipes['J'] && !seenPipes['7']
		default:
			return pipe != '-'
		}
	default:
		return false
	}
}

func checkInsideLoop(tiles []string, loop []position, p position, knownOutside map[position]bool) bool {
	// fmt.Printf("Checking tile %c @ %v\n", getTileChar(tiles, p), p)

	if isLoopTile(p.x, p.y, loop) {
		// fmt.Println("    is part of loop")
		return true
	}

	if p.y < 0 || p.y >= len(tiles) || p.x < 0 || p.x >= len(tiles[p.y]) {
		// fmt.Println("    is out of bounds")
		return false
	}

	neighbors := []position{
		{p.x - 1, p.y},
		{p.x + 1, p.y},
		{p.x, p.y - 1},
		{p.x, p.y + 1},
	}
	for _, pNeighbor := range neighbors {
		if knownOutside[pNeighbor] && !isLoopTile(pNeighbor.x, pNeighbor.y, loop) {
			knownOutside[p] = true
			// fmt.Printf("    has neighbor that is outside: %v\n", pNeighbor)
			return false
		}
	}

	numLoopSides := 0
	seenPipes := map[byte]bool{}
	for i := p.x + 1; i < len(tiles[p.y]); i++ {
		if isLoopTile(i, p.y, loop) {
			pipe := tiles[p.y][i]
			if !canPassPipe(directionRight, pipe, seenPipes) {
				numLoopSides++
				break
			}
			// fmt.Printf("    can pass pipe %c @ %v\n", pipe, position{i, p.y})
			seenPipes[pipe] = true
		} else if knownOutside[position{i, p.y - 1}] || knownOutside[position{i, p.y + 1}] {
			break
		}
	}
	seenPipes = map[byte]bool{}
	for i := p.x - 1; i >= 0; i-- {
		if isLoopTile(i, p.y, loop) {
			pipe := tiles[p.y][i]
			if !canPassPipe(directionLeft, pipe, seenPipes) {
				numLoopSides++
				break
			}
			// fmt.Printf("    can pass pipe %c @ %v\n", pipe, position{i, p.y})
			seenPipes[pipe] = true
		} else if knownOutside[position{i, p.y - 1}] || knownOutside[position{i, p.y + 1}] {
			break
		}
	}
	seenPipes = map[byte]bool{}
	for i := p.y + 1; i < len(tiles); i++ {
		if isLoopTile(p.x, i, loop) {
			pipe := tiles[i][p.x]
			if !canPassPipe(directionDown, pipe, seenPipes) {
				numLoopSides++
				break
			}
			// fmt.Printf("    can pass pipe %c @ %v\n", pipe, position{p.x, i})
			seenPipes[pipe] = true
		} else if knownOutside[position{p.x - 1, i}] || knownOutside[position{p.x + 1, i}] {
			break
		}
	}
	seenPipes = map[byte]bool{}
	for i := p.y - 1; i >= 0; i-- {
		if isLoopTile(p.x, i, loop) {
			pipe := tiles[i][p.x]
			if !canPassPipe(directionUp, pipe, seenPipes) {
				numLoopSides++
				break
			}
			// fmt.Printf("    can pass pipe %c @ %v\n", pipe, position{p.x, i})
			seenPipes[pipe] = true
		} else if knownOutside[position{p.x - 1, i}] || knownOutside[position{p.x + 1, i}] {
			break
		}
	}

	if numLoopSides == 4 {
		return true
	} else {
		// fmt.Println("    is outside itself")
		knownOutside[p] = true
		return false
	}
}

func countInsideLoop(tiles []string, loop []position, knownOutside map[position]bool) int {
	count := 0
	for {
		c := 0
		for i := range tiles {
			for j := range tiles[i] {
				if isLoopTile(j, i, loop) {
					continue
				}
				if checkInsideLoop(tiles, loop, position{j, i}, knownOutside) {
					c++
				}
			}
		}
		if c == count {
			break // no neighbor tiles were updated
		}
		count = c
	}
	return count
}

*/

/*

func checkInsideLoop(tiles []string, loop []position, p position, knownOutside map[position]bool) bool {
	visited := map[position]bool{}
	return checkInsideLoopImpl(tiles, loop, p, knownOutside, visited)
}

func checkInsideLoopImpl(tiles []string, loop []position, p position, knownOutside map[position]bool,
	visited map[position]bool) bool {
	fmt.Printf("Checking tile %c @ %v\n", getTileChar(tiles, p), p)

	if p.y < 0 || p.y >= len(tiles) || p.x < 0 || p.x >= len(tiles[p.y]) {
		color.Red("    %v no - is out of bounds", p)
		return false
	}
	// if isLoopTilePos(p, loop) {
	// 	color.Green("    %v yes - is loop part", p)
	// 	return true
	// }
	if knownOutside[p] {
		color.Red("    %v no - is marked as outside", p)
		return false
	}

	{
		newVisited := map[position]bool{}
		for k, v := range visited {
			newVisited[k] = v
		}
		newVisited[p] = true
		visited = newVisited
	}

	neighbors := []struct {
		position
		blockingPipe byte
	}{
		{position{p.x - 1, p.y}, '|'},
		{position{p.x + 1, p.y}, '|'},
		{position{p.x, p.y - 1}, '-'},
		{position{p.x, p.y + 1}, '-'},
	}
	blockedCount := 0
	for _, neighbor := range neighbors {
		neighborPos := neighbor.position
		neighborTile := getTileChar(tiles, neighborPos)
		if visited[neighborPos] {
			blockedCount++
		} else if isLoopTilePos(neighborPos, loop) && neighborTile == neighbor.blockingPipe {
			blockedCount++
		} else if checkInsideLoopImpl(tiles, loop, neighborPos, knownOutside, visited) {
			blockedCount++
		}
	}

	if blockedCount < 4 {
		color.Red("    %v no - is not blocked on all sides", p)
		knownOutside[p] = true
		return false
	} else {
		color.Green("    %v yes", p)
		return true
	}
}

func countInsideLoop(tiles []string, loop []position, knownOutside map[position]bool) int {
	count := 0
	for {
		c := 0
		for i := range tiles {
			for j := range tiles[i] {
				if isLoopTile(j, i, loop) {
					continue
				}
				if checkInsideLoop(tiles, loop, position{j, i}, knownOutside) {
					c++
				}
			}
		}
		if c == count {
			break // no neighbor tiles were updated
		}
		count = c
	}
	return count
}

*/

func part2(input io.Reader) {
	tiles := readTiles(input)

	animalPos := findAnimal(tiles)

	var loop []position
	_ = findLoop(tiles, animalPos.x, animalPos.y, &loop)

	// fmt.Println("Loop:", loop)

	left, right := map[position]bool{}, map[position]bool{}

	type sidedNeighbor struct {
		side    int
		xOffset int
		yOffset int
	}
	neighborsByDirection := map[int]map[byte][]sidedNeighbor{
		directionRight: {
			'-': {{sideLeft, 0, -1}, {sideRight, 0, +1}},
			'J': {{sideLeft, 0, -1}, {sideRight, 0, +1}, {sideRight, +1, 0}},
			'7': {{sideLeft, 0, -1}, {sideRight, 0, +1}, {sideLeft, +1, 0}},
		},
		directionLeft: {
			'-': {{sideRight, 0, -1}, {sideLeft, 0, +1}},
			'L': {{sideRight, 0, -1}, {sideLeft, 0, +1}, {sideLeft, -1, 0}},
			'F': {{sideRight, 0, -1}, {sideLeft, 0, +1}, {sideRight, -1, 0}},
		},
		directionUp: {
			'|': {{sideLeft, -1, 0}, {sideRight, +1, 0}},
			'7': {{sideLeft, -1, 0}, {sideRight, +1, 0}, {sideRight, 0, -1}},
			'F': {{sideLeft, -1, 0}, {sideRight, +1, 0}, {sideLeft, 0, -1}},
		},
		directionDown: {
			'|': {{sideRight, -1, 0}, {sideLeft, +1, 0}},
			'L': {{sideRight, -1, 0}, {sideLeft, +1, 0}, {sideRight, 0, +1}},
			'J': {{sideRight, -1, 0}, {sideLeft, +1, 0}, {sideLeft, 0, +1}},
		},
	}

	for i := 1; i < len(loop); i++ {
		p := loop[i]
		pPrev := loop[i-1]
		direction := 0
		if p.x > pPrev.x {
			direction = directionRight
		} else if p.x < pPrev.x {
			direction = directionLeft
		} else if p.y < pPrev.y {
			direction = directionUp
		} else if p.y > pPrev.y {
			direction = directionDown
		}
		if direction == 0 {
			panic("could not determine direction")
		}
		neighbors := neighborsByDirection[direction][getTile(tiles, p)]
		for _, neighbor := range neighbors {
			nX := p.x + neighbor.xOffset
			nY := p.y + neighbor.yOffset
			var side map[position]bool
			switch neighbor.side {
			case sideLeft:
				side = left
			case sideRight:
				side = right
			}
			if nX > 0 && !isLoopTile(nX, nY, loop) {
				side[position{nX, nY}] = true
			}
		}
	}

	for {
		numTakenSide := 0
		for i := range tiles {
			for j := range tiles[i] {
				p := position{j, i}
				if isLoopTilePos(p, loop) {
					continue
				}
				if !left[p] && !right[p] {
					continue
				}
				side := left
				if right[p] {
					side = right
				}
				if i > 0 && !side[position{j, i - 1}] && !isLoopTile(j, i-1, loop) {
					side[position{j, i - 1}] = true
					numTakenSide++
				}
				if i < len(tiles)-1 && !side[position{j, i + 1}] && !isLoopTile(j, i+1, loop) {
					side[position{j, i + 1}] = true
					numTakenSide++
				}
				if j > 0 && !side[position{j - 1, i}] && !isLoopTile(j-1, i, loop) {
					side[position{j - 1, i}] = true
					numTakenSide++
				}
				if j < len(tiles[i])-1 && !side[position{j + 1, i}] && !isLoopTile(j+1, i, loop) {
					side[position{j + 1, i}] = true
					numTakenSide++
				}
			}
		}
		if numTakenSide == 0 {
			break
		}
	}

	white := color.New(color.FgWhite)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	// blue := color.New(color.FgBlue)
	// gray := color.New(color.FgHiBlack)
	green := color.New(color.FgGreen)
	magenta := color.New(color.FgMagenta)
	for i := range tiles {
		for j := range tiles[i] {
			p := position{j, i}
			tile := tiles[i][j]
			var tileColor = white
			// isOutside := outside[position{j, i}]
			if left[p] {
				tileColor = green
			} else if right[p] {
				if left[p] {
					panic("tile is right and left simultaneously")
				}
				tileColor = magenta
			} else if isLoopTile(j, i, loop) {
				tileColor = yellow
			} else {
				panic(fmt.Sprintf("tile %c @ %v is neither left nor right!", getTileChar(tiles, p), p))
			}
			// if isOutside {
			// 	tileColor = gray
			// 	tile = 'O'
			// } else {
			// 	tileColor = blue
			// 	tile = 'I'
			// }
			if tile == 'S' {
				tileColor = red
			}
			printedTile := rune(tile)
			switch tile {
			case '.', 'I', 'O':
				printedTile = '●'
			case 'S':
				printedTile = '★'
			case '-':
				printedTile = '━'
			case '|':
				printedTile = '┃'
			case 'L':
				printedTile = '┗'
			case 'J':
				printedTile = '┛'
			case '7':
				printedTile = '┓'
			case 'F':
				printedTile = '┏'
			}
			tileColor.Printf("%c", printedTile)
		}
		fmt.Println("")
	}

	loopDirection := loopDirectionUnknown
	for i := range tiles {
		for j := range tiles[i] {
			if i == 0 || j == 0 {
				p := position{j, i}
				if left[p] {
					loopDirection = loopDirectionClockwise
					break
				}
				if right[p] {
					loopDirection = loopDirectionCounterClockwise
					break
				}
			}
		}
	}
	insideCount := 0
	switch loopDirection {
	case loopDirectionClockwise:
		insideCount = len(right)
	case loopDirectionCounterClockwise:
		insideCount = len(left)
	default:
		insideCount = len(left) + len(right) // no tiles outside main loop
	}

	fmt.Println(insideCount)
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
