package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
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

	f, err := os.Create("profile.pf")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// Part 1: Move blocks and calculate checksum
	moved := true
	finishedBlocks := append([]int{}, blocks...) // Copy of blocks for manipulation
	for moved {
		// Move blocks until no movement occurs
		finishedBlocks, moved = move(finishedBlocks)
		// Create a new copy to avoid mutating the original
		finishedBlocks = append([]int{}, finishedBlocks...)
	}
	fmt.Println("part 1 checksum:", calcChecksum(finishedBlocks)) // Calculate and print checksum for Part 1

	// Part 2: Rearrange blocks and calculate checksum
	idToSize := getIDToSize(blocks)              // Map block IDs to their sizes
	fileEndIndices := getLastFileIndices(blocks) // Get indices of the last files
	fsis := getFreeSpaceIndices(blocks)          // Get indices of free space

	// Move entire blocks based on available free space
	for _, fei := range fileEndIndices {
		blocks = moveWholeBlock(blocks, fsis, fei, idToSize) // Move the block
	}

	// Calculate and print checksum for Part 2
	fmt.Println("part 2 checksum:", calcChecksum(blocks))
}

// Reads the input file and parses it into a slice of blocks
func readInput(fname string) ([]int, error) {
	ids := make([]int, 0) // Slice to store block IDs

	// Open the input file
	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file:  %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file) // Reader for efficient file reading
	var fileIDIdx int               // Counter for file IDs

	for {
		// Read one character (file length) at a time
		lenFileRu, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" { // End of file
				break
			}
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		if lenFileRu == '\n' { // End of a line
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
		if lenFreeRu == '\n' { // End of a line
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

// Converts blocks to a string representation for debugging
func blocksToString(bs []int) string {
	var export string
	for _, b := range bs {
		if b == freeSpaceVal {
			export += "." // Represent free space with a dot
			continue
		}
		export += fmt.Sprintf("%d", b) // Represent file block with its ID
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
	blocksCopy := append([]int{}, blocks...) // Create a copy to avoid modifying input

	// Find indices of the first free space and the last file
	freeStartIndex := getFirstFreeSpaceIndex(blocksCopy)
	fileIndex := getLastFileIndex(blocks)

	// If all free spaces are after files, no movement needed
	if fileIndex < freeStartIndex {
		return blocks, false
	}

	// Swap the free space and file
	return swap(blocksCopy, freeStartIndex, fileIndex), true
}

// Creates a map of block IDs to their sizes
func getIDToSize(blocks []int) map[int]int {
	idToSizeMap := make(map[int]int)

	for _, id := range blocks {
		idToSizeMap[id]++ // Increment size for each block ID
	}
	return idToSizeMap
}

// Calculates the length of consecutive free blocks starting at a given index
func getFreeBlockLengths(blocks []int) []int {
	freeLengthCounts := make([]int, 0)
	var freeLengthCount int

	for _, block := range blocks {
		if block != freeSpaceVal {
			freeLengthCount = 0
			continue
		}
		freeLengthCount++
		freeLengthCounts = append(freeLengthCounts, freeLengthCount)
	}
	return freeLengthCounts
}

// Moves an entire block to available free space
func moveWholeBlock(blocks []int, freeStartIndices []int, lastFileIndex int, idToSize map[int]int) []int {
	blocksCopy := append([]int{}, blocks...) // Create a copy of the blocks

	freeLengthCounts := getFreeBlockLengths(blocksCopy)
	for i, freeStartIndex := range freeStartIndices {
		// Get free space length
		lastFileSize := idToSize[blocks[lastFileIndex]] // Get size of the last file block

		// Move block if there is enough free space and it's valid to move
		if freeLengthCounts[i] >= lastFileSize && lastFileIndex > freeStartIndex {
			blocksCopy = moveBlock(blocks, lastFileSize, lastFileIndex, freeStartIndex)
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
	return -1 // No free space found
}

// Moves a block from one position to another
func moveBlock(blocks []int, lastFileSize, lastFileIndex, freeStartIndex int) []int {
	for i := 0; i < lastFileSize; i++ {
		blocks = swap(blocks, freeStartIndex+i, lastFileIndex-i)
	}
	return blocks
}

// Returns a slice of indices of all free spaces
func getFreeSpaceIndices(blocks []int) []int {
	freeSpaceIndices := make([]int, 0)
	for i := range blocks {
		if blocks[i] == freeSpaceVal {
			freeSpaceIndices = append(freeSpaceIndices, i)
		}
	}
	return freeSpaceIndices
}

// Returns a slice of indices of the last files in the blocks
func getLastFileIndices(blocks []int) []int {
	lastFileIndices := make([]int, 0)
	var lastVal int
	for i := len(blocks) - 1; i >= 0; i-- {
		if blocks[i] != freeSpaceVal && blocks[i] != lastVal {
			lastFileIndices = append(lastFileIndices, i)
			lastVal = blocks[i]
		}
	}
	return lastFileIndices
}

// Returns the index of the last file in the blocks
func getLastFileIndex(blocks []int) int {
	for i := len(blocks) - 1; i >= 0; i-- {
		if blocks[i] != freeSpaceVal {
			return i
		}
	}
	return -1 // No file found
}

// Calculates a checksum based on block indices and IDs
func calcChecksum(blocks []int) int {
	var checksum int
	for i, id := range blocks {
		if id != freeSpaceVal {
			checksum += i * id // Multiply index and ID for checksum
		}
	}
	return checksum
}
