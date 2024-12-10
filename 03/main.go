package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Constants for configuration
const (
	filename string = "input.txt" // Path to the input file
	prefix   string = "mul("      // Denotes the start of a multiplication function call
	suffix   string = ")"         // Denotes the end of a multiplication function call
	delim    string = ","         // Separator between arguments in the multiplication function call
	doFunc   string = "do()"      // Keyword to enable processing
	dontFunc string = "don't()"   // Keyword to disable processing
)

func main() {
	// Read instructions from the input file
	instructions, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Part 1: Calculate the sum of all valid multiplications
	var sumOfMuls int     // Total sum of all multiplication results
	var mulCalls []string // Stores each valid multiplication function call
	for _, ins := range instructions {
		// Identify all "mul" function calls in the current instruction line
		mulCalls = getMulFuncCalls(ins)
		for _, mc := range mulCalls {
			// Extract and validate arguments from the "mul" function call
			valid, x, y := getMulArgs(mc)
			if !valid {
				continue // Skip invalid "mul" calls
			}
			// Multiply the arguments and add the result to the total sum
			sumOfMuls += x * y
		}
	}
	// Output the total for Part 1
	fmt.Println(sumOfMuls)

	// Part 2: Calculate the sum while respecting "do" and "don't" toggles
	sumOfMuls = 0 // Reset the sum for Part 2
	do := true
	for _, ins := range instructions {
		// Identify "mul" function calls while respecting toggles
		mulCalls = getMulFuncCallsWithDoOrDont(ins, do)
		for _, mc := range mulCalls {
			// Extract and validate arguments from the "mul" function call
			valid, x, y := getMulArgs(mc)
			if !valid {
				continue // Skip invalid "mul" calls
			}
			// Multiply the arguments and add the result to the total sum
			sumOfMuls += x * y
		}
	}
	// Output the total for Part 2
	fmt.Println(sumOfMuls)
}

// Reads the input file and splits it into lines of instructions
func readInput(fname string) ([]string, error) {
	var instructions []string // Holds all instruction lines

	// Open the input file
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file [%s]: %w", fname, err)
	}
	defer f.Close() // Ensure the file is closed when the function exits

	// Read the file line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Append each line of the file to the list of instructions
		instructions = append(instructions, scanner.Text())
	}
	return instructions, nil // Return the parsed instructions
}

// Part 1: Extracts all valid "mul" function calls from a line
func getMulFuncCalls(li string) []string {
	var export []string                   // Stores valid "mul" function calls
	mulCalls := strings.Split(li, prefix) // Split the line by the "mul(" prefix

	// Reassemble valid "mul" calls by finding their suffix
	for i := range mulCalls {
		indexEndParen := strings.Index(mulCalls[i], suffix) // Look for the closing parenthesis
		if indexEndParen < 0 {
			continue // Skip parts without a valid closing parenthesis
		}
		// Append the reconstructed "mul" call to the result
		export = append(export, fmt.Sprintf("%s%s", prefix, mulCalls[i][:indexEndParen+1]))
	}
	return export // Return the valid "mul" function calls
}

// Extracts arguments from a "mul" function call and validates them
func getMulArgs(call string) (bool, int, int) {
	// Remove the "mul(" prefix
	m, found := strings.CutPrefix(call, prefix)
	if !found {
		return false, 0, 0 // Invalid if prefix is missing
	}

	// Remove the ")" suffix
	m, found = strings.CutSuffix(m, suffix)
	if !found {
		return false, 0, 0 // Invalid if suffix is missing
	}

	// Split the remaining string into arguments
	nums := strings.Split(m, delim)

	// Validate the number of arguments and their lengths
	if len(nums) != 2 || len(nums[0]) > 3 || len(nums[1]) > 3 {
		return false, 0, 0 // Invalid if incorrect number of arguments or length
	}

	// Convert arguments to integers
	a, err := strconv.Atoi(nums[0])
	if err != nil {
		return false, 0, 0 // Invalid if the first argument is not a number
	}
	b, err := strconv.Atoi(nums[1])
	if err != nil {
		return false, 0, 0 // Invalid if the second argument is not a number
	}

	return true, a, b // Return validity and the parsed arguments
}

// Part 2: Extracts valid "mul" calls while respecting "do" and "don't" toggles
func getMulFuncCallsWithDoOrDont(li string, do bool) []string {
	var export []string // Stores valid "mul" function calls
	// rule 0: starts with mul(
	mulCalls := strings.Split(li, prefix) // Split the line by the "mul(" prefix
	valid := true

	// Process each part of the split input
	for i := range mulCalls {
		// validation
		valid = true
		// rule 1: end parenthesis exists
		indexEndParen := strings.Index(mulCalls[i], suffix) // Look for the closing parenthesis
		if indexEndParen < 0 {
			valid = false // Skip invalid parts
		} else {
			// Split the remaining string into arguments
			nums := strings.Split(mulCalls[i][:indexEndParen], delim)
			// rule 2: each argument is 1-3 decimals long
			if valid && len(nums) != 2 || len(nums[0]) > 3 || len(nums[1]) > 3 {
				// Validate the number of arguments and their lengths{
				valid = false
			}
		}

		fmt.Println(valid)

		if do && valid {
			// Add the "mul" call to the output if processing is enabled
			// Construct the full "mul" function call
			call := fmt.Sprintf("%s%s%s", prefix, mulCalls[i][:indexEndParen], suffix)
			export = append(export, call)
		}

		// Toggle processing mode based on the presence of "do()" or "don't()"
		// indexOfLastDo := strings.LastIndex(mulCalls[i], doFunc)
		// indexOfLastDont := strings.LastIndex(mulCalls[i], dontFunc)

		// Check for toggles (do() or don't()) in the current segment
		if strings.Contains(mulCalls[i], dontFunc) {
			do = false // Disable processing
		}
		if strings.Contains(mulCalls[i], doFunc) {
			do = true // Enable processing
		}
	}
	return export // Return the valid "mul" function calls
}
