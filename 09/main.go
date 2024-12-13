package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Advent of Code 2024 - Day 9: Challenge
// Link: https://adventofcode.com/2024/day/9

const (
	filename     string = "input.txt" // Name of the input file containing data
	freeSpaceVal int    = -1          // Represents free space in the blocks
)

// Entry point for the program
func main() {
	// Read and parse the input file into a slice of blocks
	blocks, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Part 1: Move blocks and calculate checksum
	moved := true
	finishedBlocks := append([]int{}, blocks...) // Copy of blocks for manipulation
	for moved {
		finishedBlocks, moved = move(finishedBlocks)        // Move blocks until no movement occurs
		finishedBlocks = append([]int{}, finishedBlocks...) // Ensure copy consistency
	}
	fmt.Println("part 1 checksum:", calcChecksum(finishedBlocks))

	// Part 2: Rearrange blocks and calculate checksum
	idToSize := getIDToSize(blocks)              // Map block IDs to their sizes
	fileEndIndices := getLastFileIndices(blocks) // Get indices of the last files
	fsis := getFreeSpaceIndices(blocks)          // Get indices of free space
	for _, fei := range fileEndIndices {
		fmt.Printf(".")                                      // Progress indicator
		blocks = moveWholeBlock(blocks, fsis, fei, idToSize) // Move entire blocks
	}
	fmt.Println("finished!\n")

	// Calculate checksum for Part 2
	fmt.Println("part 2 checksum:", calcChecksum(blocks))
}

// Reads the input file and parses it into a slice of blocks
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
		// Read one character (file length) at a time
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

		// Parse file length
		var lenFile, lenFree int
		lenFile, err = strconv.Atoi(string(lenFileRu))
		if err != nil {
			return nil, fmt.Errorf("error converting file length input to int: %w", err)
		}

		// Add blocks for the file
		for i := 0; i < lenFile; i++ {
			ids = append(ids, fileIDIdx)
		}
		fileIDIdx++

		// Parse free space length
		lenFreeRu, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
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

		// Add blocks for the free space
		for i := 0; i < lenFree; i++ {
			ids = append(ids, freeSpaceVal)
		}
	}
	return ids, nil
}

// Converts blocks to a string representation
func blocksToString(bs []int) string {
	var export string
	for _, b := range bs {
		if b == -1 {
			export += "." // Free space
			continue
		}
		export += fmt.Sprintf("%d", b) // File block
	}
	return export
}

// Swaps two elements in a slice
func swap(blocks []int, a, b int) []int {
	blocks[a], blocks[b] = blocks[b], blocks[a]
	return blocks
}

// Moves blocks to fill free spaces and returns updated blocks and movement status
func move(blocks []int) ([]int, bool) {
	blocksCopy := append([]int{}, blocks...)

	freeStartIndex := getFirstFreeSpaceIndex(blocksCopy)
	fileIndex := getLastFileIndex(blocks)

	// All free spaces are after file spaces, so no movement needed
	if fileIndex < freeStartIndex {
		return blocks, false
	}

	return swap(blocksCopy, freeStartIndex, fileIndex), true
}

// Creates a map of block IDs to their sizes
func getIDToSize(blocks []int) map[int]int {
	idToSizeMap := make(map[int]int)

	for _, id := range blocks {
		idToSizeMap[id]++
	}
	return idToSizeMap
}

// Moves an entire block to available free space
func moveWholeBlock(blocks []int, freeStartIndices []int, lastFileIndex int, idToSize map[int]int) []int {
	blocksCopy := append([]int{}, blocks...)

	for _, freeStartIndex := range freeStartIndices {
		fi := freeStartIndex
		var freeLengthCount int
		var currBlock = -1
		for {
			currBlock = blocksCopy[fi]
			fi++
			if currBlock != freeSpaceVal || fi == len(blocksCopy) {
				break
			}
			freeLengthCount++
		}

		lastFileSize := idToSize[blocks[lastFileIndex]]

		// Move block if sufficient free space is available
		if freeLengthCount >= lastFileSize && lastFileIndex > freeStartIndex {
			mod := 0
			for i := 0; i < lastFileSize; i++ {
				blocksCopy = swap(blocksCopy, freeStartIndex+mod, lastFileIndex-mod)
				mod++
			}
			return blocksCopy
		}
	}
	return blocksCopy
}

// Returns the index of the first free space
func getFirstFreeSpaceIndex(blocks []int) int {
	for i := range blocks {
		if blocks[i] == freeSpaceVal {
			return i
		}
	}
	return -1
}

// Returns a slice of all free space indices
func getFreeSpaceIndices(blocks []int) []int {
	freeSpaceIndices := make([]int, 0)
	for i := range blocks {
		if blocks[i] == freeSpaceVal {
			freeSpaceIndices = append(freeSpaceIndices, i)
		}
	}
	return freeSpaceIndices
}

// Returns a slice of indices of the last files in blocks
func getLastFileIndices(blocks []int) []int {
	lastFileIndices := make([]int, 0)
	var lastVal int
	for i := len(blocks) - 1; i >= 0; i-- {
		if blocks[i] == freeSpaceVal || lastVal == blocks[i] {
			continue
		}
		lastFileIndices = append(lastFileIndices, i)
		lastVal = blocks[i]
	}
	return lastFileIndices
}

// Returns the index of the last file in blocks
func getLastFileIndex(blocks []int) int {
	for i := len(blocks) - 1; i >= 0; i-- {
		if blocks[i] != freeSpaceVal {
			return i
		}
	}
	return -1
}

// Calculates a checksum based on the block indices and IDs
func calcChecksum(blocks []int) int {
	var checksum int
	for i, id := range blocks {
		if id != freeSpaceVal {
			checksum += i * id
		}
	}
	return checksum
}
