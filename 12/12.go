package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type placement struct {
	pattern string
	end     int
}

func canPlaceGroup(pattern string, index int, size int) bool {
	if index > 0 && (pattern[index-1] == 'X' || pattern[index-1] == '#') {
		return false
	}

	patternLen := len(pattern)
	if index+size > patternLen {
		return false
	}
	if index+size < patternLen && (pattern[index+size] == 'X' || pattern[index+size] == '#') {
		return false
	}

	matchCount := 0
	for cur := index; cur < patternLen; cur++ {
		if matchCount == size {
			break
		}
		if pattern[cur] != '#' && pattern[cur] != '?' {
			break
		}
		matchCount++
	}
	return matchCount == size
}

func generatePlacements(pattern string, start int, groupSizes []int, groupIndex int) []placement {
	placements := []placement{}

	restSize := 0
	for i := groupIndex + 1; i < len(groupSizes); i++ {
		restSize += groupSizes[i]
	}

	groupSize := groupSizes[groupIndex]
	for i := start; i < len(pattern)-restSize; i++ {
		if canPlaceGroup(pattern, i, groupSize) {
			// fmt.Println("Can place group of", size, "at", i, "in", pattern)
			tempPattern := []byte(strings.Clone(pattern))
			if i > 0 {
				if tempPattern[i-1] == 'X' || tempPattern[i-1] == '#' {
					panic("unexpected preceding symbol")
				}
				tempPattern[i-1] = '.'
			}
			for j := 0; j < groupSize; j++ {
				switch tempPattern[i+j] {
				case '?', '#':
					tempPattern[i+j] = 'X'
				case '.':
					panic("unexpected element while placing group")
				}
			}
			if i+groupSize < len(tempPattern) {
				switch tempPattern[i+groupSize] {
				case '?', '.':
					tempPattern[i+groupSize] = '.'
				case '#', 'X':
					panic("unexpected following symbol")
				}
			}
			placements = append(placements, placement{
				string(tempPattern),
				i + groupSize,
			})
		}
	}
	// fmt.Println("Generated placements:", placements)
	return placements
}

func generateAllPlacements(pattern string, start int, groupSizes []int, groupIndex int, placements *[]string) {
	// fmt.Println("Processing group @", index, "-", size)
	for _, p := range generatePlacements(pattern, start, groupSizes, groupIndex) {
		if groupIndex+1 < len(groupSizes) {
			generateAllPlacements(p.pattern, p.end, groupSizes, groupIndex+1, placements)
		} else {
			if strings.Contains(p.pattern, "#") {
				continue // inserted all groups but there are extra ones over from input pattern
			}
			*placements = append(*placements, p.pattern)
		}
	}
}

func finalizePattern(pattern string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(pattern, "X", "#"),
		"?",
		".")
}

func part1(input io.Reader) {
	type record struct {
		condition         string
		damagedGroupSizes []int
	}
	records := []record{}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		var condition string
		var damagedGroupSizesStr string
		_, err := fmt.Sscanf(line, "%s %s", &condition, &damagedGroupSizesStr)
		if err != nil {
			panic(err)
		}
		damagedGroupSizes := []int{}
		for _, sizeStr := range strings.Split(damagedGroupSizesStr, ",") {
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				panic(err)
			}
			damagedGroupSizes = append(damagedGroupSizes, size)
		}
		records = append(records, record{condition, damagedGroupSizes})
	}

	sum := 0
	for _, r := range records {
		// fmt.Println("Processing record:", r)

		placements := []string{}
		generateAllPlacements(r.condition, 0, r.damagedGroupSizes, 0, &placements)
		// fmt.Printf("Possible arrangements (%d):\n", len(placements))
		// for _, p := range placements {
		// 	fmt.Println(finalizePattern(p))
		// }
		sum += len(placements)

		for _, p := range placements {
			groups := strings.FieldsFunc(p, func(c rune) bool { return c == '.' || c == '?' })
			if len(groups) != len(r.damagedGroupSizes) {
				panic(fmt.Sprintf("unexpected number of damaged spring groups: expected %d, got %d,",
					len(r.damagedGroupSizes), len(groups)))
			}
			for i, group := range groups {
				if len(group) != r.damagedGroupSizes[i] {
					panic(fmt.Sprintf("unexpected group length: expected %d, got %d",
						r.damagedGroupSizes[i], len(group)))
				}
			}
		}
	}

	fmt.Println(sum)
}

