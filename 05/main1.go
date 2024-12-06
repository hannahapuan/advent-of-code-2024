// package main

// import (
// 	"bufio"
// 	"container/heap"
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
// 	pq, uls, err := readInputAndHeapify(fileName)
// 	if err != nil {
// 		fmt.Println("error reading input and heapifying: %w", err)
// 		os.Exit(1)
// 	}
// 	for _, p := range pq {
// 		fmt.Println(toString(*p))
// 	}
// 	pim := pqIndexMap(pq)
// 	fmt.Println(pim)
// 	fmt.Println(uls)
// 	fmt.Println(isValidUpdateList(uls[0], pim))
// }

// // Reads the input file
// func readInputAndHeapify(fname string) (PriorityQueue, [][]int, error) {
// 	// var protocols []protocol
// 	pq := make(PriorityQueue, 0)
// 	updateLists := make([][]int, 0)

// 	// Open the input file
// 	f, err := os.Open(fname)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("error opening file [%s]: %w", fname, err)
// 	}
// 	defer f.Close()

// 	scanner := bufio.NewScanner(f)

// 	var pqIndex int
// 	for scanner.Scan() {
// 		// Read the line and split it into strings based on the delimiter
// 		line := scanner.Text()

// 		// read page ordering rules of input first
// 		var page, pri int
// 		if strings.Contains(line, delimPipe) {
// 			values := strings.Split(line, delimPipe)
// 			if len(values) != 2 {
// 				return nil, nil, fmt.Errorf("incorrect format, expected int|int, found %s", line)
// 			}

// 			page, err = strconv.Atoi(values[0])
// 			if err != nil {
// 				return nil, nil, fmt.Errorf("incorrect format, could not parse first int value %s: %w", values[0], err)
// 			}
// 			pri, err = strconv.Atoi(values[1])
// 			if err != nil {
// 				return nil, nil, fmt.Errorf("incorrect format, could not parse second int value %s: %w", values[1], err)
// 			}
// 			pro := protocol{
// 				page:     page,
// 				priority: pri,
// 				index:    pqIndex,
// 			}
// 			heap.Push(&pq, &pro)
// 			pqIndex++
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

// 	heap.Init(&pq)
// 	return pq, updateLists, nil // Return the list of parsed reports and update lists
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

// // pqIndexMap from a PQ returns a map with page:index
// func pqIndexMap(pq PriorityQueue) map[int]int {
// 	for _, p := range pq {
// 		fmt.Println(toString(*p))
// 	}
// 	export := make(map[int]int)
// 	var count int
// 	for pq.Len() > 0 {
// 		p := heap.Pop(&pq).(*protocol)
// 		export[p.page] = p.priority
// 		count++
// 	}

// 	return export
// }

// func isValidUpdateList(ul []int, pti map[int]int) bool {
// 	prevPageIndex := -1
// 	for _, page := range ul {
// 		if pti[page] < prevPageIndex {
// 			return false
// 		}
// 		prevPageIndex = pti[page]
// 	}
// 	return true
// }

// //////// Modified Priority Queue ////////

// // modified version of https://pkg.go.dev/container/heap

// type protocol struct {
// 	page     int
// 	priority int
// 	index    int
// }

// // A PriorityQueue implements heap.Interface and holds protocols.
// type PriorityQueue []*protocol

// func (pq PriorityQueue) Len() int { return len(pq) }

// func (pq PriorityQueue) Less(i, j int) bool {
// 	return pq[i].priority < pq[j].priority
// }

// func (pq PriorityQueue) Swap(i, j int) {
// 	pq[i], pq[j] = pq[j], pq[i]
// 	pq[i].index = i
// 	pq[j].index = j
// }

// func (pq *PriorityQueue) Push(x any) {
// 	n := len(*pq)
// 	item := x.(*protocol)
// 	item.index = n
// 	*pq = append(*pq, item)
// }

// func (pq *PriorityQueue) Pop() any {
// 	old := *pq
// 	n := len(old)
// 	item := old[n-1]
// 	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
// 	item.index = -1 // for safety
// 	*pq = old[0 : n-1]
// 	return item
// }

// //////// Debug Funcs ////////

// func toString(pro protocol) string {
// 	return fmt.Sprintf("%#v", pro)
// }
