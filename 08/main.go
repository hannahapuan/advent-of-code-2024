package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// Constants for file input and map symbols
const (
	filename    string = "input.txt"   // Input file containing the grid
	openVal     rune   = '.'           // Open cell value
	antinodeVal rune   = '#'           // Antinode marker
	regex       string = "[a-zA-Z0-9]" // Regex for valid antenna characters
)

// Structure representing each cell in the grid
type cell struct {
	frequency             rune // Antenna frequency or open/antinode value
	x, y                  int  // Coordinates of the cell
	isAntenna, isAntinode bool // Flags for antenna or antinode status
}

// Represents a change in position (dx, dy) between two cells
type change struct {
	dx, dy int
}

func main() {
	// Read the grid from the input file
	m, err := readInput(filename)
	if err != nil {
		os.Exit(1) // Exit if the file cannot be read
	}

	// Flatten the 2D grid into a single slice
	mf := flatten2dSlice(m)

	// Calculate antenna pairs based on frequency and position
	pairs := calcAntennaPairs(mf, mf)

	// Print the initial grid with coordinates
	fmt.Println(mapToString(m))

	// Calculate and display antinodes without resonant harmonics
	an := getAllAntinodes(pairs, m, false)
	uan := unique(an)
	mapWithAntinodes := updateMapWithAntinodes(m, uan)
	fmt.Println(mapToString(mapWithAntinodes))
	fmt.Println(len(uan))

	// Calculate and display antinodes with resonant harmonics
	anrh := getAllAntinodes(pairs, m, true)
	uanrh := unique(anrh)
	mapWithAntinodesrh := updateMapWithAntinodes(m, uanrh)
	fmt.Println(mapToString(mapWithAntinodesrh))
	fmt.Println(len(uanrh))
}

// Reads the input file and converts it into a 2D grid of cells
func readInput(fname string) ([][]cell, error) {
	cells := make([][]cell, 0) // 2D grid

	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file:  %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var row []cell
	var i, j int

	// Compile the regex for antenna characters
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex string: %s", err)
	}

	for {
		// Read one character at a time
		char, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" { // End of file
				if len(row) > 0 {
					cells = append(cells, row) // Append the last row
				}
				break
			}
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		if char == '\n' { // Newline indicates the end of a row
			cells = append(cells, row)
			row = make([]cell, 0)
			j++   // Increment y-coordinate
			i = 0 // Reset x-coordinate
			continue
		}

		// Create a cell for the current character
		currCell := cell{x: i, y: j, frequency: char}

		// Mark as an antenna if it matches the regex
		if r.MatchString(string(char)) {
			currCell.isAntenna = true
		}

		row = append(row, currCell)
		i++
	}
	return cells, nil
}

// Converts the 2D grid into a printable string with axis labels
func mapToString(s [][]cell) string {
	var export string

	// Add X-axis labels
	export += "   " // Padding for Y-axis labels
	for x := 0; x < len(s[0]); x++ {
		export += fmt.Sprintf("%2d ", x)
	}
	export += "\n"

	// Add each row with Y-axis labels
	for y, row := range s {
		export += fmt.Sprintf("%2d ", y) // Y-axis label
		for _, cell := range row {
			export += fmt.Sprintf("%s  ", string(cell.frequency))
		}
		export += "\n"
	}
	return export
}

// Finds all pairs of antennas with the same frequency
func calcAntennaPairs(m1 []cell, m2 []cell) map[cell][]change {
	pairs := make(map[cell][]change)

	for _, cell1 := range m1 {
		for _, cell2 := range m2 {
			// Skip invalid or duplicate cells
			if isSameCell(cell1, cell2) || !cell1.isAntenna || !cell2.isAntenna || cell1.frequency != cell2.frequency {
				continue
			}
			// Calculate the relative position (dx, dy)
			s := change{
				dx: cell2.x - cell1.x,
				dy: cell2.y - cell1.y,
			}
			// Add the pair
			pairs[cell1] = append(pairs[cell1], s)
		}
	}

	return pairs
}

// Compares two cells for equality
func isSameCell(a, b cell) bool {
	return a.x == b.x && a.y == b.y
}

// Flattens a 2D grid into a 1D slice
func flatten2dSlice(s [][]cell) []cell {
	fs := make([]cell, 0)
	for _, row := range s {
		fs = append(fs, row...)
	}
	return fs
}

// Retrieves all valid antinodes
// note that the last argument is a flag determining whether resonant harmonics should be calculated
func getAllAntinodes(pairs map[cell][]change, m [][]cell, withResonantHarmonics bool) []cell {
	antinodes := make([]cell, 0)
	for c, changes := range pairs {
		for _, sl := range changes {
			an := make([]cell, 0)
			if withResonantHarmonics {
				an = validAntinodesWithResonantHarmonics(c, sl, m)
			} else {
				an = validAntinodes(c, sl, m)
			}
			antinodes = append(antinodes, an...)
		}
	}
	return antinodes
}

// Calculates valid antinodes without resonant harmonics
func validAntinodes(c cell, s change, m [][]cell) []cell {
	antinodes := make([]cell, 0)

	if s.dx == 0 { // Vertical case
		newY := c.y + (s.dy * 2)
		if inBounds(c.x, newY, m) {
			antinodes = append(antinodes, cell{frequency: '#', x: c.x, y: newY, isAntinode: true})
		}
		return antinodes
	}

	// Normal case
	newX := c.x + (s.dx * 2)
	newY := c.y + (s.dy * 2)
	if inBounds(newX, newY, m) {
		antinodes = append(antinodes, cell{frequency: '#', x: newX, y: newY, isAntinode: true})
	}

	return antinodes
}

// Calculates valid antinodes with resonant harmonics
func validAntinodesWithResonantHarmonics(c cell, s change, m [][]cell) []cell {
	antinodes := make([]cell, 0)

	if s.dx == 0 { // Vertical case
		newY := c.y + s.dy
		for inBounds(c.x, newY, m) {
			antinodes = append(antinodes, cell{frequency: '#', x: c.x, y: newY, isAntinode: true})
			newY += s.dy
		}
		return antinodes
	}

	// Normal case
	newX := c.x + s.dx
	newY := c.y + s.dy
	for inBounds(newX, newY, m) {
		antinodes = append(antinodes, cell{frequency: '#', x: newX, y: newY, isAntinode: true})
		newX += s.dx
		newY += s.dy
	}

	return antinodes
}

// Filters out duplicate antinodes
func unique(cs []cell) map[[2]int]bool {
	mcs := make(map[[2]int]bool)
	for _, c := range cs {
		mcs[[2]int{c.x, c.y}] = true
	}
	return mcs
}

// Checks if a coordinate is within bounds
func inBounds(x, y int, m [][]cell) bool {
	return x >= 0 && x < len(m[0]) && y >= 0 && y < len(m)
}

// Updates the grid to include antinodes
func updateMapWithAntinodes(m [][]cell, antiNodes map[[2]int]bool) [][]cell {
	for antinode := range antiNodes {
		x, y := antinode[0], antinode[1]
		if inBounds(x, y, m) {
			m[y][x].isAntinode = true
			m[y][x].frequency = antinodeVal
		}
	}
	return m
}