func countAllPlacements2(
	pattern string, start int, groupSizes []int, groupIndex int, groupCache map[int][]int, count *int) {
	restSize := 0
	for i := groupIndex + 1; i < len(groupSizes); i++ {
		restSize += groupSizes[i] + 1
	}

	groupSize := groupSizes[groupIndex]
	patternLen := len(pattern)

	processIndex := func(i int) {
		tempPattern := []byte(strings.Clone(pattern))
		if i > 0 {
			if tempPattern[i-1] == 'X' || tempPattern[i-1] == '#' {
				panic("unexpected preceding symbol")
			}
			tempPattern[i-1] = '.'
		}
		for j := 0; j < groupSize; j++ {
			switch tempPattern[i+j] {
			case '?', '#':
				tempPattern[i+j] = 'X'
			case '.':
				panic("unexpected element while placing group")
			}
		}
		if i+groupSize < patternLen {
			switch tempPattern[i+groupSize] {
			case '?', '.':
				tempPattern[i+groupSize] = '.'
			case '#', 'X':
				panic("unexpected following symbol")
			}
		}

		tempPatternStr := string(tempPattern)
		if groupIndex+1 < len(groupSizes) {
			countAllPlacements2(tempPatternStr, i+groupSize+1, groupSizes, groupIndex+1, groupCache, count)
		} else {
			if strings.Contains(tempPatternStr, "#") {
				return // inserted all groups but there are extra ones over from input pattern
			}
			*count++
		}
	}

	if groupCache[groupSize] != nil {
		for _, i := range groupCache[groupSize] {
			if i >= patternLen-restSize {
				break
			}
			if i >= start {
				processIndex(i)
			}
		}
	} else {
		groupCache[groupSize] = []int{}
		for i := start; i < patternLen; i++ {
			if canPlaceGroup(pattern, i, groupSize) {
				groupCache[groupSize] = append(groupCache[groupSize], i)
			}
		}
		for _, i := range groupCache[groupSize] {
			if i >= patternLen-restSize {
				break
			}
			if i >= start {
				processIndex(i)
			}
		}
	}
}

func part2(input io.Reader, parallel bool, noUnfold bool) {
	type record struct {
		condition         string
		damagedGroupSizes []int
	}
	records := []record{}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		var foldedCondition string
		var foldedDamagedGroupSizesStr string
		_, err := fmt.Sscanf(line, "%s %s", &foldedCondition, &foldedDamagedGroupSizesStr)
		if err != nil {
			panic(err)
		}
		condition := foldedCondition
		damagedGroupSizesStr := foldedDamagedGroupSizesStr
		if !noUnfold {
			for i := 0; i < 4; i++ {
				condition += "?" + foldedCondition
				damagedGroupSizesStr += "," + foldedDamagedGroupSizesStr
			}
		}
		damagedGroupSizes := []int{}
		for _, sizeStr := range strings.Split(damagedGroupSizesStr, ",") {
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				panic(err)
			}
			damagedGroupSizes = append(damagedGroupSizes, size)
		}
		records = append(records, record{condition, damagedGroupSizes})
	}

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	sum := 0

	process := func(i int, r record) {
		fmt.Printf("Processing record %d: %v\n", i, r)
		if parallel {
			defer wg.Done()
		}

		count := 0
		countAllPlacements2(r.condition, 0, r.damagedGroupSizes, 0, map[int][]int{}, &count)

		fmt.Printf("Possible arrangements for record %d: %d\n", i, count)
		mutex.Lock()
		sum += count
		mutex.Unlock()
	}

	for i, r := range records {
		if parallel {
			wg.Add(1)
			go process(i, r)
		} else {
			process(i, r)
		}
	}
	if parallel {
		wg.Wait()
	}

	fmt.Println(sum)
}

