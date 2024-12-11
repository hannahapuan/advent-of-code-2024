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
}

type slope struct {
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

// Converts the grid into a printable string representation
func mapToString(s [][]cell) string {
	export := "  "
	for i := 0; i < len(s); i++ {
		export += fmt.Sprintf("%d ", i)
	}
	export += "\n"
	var y int
	for _, row := range s {
		export += fmt.Sprintf("%d ", y)
		y++
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
			s := slope{
				dx: cell2.x - cell1.x,
				dy: cell2.y - cell1.y,
			}
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
	rateOfChange := s.dy / s.dx * 2
	newX := c.x + rateOfChange
	newY := c.y + rateOfChange
	newNegX := c.x - rateOfChange
	newNegY := c.y - rateOfChange
	inBounds1 := inBounds(newX, newY, m)
	inBounds2 := inBounds(newNegX, newNegY, m)

	if inBounds1 {
		c1 := cell{
			frequency:  '#',
			x:          newX,
			y:          newY,
			isAntinode: true,
		}
		antinodes = append(antinodes, c1)
	}

	if inBounds2 {
		c2 := cell{
			frequency:  '#',
			x:          newNegX,
			y:          newNegY,
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

func updateMapWithAntinodes(m [][]cell, antiNodes map[cell]bool) [][]cell {
	for antinode := range antiNodes {
		m[antinode.x][antinode.y].isAntinode = true
		m[antinode.x][antinode.y].frequency = antinodeVal
	}
	return m
}
