package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	filename string = "input.txt" // Name of the input file
	solution string = "XMAS"      // Target solution to find
)

// direction represents a movement in the grid
type direction struct {
	dx, dy int // Changes in x and y coordinates for the direction
}

// Possible movement directions for adjacency search
var (
	directions = []direction{
		{-1, -1}, // Top-left
		{-1, 0},  // Top
		{-1, 1},  // Top-right
		{0, -1},  // Left
		{0, 1},   // Right
		{1, -1},  // Bottom-left
		{1, 0},   // Bottom
		{1, 1},   // Bottom-right
	}
)

// Represents a cell in the grid with its coordinates and rune value
type cell struct {
	x, y int
	val  rune
}

func main() {
	// Read the puzzle grid from the input file
	puzzle := readInput(filename)

	// Find all paths in the grid that match the target solution
	solutions := findSolutions(puzzle, solution)

	// Print the number of solutions found
	fmt.Println(len(solutions))
}

// readInput reads the puzzle grid from a file and returns a 2D slice of cells
func readInput(fname string) [][]cell {
	cells := make([][]cell, 0) // Initialize the 2D slice

	file, err := os.Open(fname)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var i, j int
	row := make([]cell, 0)

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				if len(row) > 0 { // Append the last row if not empty
					cells = append(cells, row)
				}
				break
			}
			fmt.Println("Error reading file:", err)
			return nil
		}

		if char == '\n' { // Handle newline
			cells = append(cells, row)
			row = make([]cell, 0)
			j++
			i = 0
			continue
		}

		row = append(row, cell{
			x:   i,
			y:   j,
			val: char,
		})
		i++
	}
	return cells
}

// findSolutions searches the grid for all paths that match the target solution
func findSolutions(puzzle [][]cell, solution string) [][]cell {
	var solutions [][]cell // List to store all matching paths

	// Iterate over all cells in the grid
	for i, row := range puzzle {
		for j := range row {
			// Check if the current cell matches the first character of the solution
			if puzzle[i][j].val == rune(solution[0]) {
				var path []cell                   // Start a new path
				path = append(path, puzzle[i][j]) // Add the current cell to the path

				// Explore all possible directions from this cell
				for _, dir := range directions {
					solutions = dfs(puzzle, i, j, solution[1:], path, dir, solutions)
				}
			}
		}
	}

	return solutions
}

// dfs performs a recursive depth-first search to find paths that match the solution
func dfs(puzzle [][]cell, x, y int, remainingSolution string, path []cell, dir direction, solutions [][]cell) [][]cell {
	// base case: If there are no more characters to match, add the path to solutions
	if len(remainingSolution) == 0 {
		solutions = append(solutions, path)
		return solutions
	}

	// next cell's coordinates
	nx, ny := x+dir.dx, y+dir.dy

	if inBounds(nx, ny, puzzle) {
		// check if the next cell matches the next character in the solution
		// and if it hasn't already been visited in the current path
		if puzzle[nx][ny].val == rune(remainingSolution[0]) && !visited(path, puzzle[nx][ny]) {
			// Create a new path including the next cell
			newPath := append([]cell{}, path...)
			newPath = append(newPath, puzzle[nx][ny])

			// recursive case: continue the search with the updated path and remaining solution
			solutions = dfs(puzzle, nx, ny, remainingSolution[1:], newPath, dir, solutions)
		}
	}

	return solutions
}

// inBounds checks if the next cell is within bounds
func inBounds(x, y int, puzzle [][]cell) bool {
	return x >= 0 && x < len(puzzle) && y >= 0 && y < len(puzzle[0])
}

// visited checks if a cell is already in the current path
func visited(path []cell, c cell) bool {
	for _, p := range path {
		if p.x == c.x && p.y == c.y {
			return true // Cell is already in the path
		}
	}
	return false
}