func part2Clever(input io.Reader, parallel bool) {
	type record struct {
		condition         string
		damagedGroupSizes []int
	}
	records := []record{}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		var condition string
		var damagedGroupSizesStr string
		_, err := fmt.Sscanf(line, "%s %s", &condition, &damagedGroupSizesStr)
		if err != nil {
			panic(err)
		}
		damagedGroupSizes := []int{}
		for _, sizeStr := range strings.Split(damagedGroupSizesStr, ",") {
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				panic(err)
			}
			damagedGroupSizes = append(damagedGroupSizes, size)
		}
		records = append(records, record{condition, damagedGroupSizes})
	}

	sum := 0

	for i, r := range records {
		fmt.Printf("Processing record %d: %v\n", i, r)

		// placements := []string{}
		// generateAllPlacements(r.condition, 0, r.damagedGroupSizes, 0, &placements)
		// count1 := len(placements)
		// fmt.Println("count1:", count1)

		// headPatterns := map[string]bool{}
		// tailPatterns := map[string]bool{}
		// for _, p := range placements {
		// 	i := 0
		// 	for i < len(p) && p[i] != 'X' {
		// 		i++
		// 	}
		// 	head := p[0:i]
		// 	if strings.Contains(head, "?") || strings.Contains(head, "#") {
		// 		headPatterns[head] = true
		// 	}
		// 	i = len(p) - 1
		// 	for i > 0 && p[i] != 'X' {
		// 		i--
		// 	}
		// 	tail := p[i+1:]
		// 	if strings.Contains(tail, "?") || strings.Contains(tail, "#") {
		// 		tailPatterns[tail] = true
		// 	}
		// }
		// headPatterns[""] = true
		// tailPatterns[""] = true
		// for p := range headPatterns {
		// 	for q := range headPatterns {
		// 		if p != q && strings.Contains(p, q) {
		// 			headPatterns[q] = false
		// 		}
		// 	}
		// }
		// for p := range tailPatterns {
		// 	for q := range tailPatterns {
		// 		if p != q && strings.Contains(p, q) {
		// 			tailPatterns[q] = false
		// 		}
		// 	}
		// }

		// // placements = []string{}
		// // generateAllPlacements(r.condition, 0, r.damagedGroupSizes, 0, &placements)
		// // count2New := len(placements)
		// count2Head := 0
		// for p, ok := range headPatterns {
		// 	if !ok {
		// 		continue
		// 	}
		// 	// fmt.Println(p)
		// 	placements = []string{}
		// 	generateAllPlacements(r.condition+"?"+p, 0, r.damagedGroupSizes, 0, &placements)
		// 	count := len(placements)
		// 	count2Head += count
		// }
		// count2Tail := 0
		// for p, ok := range tailPatterns {
		// 	if !ok {
		// 		continue
		// 	}
		// 	// fmt.Println(p)
		// 	placements = []string{}
		// 	generateAllPlacements(p+"?"+r.condition, 0, r.damagedGroupSizes, 0, &placements)
		// 	count := len(placements)
		// 	count2Tail += count
		// }
		// count2New := count2Tail * count2Head
		// fmt.Println("count2New:", count2New)

		// placements2 := []string{}
		pattern2 := r.condition + "?" + r.condition
		groupSizes2 := append(r.damagedGroupSizes, r.damagedGroupSizes...)
		// generateAllPlacements(pattern2, 0, groupSizes2, 0, &placements2)
		// count2 := len(placements2)
		// fmt.Println("count2:", count2)

		// placements3 := []string{}
		pattern3 := pattern2 + "?" + r.condition
		groupSizes3 := append(groupSizes2, r.damagedGroupSizes...)
		// generateAllPlacements(pattern3, 0, groupSizes3, 0, &placements3)
		// count3 := len(placements3)
		// fmt.Println("count3:", count2)

		pattern4 := pattern3 + "?" + r.condition
		groupSizes4 := append(groupSizes3, r.damagedGroupSizes...)
		pattern5 := pattern4 + "?" + r.condition
		groupSizes5 := append(groupSizes4, r.damagedGroupSizes...)

		// fmt.Println("count2/count1:", count2/count1)
		// fmt.Println("count3/count2:", count3/count2)
		// fmt.Println("count3/count1^2:", count3/count1/count1)

		// // a := count2 - (count1 * count1)
		// // fmt.Println("count5:", count1*count1*count1*count1*count1+(a*a*a*a))

		// // count := count3 * (count3 / count2) * (count3 / count2)
		// count := count2 * (count2 / count1) * (count2 / count1) * (count2 / count1)

		count := 0
		countAllPlacements2(pattern5, 0, groupSizes5, 0, map[int][]int{}, &count)

		fmt.Printf("Possible arrangements for record %d: %d\n", i, count)
		sum += count
	}

	fmt.Println(sum)
}

