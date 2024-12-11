package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	filename    string = "example.txt"
	openVal     rune   = '.'
	antinodeVal rune   = '#'
)

type cell struct {
	val  rune
	x, y int
}

func main() {

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

		currCell := cell{x: i, y: j, val: char}
		row = append(row, currCell)
		i++
	}
	return cells, nil
}
