package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// https://adventofcode.com/2024/day/3

const (
	filename           string = "input.txt"                                      // File containing the input data
	prefix             string = "mul("                                           // Prefix for a valid "mul" function call
	suffix             string = ")"                                              // Suffix for a valid "mul" function call
	getAllMuls         string = `(mul\(\d{1,3},\d{1,3}\))`                       // Regex to match all "mul" calls
	getAllMulsDosDonts string = `(mul\(\d{1,3},\d{1,3}\))|(do\(\))|(don\'t\(\))` // Regex to match "mul", "do()", and "don't()" calls
	delim              string = ","                                              // Delimiter separating arguments in "mul" calls
	doFunc             string = "do()"                                           // Keyword to enable processing
	dontFunc           string = "don't()"                                        // Keyword to disable processing
)

func main() {
	// Read and parse instructions from the input file
	instructions, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Part 1: Calculate the sum of valid multiplication results
	var sumProducts int
	mcs, err := extractCalls(getAllMuls, instructions)
	if err != nil {
		fmt.Printf("error extracting mul calls: %w\n", err)
		os.Exit(1)
	}

	for _, mc := range mcs {
		// Extract and validate arguments for each "mul" call
		valid, a, b := getMulArgs(mc)
		if !valid {
			fmt.Printf("error getting mul args from mul call: %s", mc)
		}
		sumProducts += a * b
	}
	fmt.Println(sumProducts)

	// Part 2: Process toggles ("do()" and "don't()") and calculate the sum
	sumProducts = 0
	do := true
	var mcdds []string
	mcdds, err = extractCalls(getAllMulsDosDonts, instructions)
	if err != nil {
		fmt.Printf("error extracting mul calls: %w\n", err)
		os.Exit(1)
	}

	for _, mcdd := range mcdds {
		// Check for "do()" to enable processing
		if strings.Contains(mcdd, doFunc) {
			do = true
			continue
		}
		// Check for "don't()" to disable processing
		if strings.Contains(mcdd, dontFunc) {
			do = false
			continue
		}

		// Process valid "mul" calls when processing is enabled
		if do {
			valid, a, b := getMulArgs(mcdd)
			if !valid {
				fmt.Printf("error getting mul args from mul call: %s", mcdd)
			}
			// Accumulate the product if valid
			sumProducts += a * b
		}
	}
	fmt.Println(sumProducts)
}

// Reads and returns the contents of the input file as a single string
func readInput(fname string) (string, error) {
	var instructions string // String to store all instruction lines

	// Open the input file
	f, err := os.Open(fname)
	if err != nil {
		return "", fmt.Errorf("error opening file [%s]: %w", fname, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Append each line to the instructions string
		instructions += scanner.Text()
	}
	return instructions, nil // Return the full content as a single string
}

// /////////////
// Part 1 & 2 //
// ////// /////

// Extracts all matches of a given regex from the input string
func extractCalls(regex, input string) ([]string, error) {
	var export []string

	// Compile the provided regex pattern
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex string: %s", getAllMuls)
	}

	// Find all matches and return them as a slice
	export = r.FindAllString(input, -1)
	return export, nil
}

// Parses and validates arguments from a "mul" function call
func getMulArgs(call string) (bool, int, int) {
	// Remove the "mul(" prefix
	m, found := strings.CutPrefix(call, prefix)
	if !found {
		return false, 0, 0 // Invalid if the prefix is missing
	}

	// Remove the ")" suffix
	m, found = strings.CutSuffix(m, suffix)
	if !found {
		return false, 0, 0 // Invalid if the suffix is missing
	}

	// Split the arguments using the delimiter ","
	nums := strings.Split(m, delim)

	// Validate the number of arguments and their lengths
	if len(nums) != 2 || len(nums[0]) > 3 || len(nums[1]) > 3 {
		return false, 0, 0 // Invalid if incorrect number or size of arguments
	}

	// Convert arguments to integers
	a, err := strconv.Atoi(nums[0])
	if err != nil {
		return false, 0, 0 // Invalid if the first argument is not an integer
	}
	b, err := strconv.Atoi(nums[1])
	if err != nil {
		return false, 0, 0 // Invalid if the second argument is not an integer
	}

	return true, a, b // Return validity and the parsed arguments
}
