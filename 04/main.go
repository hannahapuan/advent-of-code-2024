package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: either something wrong with the calculation of the adj cells or the tracking through them in the sols
const (
	filename    string = "example.txt" // Name of the input file
	solution    string = "XMAS"        // Target solution to find
	newlineRune        = '\n'          // Newline character
)

var (
	rSol       []rune = []rune(solution) // Solution as a slice of runes
	directions        = []direction{     // Directions for adjacency search
		{-1, -1}, // Top-left
		{-1, 1},  // Top-right
		{1, -1},  // Bottom-left
		{1, 1},   // Bottom-right
		{-1, 0},  // Top
		{0, -1},  // Left
		{1, 0},   // Bottom
		{0, 1},   // Right
	}
)

// cell represents a position in the puzzle with x, y coordinates and a rune value
type cell struct {
	x   int
	y   int
	val rune
}

// direction represents a movement in the grid
type direction struct {
	dx, dy int
}

func main() {
	// Read the puzzle grid from the input file
	puzzle := readInput(filename)

	// Generate a map of cells and their adjacent cells by rune and direction
	cellsToAdjCells := genCellToAdjCellsMap(puzzle)

	// Print the adjacency map for debugging purposes
	printCellToAdjCellsMap(cellsToAdjCells)

	// Find and print all solutions
	fmt.Println()
	fmt.Println(findSols(cellsToAdjCells))
	fmt.Println(len(findSols(cellsToAdjCells)))
}

// readInput reads the puzzle grid from a file and returns a 2D slice of cells
func readInput(fname string) [][]cell {
	cells := make([][]cell, 0) // Initialize the 2D slice

	file, err := os.Open(fname) // Open the file
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	reader := bufio.NewReader(file) // Create a buffered reader

	var i, j int           // Row and column indices
	row := make([]cell, 0) // Current row of cells

	for {
		char, _, err := reader.ReadRune() // Read the next rune
		if err != nil {
			if err.Error() == "EOF" { // End of file
				if len(row) > 0 { // Append the last row if not empty
					cells = append(cells, row)
				}
				break
			}
			fmt.Println("Error reading file:", err)
			return nil
		}

		if char == '\n' { // Handle newline
			cells = append(cells, row) // Append completed row to cells
			row = make([]cell, 0)      // Reset the row
			j++                        // Increment row index
			i = 0                      // Reset column index
			continue
		}

		// Add the character to the current row as a cell
		row = append(row, cell{
			x:   i,
			y:   j,
			val: char,
		})
		i++ // Increment column index
	}
	return cells
}

// genCellToAdjCellsMap creates a map of runes to their adjacent cells by direction
func genCellToAdjCellsMap(puzzle [][]cell) map[rune]map[direction][]cell {
	cellToAdjCells := make(map[rune]map[direction][]cell)
	for i, row := range puzzle {
		for j := range row {
			ru := puzzle[i][j].val
			if _, ok := cellToAdjCells[ru]; !ok {
				cellToAdjCells[ru] = make(map[direction][]cell)
			}
			newAdjCells := populateAdjCells(i, j, puzzle)
			for dir, cells := range newAdjCells {
				cellToAdjCells[ru][dir] = append(cellToAdjCells[ru][dir], cells...)
			}
		}
	}
	return cellToAdjCells
}

// populateAdjCells finds all valid adjacent cells for a given cell in the puzzle
func populateAdjCells(x, y int, puzzle [][]cell) map[direction][]cell {
	adjCells := make(map[direction][]cell)
	height := len(puzzle)      // Number of rows
	width := len(puzzle[0])    // Number of columns
	curVal := puzzle[x][y].val // Value of the current cell

	// Iterate over all possible directions
	for _, dir := range directions {
		if inRange(x+dir.dx, y+dir.dy, width, height) {
			for i, r := range rSol[:len(rSol)-1] { // Match against solution runes
				if curVal == r {
					adjCells = appendIfVal(puzzle[x+dir.dx][y+dir.dy], adjCells, dir, rSol[i+1])
				}
			}
		}
	}
	return adjCells
}

