package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Constants
const (
	filename string = "input.txt" // Input file containing instructions
	prefix   string = "mul("      // Prefix indicating the start of a multiplication function call
	suffix   string = ")"         // Suffix indicating the end of a multiplication function call
	delim    string = ","         // Delimiter separating arguments in the multiplication function call
)

func main() {
	// Read the input file and parse it into a slice of strings (instructions)
	instructions, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Part 1
	// Process multiplication function calls and calculate their sum
	var sumOfMuls int     // Variable to store the sum of all products
	var mulCalls []string // Slice to store individual "mul" function calls
	for _, ins := range instructions {
		// Extract all "mul" function calls from the instruction line
		mulCalls = getMulFuncCalls(ins)
		for _, mc := range mulCalls {
			// Parse arguments of the "mul" function call
			valid, x, y := getMulArgs(mc)
			if !valid {
				// Skip if the format of the function call is invalid
				continue
			}
			// Calculate the product of the two arguments
			product := x * y
			// Add the product to the running sum
			sumOfMuls += product
		}
	}
	// Output the sum of all valid "mul" function call products
	fmt.Println(sumOfMuls)

	// Part 2
	// TODO
}

// Reads the input file and returns a slice of strings representing each instruction line
func readInput(fname string) ([]string, error) {
	var instructions []string // Slice to hold all instruction lines

	// Open the input file
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file [%s]: %w", fname, err)
	}
	defer f.Close() // Ensure the file is closed after the function completes

	scanner := bufio.NewScanner(f) // Create a scanner to read the file line by line
	for scanner.Scan() {
		line := scanner.Text()                    // Read the current line as a string
		instructions = append(instructions, line) // Append the line to the list of instructions
	}
	return instructions, nil // Return the list of parsed instructions
}

// //////////
// Part 1 //
// //////////

// Extracts all "mul" function calls from a line of instruction
func getMulFuncCalls(li string) []string {
	var export []string                   // Slice to store extracted function calls
	mulCalls := strings.Split(li, prefix) // Split the line using the "mul(" prefix

	// Iterate over the split parts to identify and reassemble valid "mul" calls
	for i := range mulCalls {
		indexEndParen := strings.Index(mulCalls[i], suffix) // Find the closing parenthesis
		if indexEndParen < 0 {
			continue // Skip if no closing parenthesis is found
		}

		// Reassemble the function call with the prefix and suffix
		export = append(export, fmt.Sprintf("%s%s", prefix, mulCalls[i][:indexEndParen+1]))
	}

	return export // Return the list of extracted function calls
}

// Parses the arguments of a "mul" function call
func getMulArgs(call string) (bool, int, int) {
	// Remove the prefix ("mul(") from the function call
	m, found := strings.CutPrefix(call, prefix)
	if !found {
		// Invalid format: prefix not found
		return false, 0, 0
	}

	// Remove the suffix (")") from the function call
	m, found = strings.CutSuffix(m, suffix)
	if !found {
		// Invalid format: suffix not found
		return false, 0, 0
	}

	// Split the remaining string into arguments using the delimiter (",")
	nums := strings.Split(m, delim)

	// Validate the number and length of arguments
	if len(nums) != 2 || len(nums[0]) > 3 || len(nums[1]) > 3 {
		// Invalid format: incorrect number of arguments or argument length exceeds 3 characters
		return false, 0, 0
	}

	// Convert the arguments to integers
	a, err := strconv.Atoi(nums[0])
	if err != nil {
		// Invalid argument: first argument is not a valid integer
		return false, 0, 0
	}
	b, err := strconv.Atoi(nums[1])
	if err != nil {
		// Invalid argument: second argument is not a valid integer
		return false, 0, 0
	}

	return true, a, b // Return validity and the parsed arguments
}

// //////////
// Part 2 //
// //////////
// TODO: Part 2 functionality placeholder