func findPlacementPositions(pattern string, groupSize int) []int {
	positions := []int{}
	for i := 0; i < len(pattern); i++ {
		if canPlaceGroup(pattern, i, groupSize) {
			positions = append(positions, i)
		}
	}
	return positions
}

func countAllPlacements3Internal(pattern string, start int, groupSizes []int, groupIndex int, groupPosLists [][]int) int {
	count := 0
	groupSize := groupSizes[groupIndex]
	posList := groupPosLists[groupSize]
	posIndex := 0
	remainingSize := 0
	for i := groupIndex + 1; i < len(groupSizes); i++ {
		remainingSize += groupSizes[i] + 1
	}
	for i := 0; i < len(pattern)-remainingSize && posIndex < len(posList); i++ {
		// fmt.Println("groupIndex =", groupIndex, "groupSize =", groupSize, "i=", i)
		pos := posList[posIndex]
		if i == pos {
			if i >= start {
				// fmt.Println("pos:", pos)
				if groupIndex == len(groupSizes)-1 {
					needMore := false
					for j := i + groupSize; j < len(pattern); j++ {
						if pattern[j] == '#' {
							needMore = true
							// fmt.Println("need more!")
							break
						}
					}
					if !needMore {
						count += 1
						// fmt.Println("Success!")
						// fmt.Println(groupSize, "placed at", pos)
					}
				} else {
					x := countAllPlacements3Internal(pattern, pos+groupSize+1, groupSizes, groupIndex+1, groupPosLists)
					// if x > 0 {
					// 	fmt.Println(groupSize, "placed at", pos)
					// }
					count += x
				}
			}
			posIndex++
		}
		if i >= start && pattern[i] == '#' {
			// fmt.Println("found # after pos ", start, "for group", groupIndex)
			break
		}
	}
	// for _, pos := range groupPosLists[groupSize] {
	// 	if pos >= start {
	// 		if slices.Contains([]byte(pattern[start:pos]), '#') {
	// 			continue
	// 		}
	// 		if groupIndex == len(groupSizes)-1 {
	// 			if !slices.Contains([]byte(pattern[pos+groupSize:]), '#') {
	// 				count += 1
	// 				// fmt.Println(groupSize, "placed at", pos)
	// 			}
	// 		} else {
	// 			x := countAllPlacements3Internal(pattern, pos+groupSize+1, groupSizes, groupIndex+1, groupPosLists)
	// 			// if x > 0 {
	// 			// 	fmt.Println(groupSize, "placed at", pos)
	// 			// }
	// 			count += x
	// 		}
	// 	}
	// }
	return count
}

func countAllPlacements3(pattern string, groupSizes []int) int {
	groupPosLists := make([][]int, slices.Max(groupSizes)+1)
	for _, groupSize := range groupSizes {
		groupPosLists[groupSize] = findPlacementPositions(pattern, groupSize)
	}
	// fmt.Println(groupPosLists)
	return countAllPlacements3Internal(pattern, 0, groupSizes, 0, groupPosLists)
}

