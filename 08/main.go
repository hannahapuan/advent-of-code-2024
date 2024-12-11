package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

const (
	filename    string = "input.txt"
	openVal     rune   = '.'
	antinodeVal rune   = '#'
	regex       string = "[a-zA-Z0-9]"
)

type cell struct {
	frequency             rune
	x, y                  int
	isAntenna, isAntinode bool
}

type change struct {
	dx, dy int
}

func main() {
	m, err := readInput(filename)
	if err != nil {
		os.Exit(1)
	}
	mf := flatten2dSlice(m)
	pairs := calcAntennaPairs(mf, mf)
	fmt.Println(mapToString(m))

	an := getAllAntinodes(pairs, m)
	uan := unique(an)

	mapWithAntinodes := updateMapWithAntinodes(m, uan)
	fmt.Println(mapToString(mapWithAntinodes))
	fmt.Println(len(uan))
}

// Reads the input file and initializes the grid and guard's starting state
func readInput(fname string) ([][]cell, error) {
	cells := make([][]cell, 0) // 2D array representing the grid

	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file:  %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var row []cell
	var i, j int

	// regex checks if the character is a-z, A-Z, 0-9
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex string: %s", err)
	}

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				if len(row) > 0 {
					cells = append(cells, row) // Append the last row to the grid
				}
				break
			}
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		if char == '\n' {
			// Start a new row when encountering a newline
			cells = append(cells, row)
			row = make([]cell, 0)
			j++
			i = 0
			continue
		}

		// add cell to grid
		currCell := cell{x: i, y: j, frequency: char}

		// antenna case
		if r.MatchString(string(char)) {
			currCell.isAntenna = true
		}

		row = append(row, currCell)
		i++
	}
	return cells, nil
}

// Converts the grid into a printable string representation with labeled coordinates
func mapToString(s [][]cell) string {
	var export string

	// Add X-axis labels
	export += "   " // Padding for Y-axis labels
	for x := 0; x < len(s[0]); x++ {
		export += fmt.Sprintf("%2d ", x)
	}
	export += "\n"

	for y, row := range s {
		// Add Y-axis label
		export += fmt.Sprintf("%2d ", y)

		// Add row content
		for _, cell := range row {
			export += fmt.Sprintf("%s  ", string(cell.frequency))
		}
		export += "\n"
	}
	return export
}
func calcAntennaPairs(m1 []cell, m2 []cell) map[cell][]change {
	pairs := make(map[cell][]change)

	for _, cell1 := range m1 {
		for _, cell2 := range m2 {
			if isSameCell(cell1, cell2) || !cell1.isAntenna || !cell2.isAntenna || cell1.frequency != cell2.frequency {
				continue
			}
			s := change{
				dx: cell2.x - cell1.x,
				dy: cell2.y - cell1.y,
			}
			_, ok1 := pairs[cell1]
			_, ok2 := pairs[cell2]
			if !ok1 {
				pairs[cell1] = make([]change, 0)
			}
			if !ok2 {
				pairs[cell2] = make([]change, 0)
			}
			pairs[cell1] = append(pairs[cell1], s)
			// pairs[cell2] = append(pairs[cell2], s)
		}
	}

	return pairs
}

func isSameCell(a, b cell) bool {
	return a.x == b.x && a.y == b.y
}

func flatten2dSlice(s [][]cell) []cell {
	fs := make([]cell, 0)

	for _, row := range s {
		fs = append(fs, row...)
	}

	return fs
}

func getAllAntinodes(pairs map[cell][]change, m [][]cell) []cell {
	antinodes := make([]cell, 0)
	for c, changes := range pairs {
		for _, sl := range changes {
			an := validAntinodes(c, sl, m)
			antinodes = append(antinodes, an...)
			// break
		}
		// break
	}
	return antinodes
}

func validAntinodes(c cell, s change, m [][]cell) []cell {
	antinodes := make([]cell, 0)

	if s.dx == 0 {
		// Prevent divide by zero; vertical line case
		rateOfChangeY := s.dy
		newY := c.y + (rateOfChangeY * 2)
		// newNegY := c.y - rateOfChangeY

		if inBounds(c.x, newY, m) {
			c1 := cell{
				frequency:  '#',
				x:          c.x,
				y:          newY,
				isAntinode: true,
			}
			antinodes = append(antinodes, c1)
		}

		// if inBounds(c.x, newNegY, m) {
		// 	c2 := cell{
		// 		frequency:  '#',
		// 		x:          c.x,
		// 		y:          newNegY,
		// 		isAntinode: true,
		// 	}
		// 	antinodes = append(antinodes, c2)
		// }
		return antinodes
	}

	newX := c.x + (s.dx * 2)
	newY := c.y + (s.dy * 2)

	if inBounds(newX, newY, m) {
		c1 := cell{
			frequency:  '#',
			x:          newX,
			y:          newY,
			isAntinode: true,
		}
		antinodes = append(antinodes, c1)
	}

	return antinodes
}

func unique(cs []cell) map[[2]int]bool {
	mcs := make(map[[2]int]bool)
	for _, c := range cs {
		mcs[[2]int{c.x, c.y}] = true
	}
	return mcs
}

// inBounds checks if a cell is within the maps's boundaries
func inBounds(x, y int, m [][]cell) bool {
	return x >= 0 && x < len(m) && y >= 0 && y < len(m[0])
}

func updateMapWithAntinodes(m [][]cell, antiNodes map[[2]int]bool) [][]cell {
	for antinode := range antiNodes {
		x, y := antinode[0], antinode[1]
		if inBounds(x, y, m) {
			m[y][x].isAntinode = true
			m[y][x].frequency = antinodeVal
		}
	}
	return m
}
