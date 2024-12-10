package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// https://adventofcode.com/2024/day/6

const (
	fileName    string = "input.txt"
	visitedRune rune   = 'X'
	openRune    rune   = '.'
	blockedRune rune   = '#'
)

// FIXME: remove these later theyre just for ref for now
var (
	arrowToDir map[rune]string = map[rune]string{
		'v': "down",
		'^': "up",
		'>': "right",
		'<': "left",
	}
	directions = []string{"up", "right", "down", "left"}
	moves      = map[string]move{
		"up":    {0, -1},
		"right": {1, 0},
		"down":  {0, 1},
		"left":  {-1, 0},
	}
)

type move struct {
	dx, dy int
}

// cell represents a single position in the grid with coordinates and value
type cell struct {
	x, y int
	val  rune
}

type guard struct {
	currPos   cell
	path      []cell
	direction string
}

func main() {
	cells, guard := readInput(fileName)
	count := 0
	for {
		var err error
		guard, cells, err = step(guard, cells)
		if err != nil {
			fmt.Printf("Finished traversal. Total steps: %d\n", count)
			fmt.Printf("\t Total distinct steps: %d\n", distinctPositions(cells))
			break
		}
		// puzzleToString(cells)
		// printGuardMove(guard)
		// fmt.Println("\n---------------------\n")
		count++
	}
}

// output:
// 2d array of map of cells
// guard path (contains only one cell value of the init value of the guard)

func readInput(fname string) ([][]cell, guard) {
	cells := make([][]cell, 0)
	guardPath := make([]cell, 0)
	// TODO: update this to handle any direction, right now it is up as the example starts at up
	gu := guard{
		path: guardPath,
	}

	file, err := os.Open(fname)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, guard{}
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
					cells = append(cells, row) // Append the last row
				}
				break
			}
			fmt.Println("Error reading file:", err)
			return nil, guard{}
		}

		if char == '\n' { // Handle newlines as row separators
			cells = append(cells, row)
			row = make([]cell, 0)
			j++ // Move to the next row
			i = 0
			continue
		}

		currCell := cell{x: i, y: j, val: char}
		// Add the character as a cell to the current row
		// FIXME: this startrune will change
		dir, ok := arrowToDir[char]
		// found where the guard starts
		if ok {
			char = visitedRune
			currCell.val = visitedRune
			gu.path = append(gu.path, currCell)
			gu.currPos = currCell
			gu.direction = dir
		}
		row = append(row, currCell)
		i++
	}
	return cells, gu
}

func step(g guard, cells [][]cell) (guard, [][]cell, error) {
	directionsTried := 0
	for directionsTried < len(directions) {
		move := moves[g.direction]
		newX := g.currPos.x + move.dx
		newY := g.currPos.y + move.dy

		if !inBounds(newX, newY, cells) {
			return g, cells, errors.New("done!")
		}
		// Check if the move is valid
		if cells[newY][newX].val != blockedRune {
			// Mark the new cell as visited
			cells[newY][newX].val = visitedRune

			// Update guard position
			g.currPos = cell{x: newX, y: newY, val: visitedRune}
			g.path = append(g.path, g.currPos)

			return g, cells, nil
		}

		// Rotate to try the next direction
		g.direction = turnRight(g.direction)
		directionsTried++
	}

	// If no valid moves are found
	return g, cells, errors.New("no valid moves available")
}

func turnRight(dir string) string {
	switch dir {
	case "up":
		return "right"
	case "right":
		return "down"
	case "down":
		return "left"
	case "left":
		return "up"
	}
	return ""
}

// puzzleToString prints the grid to the console
func puzzleToString(s [][]cell) {
	for _, row := range s {
		for _, cell := range row {
			fmt.Printf("%s ", string(cell.val))
		}
		fmt.Println()
	}
}

func inBounds(x, y int, cells [][]cell) bool {
	return x >= 0 && x < len(cells) && y >= 0 && y < len(cells[x])
}

func distinctPositions(s [][]cell) int {
	var count int
	for _, row := range s {
		for _, cell := range row {
			if cell.val == visitedRune {
				count++
			}
		}
	}
	return count
}

func printGuardMove(g guard) {
	fmt.Printf("guard pos: (%d,%d)\n", g.currPos.x, g.currPos.y)
}
