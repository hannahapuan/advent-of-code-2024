// package main

// import (
// 	"bufio"
// 	"container/list"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// // https://adventofcode.com/2024/day/5

// const (
// 	fileName   string = "example.txt"
// 	delimPipe  string = "|"
// 	delimComma string = ","
// )

// func main() {
// 	protocols, updateLists, err := readInput(fileName)
// 	if err != nil {
// 		fmt.Println("error reading input: %w", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(protocols)
// 	fmt.Println(updateLists)

// 	pageToAdjPages, pageToDependencyNumber, err := processProtocols(protocols)
// 	if err != nil {
// 		fmt.Println("error processing protocols: %w", err)
// 		os.Exit(1)
// 	}

// 	fmt.Println(pageToDependencyNumberString(pageToDependencyNumber))
// 	fmt.Println("pageToAdjPages")
// 	fmt.Println(pageToAdjPages)
// 	fmt.Println("\nupdateLists")
// 	fmt.Println(updateLists)
// 	fmt.Println()

// 	var middlePageSums int
// 	for _, updateList := range updateLists {
// 		fmt.Println("\noooooooooooooooooooooooooooooooooooooooo")
// 		fmt.Println("new update list")
// 		fmt.Println("oooooooooooooooooooooooooooooooooooooooo")

// 		valid, middlePage := validateUpdateList(pageToAdjPages, pageToDependencyNumber, updateList)
// 		if valid {
// 			middlePageSums += middlePage
// 		}
// 	}
// 	fmt.Printf("\nanswer: %d\n", middlePageSums)
// }

// // Reads the input file
// // outputs:
// //
// //	protocols: list of pairs where [A,B] A must come before B
// //	updateLists: list of lists which are the order to be updated
// func readInput(fname string) ([][]int, [][]int, error) {
// 	protocols := make([][]int, 0)
// 	updateLists := make([][]int, 0)

// 	// Open the input file
// 	f, err := os.Open(fname)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("error opening file [%s]: %w", fname, err)
// 	}
// 	defer f.Close()

// 	scanner := bufio.NewScanner(f)

// 	for scanner.Scan() {
// 		// Read the line and split it into strings based on the delimiter
// 		line := scanner.Text()

// 		// read page ordering rules of input first
// 		if strings.Contains(line, delimPipe) {
// 			values := strings.Split(line, delimPipe)
// 			if len(values) != 2 {
// 				return nil, nil, fmt.Errorf("incorrect format, expected int|int, found %s", line)
// 			}

// 			var page, rule int
// 			page, err = strconv.Atoi(values[0])
// 			if err != nil {
// 				return nil, nil, fmt.Errorf("incorrect format, could not parse first int value %s: %w", values[0], err)
// 			}
// 			rule, err = strconv.Atoi(values[1])
// 			if err != nil {
// 				return nil, nil, fmt.Errorf("incorrect format, could not parse second int value %s: %w", values[1], err)
// 			}
// 			protocols = append(protocols, []int{page, rule})

// 		} else if strings.Contains(line, delimComma) {
// 			updateList := strings.Split(line, delimComma)
// 			var is []int
// 			is, err = stringSliceToIntSlice(updateList)
// 			if err != nil {
// 				return nil, nil, err
// 			}
// 			updateLists = append(updateLists, is)
// 		}
// 	}

// 	// Check for any errors encountered during scanning
// 	if err := scanner.Err(); err != nil {
// 		return nil, nil, fmt.Errorf("error reading file [%s]: %w", fname, err)
// 	}

// 	return protocols, updateLists, nil // Return the list of parsed reports and update lists
// }

// func stringSliceToIntSlice(ss []string) ([]int, error) {
// 	export := make([]int, len(ss))
// 	var err error
// 	for i := range ss {
// 		export[i], err = strconv.Atoi(ss[i])
// 		if err != nil {
// 			return nil, fmt.Errorf("error converting string to int: %s", ss[i])
// 		}
// 	}
// 	return export, nil
// }

// // outputs:
// //
// //	pageToAdjPages: adjacency list
// //	pageToDependencyNumber: in-degree map
// func processProtocols(protocols [][]int) (map[int][]int, map[int]int, error) {
// 	pageToAdjPages := make(map[int][]int)
// 	pageToDependencyNumber := make(map[int]int)

// 	for _, protocol := range protocols {
// 		if len(protocol) != 2 {
// 			return nil, nil, fmt.Errorf("protocol not of length 2, length is %d", len(protocol))
// 		}

// 		pageToAdjPages[protocol[0]] = append(pageToAdjPages[protocol[0]], protocol[1])
// 		pageToDependencyNumber[protocol[1]]++
// 		if _, ok := pageToDependencyNumber[protocol[0]]; !ok {
// 			pageToDependencyNumber[protocol[0]] = 0
// 		}
// 	}

// 	return pageToAdjPages, pageToDependencyNumber, nil
// }

// // topological sort using Kahn's algorithm
// // output:
// //
// //	valid: if it is a valid updateList
// //	middle page
// func validateUpdateList(pageToAdjPages map[int][]int, pageToDependencyNumber map[int]int, updateList []int) (bool, int) {
// 	ptdn := make(map[int]int)
// 	readOrder := make([]int, 0)

// 	for page, numDeps := range pageToDependencyNumber {
// 		ptdn[page] = numDeps
// 	}

// 	pagesWithNoDependencies := list.New()
// 	for page, degree := range ptdn {
// 		if degree == 0 {
// 			pagesWithNoDependencies.PushBack(page)
// 		}
// 	}

// 	fmt.Printf("updateList: %v\n", updateList)
// 	// print pageswithnodeps
// 	for pageWithNoDependencies := pagesWithNoDependencies.Front(); pageWithNoDependencies != nil; pageWithNoDependencies = pageWithNoDependencies.Next() {
// 		if p, ok := pageWithNoDependencies.Value.(int); ok {
// 			fmt.Printf("p=%v\n", p)
// 		}
// 	}
// 	for _, page := range updateList {
// 		// check if the page has a zero degree value
// 		found := false
// 		for pageWithNoDependencies := pagesWithNoDependencies.Front(); pageWithNoDependencies != nil; pageWithNoDependencies = pageWithNoDependencies.Next() {
// 			if pageWithNoDependencies.Value.(int) == page {
// 				found = true
// 				readOrder = append(readOrder, page)
// 				fmt.Printf("readOrder: %v\n", readOrder)

// 				pagesWithNoDependencies.Remove(pageWithNoDependencies)
// 				break
// 			}
// 		}
// 		if !found {
// 			fmt.Println("grr")
// 			fmt.Println(page)
// 			return false, 0
// 		}
// 		fmt.Printf("pageToAdjPages[%d]=%v\n", page, pageToAdjPages[page])
// 		for _, adjPage := range pageToAdjPages[page] {
// 			fmt.Printf("adjPage: %d\n", page)
// 			ptdn[adjPage]--
// 			if ptdn[adjPage] == 0 {
// 				pagesWithNoDependencies.PushBack(adjPage)
// 			}
// 		}
// 	}

// 	fmt.Printf("readOrder: %v\n", readOrder)
// 	fmt.Println("valid")
// 	return true, readOrder[(len(readOrder) / 2)]
// }

// func pageToDependencyNumberString(ptdn map[int]int) string {
// 	export := "\n"
// 	for p, dn := range ptdn {
// 		export += fmt.Sprintf("page=[%d]:[%d]\n", p, dn)
// 	}
// 	return export
// }
