package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// https://adventofcode.com/2024/day/5

const (
	fileName   string = "easy-example.txt"
	delimPipe  string = "|"
	delimComma string = ","
)

func main() {
	rules, updateLists, err := readInput(fileName)
	if err != nil {
		fmt.Println("error reading input: %w", err)
		os.Exit(1)
	}

	for _, updateList := range updateLists {
		fmt.Println("---------------------")
		applicableRules := calcApplicableRules(rules, updateList)
		dependencyCount := calcDependencyCount(applicableRules)
		// pageToDepsMap := calcPageToDepsMap(applicableRules)
		zeroRulesPages := initZeroRulesPagesQueue(dependencyCount, updateList)
		valid, middleNum := validateUpdateList(updateList, dependencyCount, zeroRulesPages)
		fmt.Printf("valid: %t\n", valid)
		fmt.Printf("middleNum: %t\n", middleNum)
	}
}

// Reads the input file
// outputs:
//
//	rules: list of pairs where [A,B] A must come before B
//	updateLists: list of lists which are the order to be updated
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
		// Read the line and split it into strings based on the delimiter
		line := scanner.Text()

		// read page ordering rules of input first
		if strings.Contains(line, delimPipe) {
			values := strings.Split(line, delimPipe)
			if len(values) != 2 {
				return nil, nil, fmt.Errorf("incorrect format, expected int|int, found %s", line)
			}

			var page, rule int
			page, err = strconv.Atoi(values[0])
			if err != nil {
				return nil, nil, fmt.Errorf("incorrect format, could not parse first int value %s: %w", values[0], err)
			}
			rule, err = strconv.Atoi(values[1])
			if err != nil {
				return nil, nil, fmt.Errorf("incorrect format, could not parse second int value %s: %w", values[1], err)
			}
			rules = append(rules, []int{page, rule})

		} else if strings.Contains(line, delimComma) {
			updateList := strings.Split(line, delimComma)
			var is []int
			is, err = stringSliceToIntSlice(updateList)
			if err != nil {
				return nil, nil, err
			}
			updateLists = append(updateLists, is)
		}
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading file [%s]: %w", fname, err)
	}

	return rules, updateLists, nil // Return the list of parsed reports and update lists
}

func stringSliceToIntSlice(ss []string) ([]int, error) {
	export := make([]int, len(ss))
	var err error
	for i := range ss {
		export[i], err = strconv.Atoi(ss[i])
		if err != nil {
			return nil, fmt.Errorf("error converting string to int: %s", ss[i])
		}
	}
	return export, nil
}

// returns page: number of dependencies
func calcDependencyCount(rules [][]int) map[int]int {
	dc := make(map[int]int)
	for _, rule := range rules {
		page := rule[0]
		dc[page]++
	}
	return dc
}

// returns page: slice of what it is dependent on (all direct children in the directed graph)
func calcPageToDepsMap(rules [][]int) map[int][]int {
	pageToDeps := make(map[int][]int)

	for _, rule := range rules {
		page := rule[0]
		dep := rule[1]
		_, ok := pageToDeps[page]
		if !ok {
			pageToDeps[page] = make([]int, 0)
		}
		pageToDeps[page] = append(pageToDeps[page], dep)
	}
	return pageToDeps
}

// returns
//  1. the rules applicable to only pages included in the update list
func calcApplicableRules(rules [][]int, updateList []int) [][]int {
	var empty struct{}
	applicablePages := make(map[int]struct{})
	applicableRules := make([][]int, 0)

	// get all numbers in the update list
	for _, page := range updateList {
		applicablePages[page] = empty
	}

	for _, rule := range rules {
		_, okPage := applicablePages[rule[0]]
		_, okDep := applicablePages[rule[1]]
		if !okPage || !okDep {
			continue
		}
		applicableRules = append(applicableRules, rule)
	}
	return applicableRules
}

// returns a queue with the pages with zero rules
func initZeroRulesPagesQueue(dependencyCount map[int]int, updateList []int) *list.List {
	zeroRulesPages := list.New()
	for _, page := range updateList {
		_, ok := dependencyCount[page]
		if !ok {
			fmt.Println(page)
			zeroRulesPages.PushBack(page)
		}
	}
	return zeroRulesPages
}

// returns:
//
//  1. middle page number
//  2. if it is valid or not
func validateUpdateList(updateList []int, dependencyCount map[int]int, zeroRulesPages *list.List) (bool, int) {
	readOrder := make([]int, 0)
	printQueue(zeroRulesPages)
	var foundPage int
	for zeroRulesPage := zeroRulesPages.Front(); zeroRulesPage != nil; zeroRulesPage = zeroRulesPage.Next() {
		var found bool
		for _, page := range updateList {
			printQueue(zeroRulesPages)
			fmt.Printf("-zeroRulesPage: %v\n", zeroRulesPage.Value.(int))
			fmt.Printf("page: %v\n", page)
			fmt.Printf("zeroRulesPage.Value.(int) == page: %t\n", zeroRulesPage.Value.(int) == page)
			if zeroRulesPage.Value.(int) == page {
				found = true
				zeroRulesPages.Remove(zeroRulesPage)
				readOrder = append(readOrder, page)
				foundPage = page
			}
		}
		if !found {
			// cycle in graph -- invalid
			return false, 0
		}
		// remove from the count of dependencies for the page found
		for dcPage := range dependencyCount {
			if dcPage == foundPage {
				fmt.Println("here")
				dependencyCount[foundPage]--

				// if this creates zero dependencies, add it to the zeroRulesPages
				if dependencyCount[foundPage] == 0 {
					zeroRulesPages.PushBack(foundPage)
				}
			}
		}
	}

	fmt.Println(readOrder)
	return true, readOrder[len(readOrder)/2]
}

func printQueue(l *list.List) {
	fmt.Println("Zero rules:")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("\t%d\n", e.Value.(int))
	}
}
