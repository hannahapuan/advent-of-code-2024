package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// https://adventofcode.com/2024/day/1

const (
	filename = "input.txt" // Input file containing pairs of integers
	delim    = "   "       // Delimiter used to split input (three spaces)
)

func main() {
	// Read input file and parse it into two lists of integers
	list0, list1, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Part 1
	// Calculate the sum of absolute differences between sorted lists
	fmt.Println("part 1")

	// Create sorted copies of the original lists
	sortedL0 := append([]int{}, list0...)
	sortedL1 := append([]int{}, list1...)
	sort.Ints(sortedL0) // Sort the first list
	sort.Ints(sortedL1) // Sort the second list

	// Calculate and print the result for Part 1
	fmt.Println(partOneBruteForceSolution(sortedL0, sortedL1))

	// Part 2
	// Calculate the similarity score
	fmt.Println("part 2")
	// Create a map from list0 with all keys initialized to 0
	m0 := listToZeroMap(list0)
	// Update the map with the cardinality (frequency) of elements in list1
	m := populateCardinalityFromList(m0, list1)
	// Calculate and print the similarity score based on list0 and the map
	simScore := calcSimScore(m, list0)
	fmt.Println(simScore)
}

// Reads the input file and parses it into two lists of integers
func readInput(fname string) ([]int, []int, error) {
	var l0, l1 []int // Lists to store the parsed integers

	// Open the input file
	f, err := os.Open(fname)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file [%s]: %w", fname, err)
	}
	defer f.Close() // Ensure the file is closed after the function completes

	scanner := bufio.NewScanner(f)
	var l []byte
	for scanner.Scan() {
		l = scanner.Bytes()                  // Read the line as a byte slice
		v := strings.Split(string(l), delim) // Split the line using the delimiter
		// Parse the first value as an integer
		v0, err := strconv.Atoi(v[0])
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing file [%s]: expected int for first value", fname)
		}
		// Parse the second value as an integer
		v1, err := strconv.Atoi(v[1])
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing file [%s]: expected int for second value", fname)
		}
		// Ensure the line contains exactly two values
		if len(v) != 2 {
			return nil, nil, fmt.Errorf("error parsing file [%s]: unexpected format", fname)
		}

		// Append the parsed values to the respective lists
		l0 = append(l0, v0)
		l1 = append(l1, v1)
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading file [%s]: %w", fname, err)
	}
	// Ensure the lists have the same length
	if len(l0) != len(l1) {
		return nil, nil, fmt.Errorf("unexpected format, lists are expected to be identical lengths")
	}

	return l0, l1, nil
}

////////////
// Part 1 //
////////////

// Calculate the sum of absolute differences between corresponding elements in two lists
func partOneBruteForceSolution(list0, list1 []int) int {
	var distance int
	// Iterate over both lists and calculate the absolute difference
	for i := range list0 {
		distance += absDiffInt(list0[i], list1[i])
	}
	return distance
}

// Helper function to calculate the absolute difference between two integers
func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

// //////////
// Part 2  //
// //////////

// Creates a map with all keys from the input list initialized to 0
func listToZeroMap(list []int) map[int]int {
	m := make(map[int]int) // Initialize an empty map
	for _, lid := range list {
		if _, ok := m[lid]; !ok {
			m[lid] = 0 // Set the initial value for the key
		}
	}
	return m
}

// Updates the map with the cardinality (frequency) of elements in the second list
func populateCardinalityFromList(m map[int]int, list []int) map[int]int {
	for _, lid := range list {
		m[lid]++ // Increment the value for the key
	}
	return m
}

// Calculates the similarity score using the map and the first list
func calcSimScore(m map[int]int, list []int) int {
	var simScore int
	// Iterate over the first list and calculate the score
	for _, lid := range list {
		if count, ok := m[lid]; ok { // Check if the key exists in the map
			simScore += count * lid // Multiply the count by the key and add to the score
		}
	}
	return simScore
}
