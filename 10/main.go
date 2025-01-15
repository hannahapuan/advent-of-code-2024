package main

import (
	"bufio"
	"fmt"
	"os"
)

// Advent of Code 2024 - Day 10: Challenge
// Link: https://adventofcode.com/2024/day/10

const (
	trailhead rune   = '0'
	fileName  string = "example.txt"
	solution  string = "012345678909876543210"
)

// direction represents movement in the grid
type direction struct {
	dx, dy int // Change in x (columns) and y (rows)
}

var (
	allDirections = []direction{
		{-1, 0}, // Top
		{0, -1}, // Left
		{0, 1},  // Right
		{1, 0},  // Bottom
	}
)

// cell represents a single position in the grid with coordinates and value
type cell struct {
	x, y int
	val  rune
}

func main() {
	tmap := readInput(fileName)
	puzzleToString(tmap)

	solutions1 := findSolutions(tmap, solution, allDirections)
	fmt.Println(len(solutions1)) // Print the number of solutions found
}

// readInput reads the grid from the file and converts it to a 2D slice of cells
func readInput(fname string) [][]cell {
	cells := make([][]cell, 0)

	file, err := os.Open(fname)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
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
			return nil
		}

		if char == '\n' { // Handle newlines as row separators
			cells = append(cells, row)
			row = make([]cell, 0)
			j++ // Move to the next row
			i = 0
			continue
		}

		// Add the character as a cell to the current row
		row = append(row, cell{x: i, y: j, val: char})
		i++
	}
	return cells
}

// findSolutions finds all paths in the grid that match the target word
func findSolutions(tmap [][]cell, solution string, validDirections []direction) [][]cell {
	var solutions [][]cell

	// Iterate over all cells in the grid
	for i, row := range tmap {
		for j := range row {
			// If the cell matches the first character, start exploring paths
			fmt.Printf("exploring %s==%s\n", string(tmap[i][j].val), string(solution[0]))
			if tmap[i][j].val == rune(solution[0]) {
				fmt.Println("\tfound")
				path := []cell{tmap[i][j]} // Start a new path
				for _, dir := range validDirections {
					// Explore all paths in the specified directions
					solutions = dfs(tmap, i, j, solution[1:], path, dir, solutions)
				}
			}
		}
	}
	return solutions
}

// inBounds checks if a cell is within the grid's boundaries
func inBounds(x, y int, tmap [][]cell) bool {
	return x >= 0 && x < len(tmap) && y >= 0 && y < len(tmap[0])
}

// move returns an update map
func dfs(tmap [][]cell, x, y int, remainingSolution string, path []cell, dir direction, solutions [][]cell) [][]cell {
	// Base case: If no more characters to match, add the path to solutions
	if len(remainingSolution) == 0 {
		solutions = append(solutions, path)
		return solutions
	}

	// Calculate the coordinates of the next cell in the current direction
	nx, ny := x+dir.dx, y+dir.dy

	// Check bounds and character match for the next cell
	fmt.Println("remainingSolution: %s", string(remainingSolution[0]))
	if inBounds(nx, ny, tmap) && tmap[nx][ny].val == rune(remainingSolution[0]) && !visited(path, tmap[nx][ny]) {
		newPath := append([]cell{}, path...)    // Create a new path
		newPath = append(newPath, tmap[nx][ny]) // Add the next cell
		// Continue exploring with updated path and remaining solution
		solutions = dfs(tmap, nx, ny, remainingSolution[1:], newPath, dir, solutions)
	}
	return solutions
}

// visited checks if a cell is already part of the current path
func visited(path []cell, c cell) bool {
	for _, p := range path {
		if p.x == c.x && p.y == c.y {
			return true
		}
	}
	return false
}

// puzzleToString prints the grid to the console
func puzzleToString(s [][]cell) {
	for _, row := range s {
		for _, cell := range row {
			fmt.Print(string(cell.val))
		}
		fmt.Println()
	}
}
