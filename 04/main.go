package main

import (
	"bufio"
	"fmt"
	"os"
)

// Advent of Code 2024 - Day 4: Challenge
// https://adventofcode.com/2024/day/4

// Constants for input filename and solutions to find
const (
	filename     string = "input.txt"
	solutionXMAS string = "XMAS"
	solutionMAS  string = "MAS"
	solutionSAM  string = "SAM"
)

// direction represents movement in the grid
type direction struct {
	dx, dy int // Change in x (columns) and y (rows)
}

// Movement directions for part 1 and part 2 of the puzzle
var (
	allDirections = []direction{
		{-1, -1}, // Top-left
		{-1, 0},  // Top
		{-1, 1},  // Top-right
		{0, -1},  // Left
		{0, 1},   // Right
		{1, -1},  // Bottom-left
		{1, 0},   // Bottom
		{1, 1},   // Bottom-right
	}
	diagonalDirections = []direction{
		{-1, -1}, // Top-left
		{-1, 1},  // Top-right
	}
)

// cell represents a single position in the grid with coordinates and value
type cell struct {
	x, y int
	val  rune
}

func main() {
	// Read the puzzle grid from the file and display it
	puzzle := readInput(filename)

	// Part 1: Find all occurrences of the word "XMAS" in the grid
	solutions1 := findSolutions(puzzle, solutionXMAS, allDirections)
	fmt.Println(len(solutions1)) // Print the number of solutions found

	// Part 2: Find occurrences of "MAS" and "SAM" in diagonal directions
	solutions2MAS := findSolutions(puzzle, solutionMAS, diagonalDirections)
	solutions2SAM := findSolutions(puzzle, solutionSAM, diagonalDirections)

	// Combine all solutions for part 2
	solutions2 := append(solutions2MAS, solutions2SAM...)

	// Union solutions by their middle value
	unioned := unionByMiddleVal(solutions2, solutions2)

	// count valid groups - a valid group contains two unique solutions that contain the same middle cell
	fmt.Println(countValid(unioned))
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

// /////////////
// Part 1 & 2 //
// ////////// //

// findSolutions finds all paths in the grid that match the target word
func findSolutions(puzzle [][]cell, solution string, validDirections []direction) [][]cell {
	var solutions [][]cell

	// Iterate over all cells in the grid
	for i, row := range puzzle {
		for j := range row {
			// If the cell matches the first character, start exploring paths
			if puzzle[i][j].val == rune(solution[0]) {
				path := []cell{puzzle[i][j]} // Start a new path
				for _, dir := range validDirections {
					// Explore all paths in the specified directions
					solutions = dfs(puzzle, i, j, solution[1:], path, dir, solutions)
				}
			}
		}
	}
	return solutions
}

// dfs explores paths recursively to find matches for the target word
func dfs(puzzle [][]cell, x, y int, remainingSolution string, path []cell, dir direction, solutions [][]cell) [][]cell {
	// Base case: If no more characters to match, add the path to solutions
	if len(remainingSolution) == 0 {
		solutions = append(solutions, path)
		return solutions
	}

	// Calculate the coordinates of the next cell in the current direction
	nx, ny := x+dir.dx, y+dir.dy

	// Check bounds and character match for the next cell
	if inBounds(nx, ny, puzzle) && puzzle[nx][ny].val == rune(remainingSolution[0]) && !visited(path, puzzle[nx][ny]) {
		newPath := append([]cell{}, path...)      // Create a new path
		newPath = append(newPath, puzzle[nx][ny]) // Add the next cell
		// Continue exploring with updated path and remaining solution
		solutions = dfs(puzzle, nx, ny, remainingSolution[1:], newPath, dir, solutions)
	}
	return solutions
}

// inBounds checks if a cell is within the grid's boundaries
func inBounds(x, y int, puzzle [][]cell) bool {
	return x >= 0 && x < len(puzzle) && y >= 0 && y < len(puzzle[0])
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

// //////////
// Part 2  //
// //////////

// unionByMiddleVal groups paths by the middle cell of their slices
func unionByMiddleVal(setA, setB [][]cell) map[cell][][]cell {
	solutions := make(map[cell][][]cell)

	// Add paths from setA grouped by their middle cell
	for i := range setA {
		middleCell := setA[i][1]
		solutions[middleCell] = append(solutions[middleCell], setA[i])
	}

	// Add paths from setB if they match and have distinct first cells
	for i := range setB {
		middleCell := setB[i][1]
		// check:
		//    1) the middle cell was a part of one of setA's solutions
		//    2) the middle cell is not a part of a solution already included in setA's solutions
		if _, ok := solutions[middleCell]; ok && setA[i][0] != solutions[middleCell][0][0] {
			// append to solution, givng an X-MAS
			solutions[middleCell] = append(solutions[middleCell], setB[i])
		}
	}

	return solutions
}

// countValid counts groups with at least two unique solutions that share the same middle cell
func countValid(us map[cell][][]cell) int {
	var count int
	for _, solution := range us {
		if len(solution) >= 2 {
			count++
		}
	}
	return count
}

// //////////////
// Debug Funcs //
// ////////// //

// printUnionSolutions displays the grouped solutions
func printUnionSolutions(us map[cell][][]cell) {
	for middleCell, solution := range us {
		if len(solution) > 1 {
			fmt.Printf("Middle Cell: (%s)\n%s\n", cellToString(middleCell), solutionToString(solution))
		}
	}
}

// cellToString formats a cell as a readable string
func cellToString(c cell) string {
	return fmt.Sprintf("[%s]: [%d,%d]", string(c.val), c.x, c.y)
}

// solutionToString formats a solution as a readable string
func solutionToString(s [][]cell) string {
	var export string
	for _, row := range s {
		export += "\t"
		for _, cell := range row {
			export += fmt.Sprintf("%s ", cellToString(cell))
		}
		export += "\n"
	}
	return export
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
