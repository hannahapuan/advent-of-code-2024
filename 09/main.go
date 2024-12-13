package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

//
// https://adventofcode.com/2024/day/9

const (
	filename     string = "example.txt"
	freeSpaceVal int    = -1
)

func main() {
	blocks, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	moved := true
	finishedBlocks := append([]int{}, blocks...)
	for moved {
		blocks, moved = move(blocks)
	}
	fmt.Println("\npart 1 checksum:", calcChecksum(finishedBlocks))

	// Part 2
	fileEndIndices := getLastFileIndices(blocks)
	idToSize := getIDToSize(blocks)
	for _, fei := range fileEndIndices {
		fsis := getFreeSpaceIndices(blocks)
		blocks = moveWholeBlock(blocks, fsis, fei, idToSize)
	}

	fmt.Println("\npart 2 checksum:", calcChecksum(blocks))
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
		return nil, false
	}

	return swap(blocksCopy, freeIndex, fileIndex), true
}

func getIDToSize(blocks []int) map[int]int {
	idToSizeMap := make(map[int]int)

	for _, id := range blocks {
		idToSizeMap[id]++
	}
	return idToSizeMap
}

// returns updated blocks after move and if a move happened
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

func getFirstFreeSpaceIndex(blocks []int) int {
	for i := range blocks {
		if blocks[i] != freeSpaceVal {
			continue
		}
		return i
	}
	return -1
}

func getFreeSpaceIndices(blocks []int) []int {
	freeSpaceIndices := make([]int, 0)
	for i := range blocks {
		if blocks[i] != freeSpaceVal {
			continue
		}
		freeSpaceIndices = append(freeSpaceIndices, i)
	}
	return freeSpaceIndices
}

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
