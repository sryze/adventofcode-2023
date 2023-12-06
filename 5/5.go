package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
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

func part1(input io.Reader) {
	seeds := []int{}
	type mapRange struct {
		src int
		dst int
		len int
	}
	type categoryMap struct {
		category    string
		ranges      []mapRange
		dstCategory string
	}
	maps := map[string]*categoryMap{}
	var curMap *categoryMap

	scanner := bufio.NewScanner(input)
	var err error
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "seeds:") {
			seeds, err = readArray(strings.Replace(line, "seeds:", "", 1))
			if err != nil {
				panic(err)
			}
		} else if strings.HasSuffix(line, "map:") {
			mapName := strings.TrimSpace(strings.Replace(line, "map:", "", 1))
			name := strings.Split(mapName, "-to-")
			if len(name) != 2 {
				panic(fmt.Sprint("invalid map name: ", mapName))
			}
			srcName := name[0]
			dstName := name[1]
			curMap = &categoryMap{
				category:    srcName,
				ranges:      []mapRange{},
				dstCategory: dstName,
			}
			maps[srcName] = curMap
		} else {
			if curMap != nil {
				if line == "" {
					curMap = nil
					continue
				}
				r, err := readArray(line)
				if err != nil {
					panic(err)
				}
				if len(r) != 3 {
					panic("range line contains more than 3 elements")
				}
				curMap.ranges = append(curMap.ranges, mapRange{r[1], r[0], r[2]})
			}
		}
	}

	seedLocations := map[int]int{}
	for _, seed := range seeds {
		cur := seed
		for curMap := maps["seed"]; curMap != nil; curMap = maps[curMap.dstCategory] {
			for i := 0; i < len(curMap.ranges); i++ {
				r := curMap.ranges[i]
				if cur >= r.src && cur <= r.src+r.len {
					cur = r.dst + (cur - r.src)
					break
				}
			}
			if curMap.dstCategory == "location" {
				seedLocations[seed] = cur
				break
			}
		}
	}

	minLoc, minLocSeed := math.MaxInt, -1
	for seed, loc := range seedLocations {
		if loc < minLoc {
			minLoc = loc
			minLocSeed = seed
		}
	}

	fmt.Println(minLocSeed, minLoc)
}

func part2(input io.Reader) {
	type seedRange struct {
		start int
		len   int
	}
	seedRanges := []seedRange{}
	type mapRange struct {
		src int
		dst int
		len int
	}
	type categoryMap struct {
		category    string
		ranges      []mapRange
		dstCategory string
	}
	maps := map[string]*categoryMap{}
	var curMap *categoryMap

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "seeds:") {
			seeds, err := readArray(strings.Replace(line, "seeds:", "", 1))
			if err != nil {
				panic(err)
			}
			if len(seeds)%2 != 0 {
				panic("seeds line contains an odd number of values")
			}
			for i := 0; i < len(seeds); i += 2 {
				seedRanges = append(seedRanges, seedRange{seeds[i], seeds[i+1]})
			}
		} else if strings.HasSuffix(line, "map:") {
			mapName := strings.TrimSpace(strings.Replace(line, "map:", "", 1))
			name := strings.Split(mapName, "-to-")
			if len(name) != 2 {
				panic(fmt.Sprint("invalid map name: ", mapName))
			}
			srcName := name[0]
			dstName := name[1]
			curMap = &categoryMap{
				category:    srcName,
				ranges:      []mapRange{},
				dstCategory: dstName,
			}
			maps[srcName] = curMap
		} else {
			if curMap != nil {
				if line == "" {
					curMap = nil
					continue
				}
				r, err := readArray(line)
				if err != nil {
					panic(err)
				}
				if len(r) != 3 {
					panic("range line contains more than 3 elements")
				}
				curMap.ranges = append(curMap.ranges, mapRange{r[1], r[0], r[2]})
			}
		}
	}

	curRanges := []seedRange{}
	for _, r := range seedRanges {
		curRanges = append(curRanges, r)
	}

	curMap = maps["seed"]
	for curMap != nil {
		newRanges := []seedRange{}
		rs := make([]seedRange, len(curRanges))
		copy(rs, curRanges)
		mrs := make([]mapRange, len(curMap.ranges))
		copy(mrs, curMap.ranges)

		for _, mr := range mrs {
			for j := 0; j < len(rs); j++ {
				r := rs[j]
				interStart := max(r.start, mr.src)
				interEnd := min(r.start+r.len, mr.src+mr.len)
				if interEnd > interStart {
					mapped := seedRange{
						mr.dst + (interStart - mr.src),
						interEnd - interStart,
					}
					newRanges = append(newRanges, mapped)
					if j < len(rs)-1 {
						rs = append(rs[:j], rs[j+1:]...)
					} else {
						rs = rs[:j]
					}
					j--
					if interStart > r.start {
						rs = append(rs, seedRange{r.start, interStart - r.start})
					}
					if interEnd < r.start+r.len {
						rs = append(rs, seedRange{interEnd, r.start + r.len - interEnd})
					}
				}
			}
		}
		if len(rs) > 0 {
			newRanges = append(newRanges, rs...)
		}

		curMap = maps[curMap.dstCategory]
		curRanges = newRanges
	}

	minLoc := math.MaxInt
	for _, locRange := range curRanges {
		if locRange.start < minLoc {
			minLoc = locRange.start
		}
	}

	fmt.Println(minLoc)
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
