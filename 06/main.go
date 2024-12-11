package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// https://adventofcode.com/2024/day/6

const (
	fileName    string = "input.txt" // Name of the input file containing the grid
	visitedRune rune   = 'X'         // Rune representing a visited cell
	openRune    rune   = '.'         // Rune representing an open (unvisited) cell
	blockedRune rune   = '#'         // Rune representing a blocked cell
)

var (
	// Maps arrow characters to their respective directions
	arrowToDir map[rune]string = map[rune]string{
		'v': "down",
		'^': "up",
		'>': "right",
		'<': "left",
	}
	// Array of possible directions for traversal
	directions = []string{"up", "right", "down", "left"}
	// Maps directions to their respective movement deltas
	moves = map[string]move{
		"up":    {0, -1},
		"right": {1, 0},
		"down":  {0, 1},
		"left":  {-1, 0},
	}
)

// Represents a movement with x and y deltas
type move struct {
	dx, dy int
}

// Represents a single position in the grid with coordinates and value
type cell struct {
	x, y int
	val  rune
}

// Represents the guard's state: current position, path traversed, and direction
type guard struct {
	currPos   cell
	path      []cell
	direction string
}

func main() {
	// Read the grid and initialize the guard's state
	cells, guard, err := readInput(fileName)
	if err != nil {
		os.Exit(1)
	}
	count := 0

	// Simulate the guard's traversal until no moves are available
	for {
		guard, cells, err = step(guard, cells)
		if err != nil {
			// Print final results when traversal ends
			fmt.Printf("Finished traversal. Total steps: %d\n", count)
			fmt.Printf("\t Total distinct steps: %d\n", distinctPositions(cells))
			break
		}
		// // Print the grid and guard's position after each step
		// fmt.Println(puzzleToString(cells))
		// fmt.Println(guardPosToString(guard))
		// fmt.Println("\n---------------------\n")
		count++
	}
}

// Reads the input file and initializes the grid and guard's starting state
func readInput(fname string) ([][]cell, guard, error) {
	cells := make([][]cell, 0) // 2D array representing the grid
	guardPath := make([]cell, 0)
	gu := guard{
		path: guardPath,
	}

	file, err := os.Open(fname)
	if err != nil {
		return nil, guard{}, fmt.Errorf("error opening file: %w", err)
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
			return nil, guard{}, fmt.Errorf("error reading file :%w", err)
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

		// Identify the guard's starting position and direction
		dir, ok := arrowToDir[char]
		if ok {
			// mark starting position as visted and add it to the guard path
			char = visitedRune
			currCell.val = visitedRune
			gu.path = append(gu.path, currCell)
			gu.currPos = currCell
			gu.direction = dir
		}
		row = append(row, currCell)
		i++
	}
	return cells, gu, nil
}

// Simulates a single step of the guard's movement
func step(g guard, cells [][]cell) (guard, [][]cell, error) {
	directionsTried := 0
	for directionsTried < len(directions) {
		move := moves[g.direction]
		newX := g.currPos.x + move.dx
		newY := g.currPos.y + move.dy

		// Stop traversal if out of bounds
		if !inBounds(newX, newY, cells) {
			return g, cells, errors.New("no more valid moves")
		}
		// Check if the move is valid
		if cells[newY][newX].val != blockedRune {
			// Mark the new cell as visited
			cells[newY][newX].val = visitedRune
			// Update the guard's position and path
			g.currPos = cell{x: newX, y: newY, val: visitedRune}
			g.path = append(g.path, g.currPos)
			return g, cells, nil
		}

		// Rotate to the next direction if the current move is invalid
		g.direction = turnRight(g.direction)
		directionsTried++
	}

	// Return an error if no valid moves are available
	return g, cells, errors.New("no valid moves available")
}

// Rotates the guard's direction 90 degrees clockwise
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

// Checks if the given coordinates are within the grid bounds
func inBounds(x, y int, cells [][]cell) bool {
	return x >= 0 && x < len(cells) && y >= 0 && y < len(cells[x])
}

// Counts the number of distinct visited cells
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

// Formats the guard's current position as a string
func guardPosToString(g guard) string {
	return fmt.Sprintf("guard pos: (%d,%d)\n", g.currPos.x, g.currPos.y)
}

// Converts the grid into a printable string representation
func puzzleToString(s [][]cell) string {
	var export string
	for _, row := range s {
		for _, cell := range row {
			export += fmt.Sprintf("%s ", string(cell.val))
		}
		export += "\n"
	}
	return export
}
