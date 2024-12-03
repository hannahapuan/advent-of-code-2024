package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// https://adventofcode.com/2024/day/2

const (
	filename string = "input.txt"
	delim    string = " "
)

func main() {
	reports, err := readInput(filename)
	if err != nil {
		fmt.Println("error reading input file: %w", err)
		os.Exit(1)
	}

	safeCount := countSafe(reports)
	fmt.Println(safeCount)
}

// Reads the input file and parses each report
func readInput(fname string) ([][]int, error) {
	var reports [][]int

	// Open the input file
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file [%s]: %w", fname, err)
	}
	defer f.Close() // Ensure the file is closed after the function completes

	scanner := bufio.NewScanner(f)
	var l []byte
	for scanner.Scan() {
		l = scanner.Bytes()                  // Read the line as a byte slice
		v := strings.Split(string(l), delim) // Split the line using the delimiter
		// Parse the values as ints
		var report []int
		for i := range v {
			l, err := strconv.Atoi(v[i])
			if err != nil {
				return nil, fmt.Errorf("error parsing file [%s]: expected int, found %s", fname, v[i])
			}
			report = append(report, l)
		}

		// Append report to the list of reports
		reports = append(reports, report)

	}
	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file [%s]: %w", fname, err)
	}

	return reports, nil
}

func countSafe(reports [][]int) int {
	var safeCount int

	for _, report := range reports {
		if isSafe(report) {
			safeCount++
		}
	}

	return safeCount
}

func isSafe(report []int) bool {
	var isDecreasing bool
	diff := report[0] - report[1]

	if diff < 0 {
		isDecreasing = true
	}

	for i := range report {
		if i == 0 {
			continue
		}
		diff = report[i-1] - report[i]
		absDiff := absDiffInt(report[i], report[i-1])

		if (absDiff < 1) || (absDiff > 3) {
			return false
		}

		if diff > 0 && isDecreasing {
			return false
		} else if diff < 0 && !isDecreasing {
			return false
		}
	}
	return true
}

// Helper function to calculate the absolute difference between two integers
func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}
