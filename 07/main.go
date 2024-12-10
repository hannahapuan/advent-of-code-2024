package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// https://adventofcode.com/2024/day/7

const (
	filename   string = "input.txt" // Input file name
	delimColon string = ":"         // Delimiter for separating target and values
	delimSpace string = " "         // Delimiter for separating values
	mult       rune   = '*'         // Multiplication operator
	add        rune   = '+'         // Addition operator
)

var (
	ops = []rune{'*', '+'} // List of possible operators
)

type equation struct {
	answer int64   // Target value of the equation
	vals   []int64 // List of values in the equation
}

func main() {
	eqs, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var sum int64
	for _, eq := range eqs {
		// Generate all possible combinations of operators for the given equation
		opCombos := generateCombinations(ops, len(eq.vals)-1)
		// Solve the equations and calculate the total sum of valid answers
		sum = doMath(opCombos, eqs)
	}
	fmt.Println(sum) // Print the total sum of valid answers
}

// Reads the input file and parses it into equations
func readInput(fname string) ([]equation, error) {
	var export []equation // List to store the parsed equations

	// Open the input file
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file [%s]: %w", fname, err)
	}
	defer f.Close() // Ensure the file is closed after the function completes

	scanner := bufio.NewScanner(f)
	var l []byte
	var eq equation
	for scanner.Scan() {
		l = scanner.Bytes() // Read the line as a byte slice
		line := strings.Split(string(l), delimColon)
		ans, err := strconv.Atoi(line[0]) // Parse the target value
		if err != nil {
			return nil, fmt.Errorf("error parsing file [%s]: expected int for first value: %w", fname, err)
		}
		if len(line) != 2 {
			return nil, fmt.Errorf("error unexpected format: %s", line)
		}

		// Parse the list of values
		vsi := make([]int64, 0)
		vs := strings.Split(line[1], delimSpace)
		for _, v := range vs {
			if v == "" {
				continue // Skip empty values
			}
			vi, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("expected int64 [%s]: %w", v, err)
			}
			vsi = append(vsi, int64(vi))
		}

		// Create an equation object
		eq = equation{
			answer: int64(ans),
			vals:   vsi,
		}
		export = append(export, eq) // Add the equation to the list
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file [%s]: %w", fname, err)
	}

	return export, nil
}

// Generates all combinations of operators for the given length
func generateCombinations(elements []rune, length int) [][]rune {
	var result [][]rune // List to store the combinations

	combination := make([]int64, length) // Array to track the current combination indices

	for {
		// Create the current combination based on indices
		current := make([]rune, length)
		for i := 0; i < length; i++ {
			current[i] = elements[combination[i]]
		}
		result = append(result, current) // Add the combination to the result

		// Increment the combination indices
		idx := length - 1
		for idx >= 0 {
			combination[idx]++
			if int(combination[idx]) < len(elements) {
				break
			}
			combination[idx] = 0
			idx--
		}

		// If we are done, break the loop
		if idx < 0 {
			break
		}
	}

	return result
}

// Solves the equations by trying all operator combinations and returns the total sum of valid answers
func doMath(operators [][]rune, equations []equation) int64 {
	var totalSum int64

	for _, equation := range equations {
		for _, operatorCombo := range operators {
			var ans int64
			for i := 0; i < len(equation.vals)-1; i++ {
				if i == len(operatorCombo) {
					break // Avoid index out of range
				}
				if i == 0 {
					ans = equation.vals[i] // Initialize the result with the first value
				}
				newAns := do(operatorCombo[i], ans, equation.vals[i+1]) // Apply the operator
				ans = newAns
			}
			if ans == equation.answer { // Check if the result matches the target
				totalSum += ans
				break // Skip further operator combinations for this equation
			}
		}
	}

	return totalSum
}

// Performs the mathematical operation based on the given operator
func do(operator rune, a, b int64) int64 {
	switch operator {
	case mult:
		return a * b
	case add:
		return a + b
	}
	os.Exit(1) // Exit if an unsupported operator is encountered
	return 0
}
