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
)

func main() {
	tmap, err := readInput(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(tmap)
}

func readInput(fname string) ([][]int, error) {
	export := make([][]int, 0)

	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var row []int
	var i, j int

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				if len(row) > 0 {
					export = append(export, row) // Append the last row
				}
				break
			}
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		if char == '\n' { // Handle newlines as row separators
			export = append(export, row)
			row = make([]int, 0)
			j++ // Move to the next row
			i = 0
			continue
		}

		// Add the character as a cell to the current row
		row = append(row, int(char-'0'))
		i++
	}

	return export, nil
}
