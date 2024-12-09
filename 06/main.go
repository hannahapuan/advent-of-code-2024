package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// https://adventofcode.com/2024/day/6

const (
	fileName    string = "example.txt"
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
	dirToArrow map[string]rune = map[string]rune{
		"down":  'v',
		"up":    '^',
		"right": '>',
		"left":  '<',
	}
	directions = []string{"up", "right", "down", "left"}
	cellVal    = []string{"X", ".", "#"}
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
	for {
		var err error
		guard, cells, err = step(guard, cells)
		if err != nil {
			fmt.Println("Finished traversal. Total steps:", len(guard.path))
			break
		}
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
		fmt.Printf("Trying direction: %s at (%d, %d)\n", g.direction, g.currPos.x, g.currPos.y)
		fmt.Printf("Cell status: %v\n", cells)
		move := moves[g.direction]
		newX := g.currPos.x + move.dx
		newY := g.currPos.y + move.dy

		// Check if the move is valid
		if inBounds(newX, newY, cells) && cells[newY][newX].val == openRune {
			// Mark the new cell as visited
			cells[newY][newX].val = visitedRune

			// Update guard position
			g.currPos = cell{x: newX, y: newY, val: visitedRune}
			g.path = append(g.path, g.currPos)
			return g, cells, nil
		}

		// Rotate to try the next direction
		g.direction = rotate90(g.direction)
		directionsTried++
	}

	// If no valid moves are found
	return g, cells, errors.New("no valid moves available")
}

// func step(g guard, cells [][]cell) (guard, [][]cell, error) {
// 	fmt.Println()
// 	fmt.Println("NEW STEP")

// 	d := g.direction
// 	fmt.Printf("Initial direction: %s\n", d)

// 	newX := g.currPos.x + moves[d].dx
// 	newY := g.currPos.y + moves[d].dy

// 	visitedDirectionsCount := 0

// 	// Try all directions until a valid move is found
// 	for visitedDirectionsCount < len(directions)-1 {
// 		if inBounds(newX, newY, cells) && !visited(newX, newY, g) && cells[newX][newY].val != blockedRune {
// 			// Found a valid move
// 			break
// 		}
// 		// Rotate to the next direction
// 		d = rotate90(d)
// 		fmt.Printf("Rotating to: %s\n", d)
// 		newX = g.currPos.x + moves[d].dx
// 		newY = g.currPos.y + moves[d].dy
// 		visitedDirectionsCount++
// 	}

// 	// If all directions are blocked
// 	if visitedDirectionsCount == len(directions) || !inBounds(newX, newY, cells) {
// 		return guard{}, nil, errors.New("no valid moves available")
// 	}

// 	// Update guard position and path
// 	g.currPos = cell{
// 		x:   newX,
// 		y:   newY,
// 		val: g.currPos.val,
// 	}
// 	g.path = append(g.path, g.currPos)

// 	// Update the grid
// 	cells[newX][newY].val = visitedRune

// 	fmt.Printf("New position: (%d,%d) | Direction: %s\n", newX, newY, d)

// 	// Update guard's direction
// 	g.direction = d

// 	return g, cells, nil
// }

func rotate90(dir string) string {
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

func isValidMove(dir string, x, y int, cells [][]cell) bool {
	mm := moves[dir]
	newX := x + mm.dx
	newY := y + mm.dy
	return inBounds(newX, newY, cells) && cells[newX][newY].val != blockedRune
}

func inBounds(x, y int, cells [][]cell) bool {
	return x >= 0 && x < len(cells) && y >= 0 && y < len(cells[x])
}

func visited(x, y int, g guard) bool {
	for _, cell := range g.path {
		if cell.x == x && cell.y == y {
			return true
		}
	}
	return false
}
