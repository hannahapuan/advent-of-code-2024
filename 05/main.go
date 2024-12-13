package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Advent of Code 2024 - Day 5: Challenge
// https://adventofcode.com/2024/day/5

const (
	fileName   string = "easy-example.txt"
	delimPipe  string = "|"
	delimComma string = ","
)

func main() {
	// Read input rules and update lists from the file
	rules, updateLists, err := readInput(fileName)
	if err != nil {
		fmt.Println("error reading input: %w", err)
		os.Exit(1)
	}

	// Iterate over each update list
	for _, updateList := range updateLists {
		fmt.Println("---------------------")
		// Filter rules to those relevant for the current update list
		applicableRules := calcApplicableRules(rules, updateList)
		fmt.Println("applicableRules")
		fmt.Println(applicableRules)

		// Create a dependency count map (page: number of dependencies)
		dependencyCount := calcDependencyCount(applicableRules)
		fmt.Println("dependencyCount")
		fmt.Println(dependencyCount)

		// Create an adjacency list (page: dependencies)
		pageToDepsMap := calcPageToDepsMap(applicableRules)
		fmt.Println("pageToDepsMap")
		fmt.Println(pageToDepsMap)

		// Initialize a queue for pages with zero dependencies
		zeroRulesPages := initZeroRulesPagesQueue(dependencyCount, updateList)
		fmt.Println("zeroRulesPages")
		printQueue(zeroRulesPages)

		// Validate the current update list against the rules
		valid, middleNum := validateUpdateList(updateList, pageToDepsMap, dependencyCount, zeroRulesPages)

		// Output the validation result and the middle number
		fmt.Printf("valid: %t\n", valid)
		fmt.Printf("middleNum: %d\n", middleNum)
	}
}

// Reads the input file
// Outputs:
// - rules: list of pairs where [A,B] A must come before B
// - updateLists: list of lists which are the order to be updated
func readInput(fname string) ([][]int, [][]int, error) {
	rules := make([][]int, 0)
	updateLists := make([][]int, 0)

	// Open the input file
	f, err := os.Open(fname)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file [%s]: %w", fname, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		// Read page ordering rules (e.g., "A|B")
		if strings.Contains(line, delimPipe) {
			values := strings.Split(line, delimPipe)
			if len(values) != 2 {
				return nil, nil, fmt.Errorf("incorrect format, expected int|int, found %s", line)
			}

			page, err1 := strconv.Atoi(values[0])
			rule, err2 := strconv.Atoi(values[1])
			if err1 != nil || err2 != nil {
				return nil, nil, fmt.Errorf("error parsing rules: %s", line)
			}
			rules = append(rules, []int{page, rule})

		} else if strings.Contains(line, delimComma) {
			// Read update list (e.g., "1,2,3")
			updateList := strings.Split(line, delimComma)
			intList, err := stringSliceToIntSlice(updateList)
			if err != nil {
				return nil, nil, err
			}
			updateLists = append(updateLists, intList)
		}
	}

	// Check for errors in reading the file
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading file [%s]: %w", fname, err)
	}

	return rules, updateLists, nil
}

// Converts a slice of strings to a slice of integers
func stringSliceToIntSlice(ss []string) ([]int, error) {
	export := make([]int, len(ss))
	for i, s := range ss {
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("error converting string to int: %s", s)
		}
		export[i] = val
	}
	return export, nil
}

// Returns a map of page to the number of dependencies (in-degree)
func calcDependencyCount(rules [][]int) map[int]int {
	dc := make(map[int]int)
	for _, rule := range rules {
		page := rule[1]
		dc[page]++ // Increment dependency count for the dependent page
	}
	return dc
}

// Returns a map of page to its direct dependencies (adjacency list)
func calcPageToDepsMap(rules [][]int) map[int][]int {
	pageToDeps := make(map[int][]int)
	for _, rule := range rules {
		page, dep := rule[0], rule[1]
		pageToDeps[page] = append(pageToDeps[page], dep)
	}
	return pageToDeps
}

