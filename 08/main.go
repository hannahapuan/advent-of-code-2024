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
	val       rune
	x, y      int
	isAntenna bool
	// matchingAntennas []cell
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

	fmt.Println(mapToString(m))
	fmt.Println(atl)

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
		currCell := cell{x: i, y: j, val: char}

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
			export += fmt.Sprintf("%s ", string(cell.val))
		}
		export += "\n"
	}
	return export
}