// inRange checks if the given x, y coordinates are within the grid boundaries
func inRange(x, y, width, height int) bool {
	return x >= 0 && x < height && y >= 0 && y < width
}

// appendIfVal adds a cell to the map if it matches the desired next rune
func appendIfVal(toInsert cell, cells map[direction][]cell, dir direction, desiredNextVal rune) map[direction][]cell {
	if toInsert.val != desiredNextVal { // Skip if value doesn't match the desired rune
		return cells
	}
	if _, ok := cells[dir]; !ok { // Initialize slice if it doesn't exist
		cells[dir] = make([]cell, 0)
	}
	cells[dir] = append(cells[dir], toInsert) // Append the cell
	return cells
}

// findSols finds all valid paths matching the solution in the puzzle
func findSols(ctac map[rune]map[direction][]cell) [][]cell {
	solutions := make([][]cell, 0)
	path := make([]cell, 0)
	desiredRuneIndex := 0
	for _, dirSearching := range directions {
		for desiredRuneIndex < len(rSol)-1 {
			// fmt.Println("here")
			// fmt.Printf("\tdesiredRuneIndex: %d\n", desiredRuneIndex)
			// fmt.Printf("\tdesiredRune: %s\n", string(rSol[desiredRuneIndex]))
			// fmt.Printf("\tdirSearching: %v\n", dirSearching)
			desiredCells := ctac[rSol[desiredRuneIndex]][dirSearching]
			ctac[rSol[desiredRuneIndex]][dirSearching] = deleteCell(desiredCells, 0)
			// fmt.Printf("\tdesiredCells: %v\n", toString(desiredCells))

			if len(desiredCells) == 0 || desiredRuneIndex == len(rSol)-1 {
				break
			}
			for _, cell := range desiredCells {
				if desiredRuneIndex == len(rSol)-1 {
					// no solution found
					break
				}
				if cell.val == rSol[desiredRuneIndex+1] {
					path = append(path, cell)
					desiredRuneIndex++
				}
			}
		}
		desiredRuneIndex = 0
		solutions = append(solutions, path)
		path = make([]cell, 0)
	}
	return solutions
}

func deleteCell(cs []cell, index int) []cell {
	return append(cs[:index], cs[index+1:]...)
}

// // searchPath recursively finds paths that match the solution
// func searchPath(ctac map[rune]map[direction][]cell, current cell, index int, path []cell, solutions *[][]cell) bool {
// 	if index == len(rSol) { // If the full solution is found
// 		*solutions = append(*solutions, append([]cell(nil), path...))
// 		return true
// 	}

// 	nextRune := rSol[index]
// 	for dir, adjCells := range ctac[current.val] { // Check adjacent cells
// 		for _, adj := range adjCells {
// 			if inRange(adj.x, adj.y, len(ctac), len(ctac[adj.val])) &&
// 				adj.val == nextRune && dir.dx == adj.x-current.x && dir.dy == adj.y-current.y {
// 				path = append(path, adj) // Add cell to path
// 				if searchPath(ctac, adj, index+1, path, solutions) {
// 					return true
// 				}
// 				path = path[:len(path)-1] // Backtrack
// 			}
// 		}
// 	}

// 	return false
// }

// printCellToAdjCellsMap prints the adjacency map for debugging
func printCellToAdjCellsMap(ctac map[rune]map[direction][]cell) {
	for r, v := range ctac {
		fmt.Println(string(r))
		for dir, cells := range v {
			fmt.Printf("\t(%d, %d): %s\n", dir.dx, dir.dy, toString(cells))
		}
	}
}

// toString converts a slice of cells to a slice of their rune values as strings
func toString(c []cell) []string {
	cs := make([]string, 0)
	for _, v := range c {
		cs = append(cs, string(v.val))
	}
	return cs
}

func toCellsString(c []cell) []string {
	cs := make([]string, 0)
	for _, v := range c {
		cs = append(cs, toCellString(v))
	}
	return cs
}

// toCellString formats a cell as a string with its coordinates and value
func toCellString(c cell) string {
	return fmt.Sprintf("(%d, %d): %s", c.x, c.y, string(c.val))
}
