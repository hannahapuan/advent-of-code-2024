package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

const (
	filename    string = "example.txt"
	openVal     rune   = '.'
	antinodeVal rune   = '#'
	regex       string = "[a-zA-Z0-9]"
)

type cell struct {
	frequency             rune
	x, y                  int
	isAntenna, isAntinode bool
	// matchingAntennas []cell
}

type antennaPair struct {
	a1X, a1Y int
	a2X, a2Y int
}

type slope struct {
	dx, dy int
}

// type antennaPair struct {
// 	a1 cell
// 	a2 cell
// }

func main() {
	m, atl, err := readInput(filename)
	if err != nil {
		os.Exit(1)
	}
	mf := flatten2dSlice(m)
	pairs := calcAntennaPairs(mf, mf)

	fmt.Println(mapToString(m))
	fmt.Println(atl)
	// fmt.Println(apToSlopeToString(aps))
	fmt.Println(pairs)
	an := getAllAntinodes(pairs, m)
	uan := unique(an)
	fmt.Println(uan)
	fmt.Println(len(uan))

}

// Reads the input file and initializes the grid and guard's starting state
func readInput(fname string) ([][]cell, map[rune][]cell, error) {
	antennaToLocation := make(map[rune][]cell)

	cells := make([][]cell, 0) // 2D array representing the grid

	file, err := os.Open(fname)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file:  %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var row []cell
	var i, j int

	// regex checks if the character is a-z, A-Z, 0-9
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, nil, fmt.Errorf("error compiling regex string: %s", err)
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
			return nil, nil, fmt.Errorf("error reading file: %w", err)
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
			_, ok := antennaToLocation[char]
			if !ok {
				antennaToLocation[char] = make([]cell, 0)
			}
			antennaToLocation[char] = append(antennaToLocation[char], currCell)
		}

		row = append(row, currCell)
		i++
	}
	return cells, antennaToLocation, nil
}

// Converts the grid into a printable string representation
func mapToString(s [][]cell) string {
	var export string
	for _, row := range s {
		for _, cell := range row {
			export += fmt.Sprintf("%s ", string(cell.frequency))
		}
		export += "\n"
	}
	return export
}

func calcAntennaPairs(m1 []cell, m2 []cell) map[cell][]slope {
	// aps := make(map[antennaPair]slope)
	pairs := make(map[cell][]slope)

	for _, cell1 := range m1 {
		for _, cell2 := range m2 {
			if isSameCell(cell1, cell2) || !cell1.isAntenna || !cell2.isAntenna || cell1.frequency != cell2.frequency {
				continue
			}
			// ap := antennaPair{
			// 	a1X: cell1.x,
			// 	a1Y: cell1.y,
			// 	a2X: cell2.x,
			// 	a2Y: cell2.y,
			// }
			// apRev := antennaPair{
			// 	a1X: cell2.x,
			// 	a1Y: cell2.y,
			// 	a2X: cell1.x,
			// 	a2Y: cell1.y,
			// }
			// _, ok := aps[ap]
			// _, okRev := aps[apRev]
			// if !ok && !okRev {
			s := slope{
				dx: cell2.x - cell1.x,
				dy: cell2.y - cell1.y,
			}
			// aps[ap] = s
			_, ok1 := pairs[cell1]
			_, ok2 := pairs[cell2]
			if !ok1 {
				pairs[cell1] = make([]slope, 0)
			}
			if !ok2 {
				pairs[cell2] = make([]slope, 0)
			}
			pairs[cell1] = append(pairs[cell1], s)
			pairs[cell2] = append(pairs[cell1], s)
			// }
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

// func apToSlopeToString(aps map[antennaPair]slope) string {
// 	var export string
// 	for ap, s := range aps {
// 		export += fmt.Sprintf("[%d, %d],[%d,%d]: %d\n", ap.a1X, ap.a1Y, ap.a2X, ap.a2Y, s)
// 	}
// 	return export
// }

func getAllAntinodes(pairs map[cell][]slope, m [][]cell) []cell {
	antinodes := make([]cell, 0)
	for c, slopes := range pairs {
		for _, sl := range slopes {
			an := validAntinodes(c, sl, m)
			antinodes = append(antinodes, an...)
		}
	}
	return antinodes
}

func validAntinodes(c cell, s slope, m [][]cell) []cell {
	antinodes := make([]cell, 0)

	doubleDx := s.dx * 2
	doubleDy := s.dy * 2

	inBounds1 := inBounds(c.x+doubleDx, c.y+doubleDy, m)
	inBounds2 := inBounds(c.x-doubleDx, c.y-doubleDy, m)

	if inBounds1 {
		c1 := cell{
			frequency:  '#',
			x:          c.x + doubleDx,
			y:          c.y + doubleDy,
			isAntinode: true,
		}
		antinodes = append(antinodes, c1)
	}

	if inBounds2 {
		c2 := cell{
			frequency:  '#',
			x:          c.x - doubleDx,
			y:          c.y - doubleDy,
			isAntinode: true,
		}
		antinodes = append(antinodes, c2)
	}

	return antinodes
}

func unique(cs []cell) map[cell]bool {
	mcs := make(map[cell]bool)
	for _, c := range cs {
		mcs[c] = true
	}
	return mcs
}

// inBounds checks if a cell is within the maps's boundaries
func inBounds(x, y int, m [][]cell) bool {
	return x >= 0 && x < len(m) && y >= 0 && y < len(m[0])
}
