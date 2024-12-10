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
	filename   string = "input.txt"
	delimColon string = ":"
	delimSpace string = " "
	mult       rune   = '*'
	add        rune   = '+'
)

var (
	ops = []rune{'*', '+'}
)

type equation struct {
	answer int
	vals   []int
}

func main() {
	eqs, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(len(eqs))

	var sum int
	for _, eq := range eqs {
		opCombos := generateCombinationsOfLength(ops, len(eq.vals)-1)
		sum = doMath(opCombos, eqs)
	}

	fmt.Println(sum)

}

// Reads the input file and parses it into two lists of integers
func readInput(fname string) ([]equation, error) {
	var export []equation // Lists to store the parsed integers

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
		ans, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing file [%s]: expected int for first value: %w", fname, err)
		}
		if len(line) != 2 {
			return nil, fmt.Errorf("error unexpected format: %s", line)
		}

		vsi := make([]int, 0)
		vs := strings.Split(line[1], delimSpace)
		for _, v := range vs {
			if v == "" {
				continue
			}
			vi, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("expected int [%s]: %w", v, err)
			}
			vsi = append(vsi, vi)
		}

		eq = equation{
			answer: ans,
			vals:   vsi,
		}
		export = append(export, eq)
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file [%s]: %w", fname, err)
	}

	return export, nil
}

func generateCombinationsOfLength(chars []rune, maxLength int) [][]rune {
	results := [][]rune{{}}

	// Build combinations iteratively until we reach the desired length
	for length := 1; length <= maxLength; length++ {
		tempResults := make([][]rune, 0)
		for _, existing := range results {
			if len(existing) == length-1 {
				for _, char := range chars {
					newCombination := append([]rune{}, existing...)
					newCombination = append(newCombination, char)
					tempResults = append(tempResults, newCombination)
				}
			}
		}
		results = tempResults
	}

	return results
}

func doMath(operators [][]rune, equations []equation) int {
	var totalSum int

	for _, equation := range equations {
		for _, operatorCombo := range operators {
			ans := 0
			for i := 0; i < len(equation.vals)-1; i++ {
				if i == len(operatorCombo) {
					break
				}
				if i == 0 {
					ans = equation.vals[i]
				}
				newAns := do(operatorCombo[i], ans, equation.vals[i+1])
				// fmt.Printf("%d %s %d = %d\n", ans, string(operatorCombo[i]), equation.vals[i+1], newAns)
				ans = newAns
			}
			// fmt.Printf("\t%d == %d = %t\n", ans, equation.answer, ans == equation.answer)
			if ans == equation.answer {
				// fmt.Println("\tsolved\n")
				totalSum += ans
				break
			}
			// fmt.Println()
		}

	}

	return totalSum
}

func do(operator rune, a, b int) int {
	switch operator {
	case mult:
		return a * b
	case add:
		return a + b
	}
	os.Exit(1)
	return 0
}