// Returns the subset of rules applicable to the given update list
func calcApplicableRules(rules [][]int, updateList []int) [][]int {
	var empty struct{}
	applicablePages := make(map[int]struct{})
	applicableRules := make([][]int, 0)

	// Mark pages present in the update list
	for _, page := range updateList {
		applicablePages[page] = empty
	}

	// Filter rules to only include pages in the update list
	for _, rule := range rules {
		if _, okPage := applicablePages[rule[0]]; okPage {
			if _, okDep := applicablePages[rule[1]]; okDep {
				applicableRules = append(applicableRules, rule)
			}
		}
	}
	return applicableRules
}

// Initializes a queue with pages that have zero dependencies
func initZeroRulesPagesQueue(dependencyCount map[int]int, updateList []int) *list.List {
	zeroRulesPages := list.New()
	for _, page := range updateList {
		if dependencyCount[page] == 0 {
			zeroRulesPages.PushBack(page)
		}
	}
	return zeroRulesPages
}

// Updates the dependency count map based on the adjacency list
func updateDepsMap(dependencyCount map[int]int, pageToDepsMap map[int][]int) map[int]int {
	for page, deps := range pageToDepsMap {
		dependencyCount[page] = len(deps)
	}
	return dependencyCount
}

// Validates the update list using Kahn's algorithm
// Outputs:
// - valid: whether the update list is valid
// - middleNum: the middle page of the update order
func validateUpdateList(updateList []int, pageToDepsMap map[int][]int, dependencyCount map[int]int, zeroRulesPages *list.List) (bool, int) {
	fmt.Println("validateUpdateList()")
	readOrder := make([]int, 0)

	printQueue(zeroRulesPages)
	// Process pages in topological order
	for _, page := range updateList {
		zeroRulesPage := zeroRulesPages.Front()
		// for zeroRulesPage := zeroRulesPages.Front(); zeroRulesPage != nil; zeroRulesPage = zeroRulesPage.Next() {
		// page := zeroRulesPage.Value.(int)
		fmt.Println("PAGE")
		fmt.Println(page)

		// Append the current page to the read order
		readOrder = append(readOrder, page)
		// fmt.Println("readOrder")
		// fmt.Println(readOrder)

		zeroRulesPages.Remove(zeroRulesPage)
		// fmt.Println("zeroRulesPages.Remove(zeroRulesPage)")
		// printQueue(zeroRulesPages)

		// Reduce the dependency count for all dependents of the current page
		fmt.Println("dependencyCount[page]")
		fmt.Println(dependencyCount[page])
		dependencyCount[page]--
		if dependencyCount[page] == 0 {
			zeroRulesPages.PushBack(page)
			for p, deps := range pageToDepsMap {
				for i, dep := range deps {
					if dep == page {
						// remove the current page from the list of existing pages
						copy := append([]int{}, remove(deps, i)...)
						// add the new deps map back to the pageToDepsMap
						pageToDepsMap[p] = copy
					}
				}
			}
			zeroRulesPages = initZeroRulesPagesQueue(dependencyCount, updateList)
		}

		// reset zerorulespages stuff
		zeroRulesPage = zeroRulesPage.Next()
		if zeroRulesPage == nil {
			break
		}

	}
	// If not all pages were processed, the update list is invalid
	if len(readOrder) != len(updateList) {
		return false, 0
	}

	// Return the middle page in the sorted order
	return true, readOrder[len(readOrder)/2]
}

// Prints the elements of a queue
func printQueue(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("\t%d\n", e.Value.(int))
	}
}

// Removes a specific value from a slice
func removeVal(s []int, v int) []int {
	i := indexOf(s, v)
	return remove(s, i)
}

// Removes an element at a specific index from a slice
func remove(s []int, i int) []int {
	return append(s[:i], s[i+1:]...)
}

// Finds the index of a value in a slice
func indexOf(s []int, v int) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}
