package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// https://adventofcode.com/2024/day/2

// Constants
const (
	filename string = "input.txt" // Input file containing the reports
	delim    string = " "         // Delimiter used to split values in the file
)

func main() {
	// Read the input file and parse it into a slice of integer slices
	reports, err := readInput(filename)
	if err != nil {
		fmt.Println("error reading input file: %w", err) // Print error if file reading fails
		os.Exit(1)
	}

	// Calculate the count of "safe" reports
	safeCount := countSafe(reports)
	fmt.Println(safeCount) // Print the number of safe reports
}

// Reads the input file and parses each report
func readInput(fname string) ([][]int, error) {
	var reports [][]int // Slice to hold all parsed reports

	// Open the input file
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file [%s]: %w", fname, err)
	}
	defer f.Close() // Ensure the file is closed after the function completes

	scanner := bufio.NewScanner(f) // Create a scanner to read the file line by line
	for scanner.Scan() {
		// Read the line and split it into strings based on the delimiter
		line := scanner.Text()
		values := strings.Split(line, delim)

		// Parse the string values into integers
		var report []int
		for _, val := range values {
			num, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("error parsing file [%s]: expected int, found %s", fname, val)
			}
			report = append(report, num) // Append the parsed integer to the report
		}

		// Append the parsed report to the list of reports
		reports = append(reports, report)
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file [%s]: %w", fname, err)
	}

	return reports, nil // Return the list of parsed reports
}

// Counts how many reports are "safe"
func countSafe(reports [][]int) int {
	var safeCount int // Counter for safe reports

	// Iterate over all reports
	for _, report := range reports {
		// Check if the report is safe
		if isSafe(report) {
			safeCount++ // Increment the counter if the report is safe
		}
	}

	return safeCount // Return the total count of safe reports
}

// Determines whether a given report is "safe"
func isSafe(report []int) bool {
	// Track if the sequence is decreasing
	var isDecreasing bool
	// Calculate the difference between the first two numbers
	diff := report[0] - report[1]

	// Determine the initial trend (increasing or decreasing)
	if diff < 0 {
		isDecreasing = true
	}

	// Iterate over the report to check the safety conditions
	for i := range report {
		if i == 0 {
			continue // Skip the first element
		}

		// Calculate the difference between the current and previous numbers
		diff = report[i-1] - report[i]
		absDiff := absDiffInt(report[i], report[i-1]) // Absolute difference

		// Condition 1: The absolute difference must be between 1 and 3 (inclusive)
		if absDiff < 1 || absDiff > 3 {
			return false
		}

		// Condition 2: The trend must not switch between increasing and decreasing
		if diff > 0 && isDecreasing { // If it starts increasing after decreasing
			return false
		} else if diff < 0 && !isDecreasing { // If it starts decreasing after increasing
			return false
		}
	}

	// If all conditions are met, the report is safe
	return true
}

// Helper function to calculate the absolute difference between two integers
func absDiffInt(x, y int) int {
	if x < y {
		return y - x // Return positive difference
	}
	return x - y // Return positive difference
}
