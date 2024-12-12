package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// https://adventofcode.com/2024/day/9

const (
	filename     string = "input.txt"
	freeSpaceVal int    = -1
)

func main() {
	blocks, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Println(blocksToString(blocks))

	fmt.Println()
	moved := true
	finishedBlocks := append([]int{}, blocks...)
	for moved {
		blocks, moved = move(blocks)
		// fmt.Println(blocksToString(blocks))
		finishedBlocks = append([]int{}, blocks...)
	}

	fmt.Println("\nchecksum:", calcChecksum(finishedBlocks))
}

func readInput(fname string) ([]int, error) {
	ids := make([]int, 0)

	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file:  %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var fileIDIdx int

	for {
		// Read one character at a time
		lenFileRu, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" { // End of file
				break
			}
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		if lenFileRu == '\n' {
			return ids, nil
		}

		var lenFile, lenFree int
		lenFile, err = strconv.Atoi(string(lenFileRu))
		if err != nil {
			return nil, fmt.Errorf("error converting file length input to int: %w", err)
		}

		for i := 0; i < lenFile; i++ {
			ids = append(ids, fileIDIdx)
		}
		fileIDIdx++

		lenFreeRu, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" { // End of file
				break
			}
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		if lenFreeRu == '\n' {
			return ids, nil
		}

		lenFree, err = strconv.Atoi(string(lenFreeRu))
		if err != nil {
			return nil, fmt.Errorf("error converting free length input to int: %w", err)
		}

		for i := 0; i < lenFree; i++ {
			ids = append(ids, freeSpaceVal)
		}

	}
	return ids, nil
}

func blocksToString(bs []int) string {
	var export string
	for _, b := range bs {
		// free space
		if b == -1 {
			export += "."
			continue
		}
		// file space
		export += fmt.Sprintf("%d", b)
	}
	return export
}

func swap(blocks []int, a, b int) []int {
	blocks[a], blocks[b] = blocks[b], blocks[a]
	return blocks
}

// returns updated blocks after move and if a move happened
func move(blocks []int) ([]int, bool) {
	blocksCopy := append([]int{}, blocks...)

	freeIndex := getFirstFreeSpaceIndex(blocksCopy)
	fileIndex := getLastFileIndex(blocks)

	// all free space is after file space, finished
	if fileIndex < freeIndex {
		return blocks, false
	}

	return swap(blocksCopy, freeIndex, fileIndex), true
}

func getFirstFreeSpaceIndex(blocks []int) int {
	for i := range blocks {
		if blocks[i] != freeSpaceVal {
			continue
		}
		return i
	}
	return -1
}

func getLastFileIndex(blocks []int) int {
	for i := len(blocks) - 1; i >= 0; i-- {
		if blocks[i] == freeSpaceVal {
			continue
		}
		return i
	}
	return -1
}

func calcChecksum(blocks []int) int {
	var checksum int
	for i, id := range blocks {
		if id == -1 {
			continue
		}
		checksum += i * id
	}
	return checksum
}