func part2Faster(input io.Reader, parallel bool, noUnfold bool) {
	type record struct {
		condition         string
		damagedGroupSizes []int
	}
	records := []record{}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		var foldedCondition string
		var foldedDamagedGroupSizesStr string
		_, err := fmt.Sscanf(line, "%s %s", &foldedCondition, &foldedDamagedGroupSizesStr)
		if err != nil {
			panic(err)
		}
		condition := foldedCondition
		damagedGroupSizesStr := foldedDamagedGroupSizesStr
		if !noUnfold {
			for i := 0; i < 4; i++ {
				condition += "?" + foldedCondition
				damagedGroupSizesStr += "," + foldedDamagedGroupSizesStr
			}
		}
		damagedGroupSizes := []int{}
		for _, sizeStr := range strings.Split(damagedGroupSizesStr, ",") {
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				panic(err)
			}
			damagedGroupSizes = append(damagedGroupSizes, size)
		}
		records = append(records, record{condition, damagedGroupSizes})
	}

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	sum := 0

	process := func(i int, r record) {
		// if i != 1 {
		// 	continue
		// }
		fmt.Printf("Processing record %d: %v\n", i, r)
		if parallel {
			defer wg.Done()
		}

		count := countAllPlacements3(r.condition, r.damagedGroupSizes)

		fmt.Printf("Possible arrangements for record %d: %d\n", i, count)
		mutex.Lock()
		sum += count
		mutex.Unlock()
	}

	for i, r := range records {
		if parallel {
			wg.Add(1)
			go process(i, r)
		} else {
			process(i, r)
		}
	}
	if parallel {
		wg.Wait()
	}

	fmt.Println(sum)
}

func part2Clever2(input io.Reader) {
	type record struct {
		condition         string
		damagedGroupSizes []int
	}
	records := []record{}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		var condition string
		var damagedGroupSizesStr string
		_, err := fmt.Sscanf(line, "%s %s", &condition, &damagedGroupSizesStr)
		if err != nil {
			panic(err)
		}
		damagedGroupSizes := []int{}
		for _, sizeStr := range strings.Split(damagedGroupSizesStr, ",") {
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				panic(err)
			}
			damagedGroupSizes = append(damagedGroupSizes, size)
		}
		records = append(records, record{condition, damagedGroupSizes})
	}

	sum := 0

	for i, r := range records {
		fmt.Printf("Processing record %d: %v\n", i, r)

		groupSizes1 := r.damagedGroupSizes
		count1 := countAllPlacements3(r.condition, groupSizes1)
		groupSizes2 := append(groupSizes1, groupSizes1...)
		count2 := countAllPlacements3(r.condition+"?"+r.condition, groupSizes2)
		groupSizes3 := append(groupSizes2, groupSizes1...)
		count3 := countAllPlacements3(r.condition+"?"+r.condition+"?"+r.condition, groupSizes3)
		groupSizes4 := append(groupSizes3, groupSizes1...)
		count4 := countAllPlacements3(r.condition+"?"+r.condition+"?"+r.condition+"?"+r.condition, groupSizes4)
		groupSizes5 := append(groupSizes4, groupSizes1...)
		count5 := countAllPlacements3(r.condition+"?"+r.condition+"?"+r.condition+"?"+r.condition+"?"+r.condition, groupSizes5)

		fmt.Println(count1, count2, count3, count4, count5, count2/count1, count2%count1 == 0, count3/count2, count3%count2 == 0, count4/count3, count4%count3 == 0, count5/count4, count5%count4 == 0)
		fmt.Println(count1*count2/count1*count2/count1*count2/count1*count2/count1, count5)

		count5Fast := count1 * count2 / count1 * count2 / count1 * count2 / count1 * count2 / count1
		// count5Fast := count4 / count3 * count4

		// fmt.Printf("Possible arrangements for record %d: %d\n", i, count)
		sum += count5Fast

		if count5 != count5Fast {
			panic("not equal")
		}
	}

	fmt.Println(sum)
}

func main() {
	partFlag := flag.Int("part", 1, "part number (1 or 2)")
	parallelFlag := flag.Bool("parallel", false, "use parallel computation if possible")
	noUnfoldFlag := flag.Bool("no-unfold", false, "don't unfold records")
	flag.Parse()

	_ = *parallelFlag
	_ = *noUnfoldFlag

	switch *partFlag {
	case 1:
		part1(os.Stdin)
	case 2:
		// part2(os.Stdin, *parallelFlag, *noUnfoldFlag)
		// part2Clever(os.Stdin)
		part2Faster(os.Stdin, *parallelFlag, *noUnfoldFlag)
		// part2Clever2(os.Stdin)
	default:
		panic("unsupported part number")
	}
}
