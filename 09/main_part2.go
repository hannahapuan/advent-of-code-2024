package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// https://adventofcode.com/2024/day/9

const (
	filename string = "example.txt"
)

type block struct {
	isFile  bool
	fileID  int // only set if it is not free
	size    int
	inPlace bool
}

func main() {
	blocks, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(blocksToString(blocks))

	freeBlockIndices := getFreeBlockIndices(blocks)
	lastFileBlockIndices := getLastFileBlockIndices(blocks)
	printValues(blocks)
	fmt.Println()
	fmt.Println("freeBlockIndices: ", freeBlockIndices)
	fmt.Println("lastFileBlockIndices: \n", lastFileBlockIndices)
	afterMove := blocks
	moved := true
	fmt.Println(blocksToString(blocks))
	fmt.Println()
	var count int
	for moved && count < 5 {
		fmt.Println("move()")
		for _, lastFileBlockIndex := range lastFileBlockIndices {
			lastFileBlockIndices = getLastFileBlockIndices(afterMove)
			freeBlockIndices = getFreeBlockIndices(afterMove)
			// var imoved bool
			afterMove, _ = move(afterMove, freeBlockIndices, lastFileBlockIndex)
			fmt.Println(blocksToString(afterMove))
			fmt.Println()
			count++
		}
	}

	fmt.Println("after")
	fmt.Println(afterMove)
}

func readInput(fname string) ([]block, error) {
	blocks := make([]block, 0)

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
		var lenFile, lenFree int
		lenFile, err = strconv.Atoi(string(lenFileRu))
		if err != nil {
			return nil, fmt.Errorf("error converting file length input to int: %w", err)
		}
		fileBlock := block{
			isFile: true,
			fileID: fileIDIdx,
			size:   lenFile,
		}
		fileIDIdx++
		blocks = append(blocks, fileBlock)

		lenFreeRu, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" { // End of file
				break
			}
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		lenFree, err = strconv.Atoi(string(lenFreeRu))
		if err != nil {
			return nil, fmt.Errorf("error converting free length input to int: %w", err)
		}

		freeBlock := block{
			isFile: false,
			size:   lenFree,
		}
		blocks = append(blocks, freeBlock)

	}
	return blocks, nil
}

// returns (blocks, false) if no changes were made
func move(blocks []block, freeBlockIndices []int, lastFileBlockIndex int) ([]block, bool) {
	blocksCopy := append([]block{}, blocks...)

	fmt.Println("freeBlockIndices:", freeBlockIndices)
	fmt.Printf("lastFileBlockIndex: %v, lastFileBlockFileId: %v\n", lastFileBlockIndex, blocksCopy[lastFileBlockIndex].fileID)
	printValues(blocksCopy)
	fmt.Println()
	lastFileBlock := blocksCopy[lastFileBlockIndex]

	var moved bool
	// check if the file block has enough free space to be moved
	for _, freeBlockIndex := range freeBlockIndices {
		if blocksCopy[freeBlockIndex].size < lastFileBlock.size {
			// can't move
			fmt.Printf("couldn't move: freeBlockIndex=[%d]  lastFileBlockID=[%d]\n", freeBlockIndex, lastFileBlock.fileID)
			blocksCopy[lastFileBlockIndex].inPlace = true
			continue
		}

		// can move
		blocksCopy[lastFileBlockIndex].inPlace = true
		moved = true

		if blocksCopy[freeBlockIndex].size == lastFileBlock.size {
			fmt.Printf("moving %d\n", lastFileBlock.fileID)
			blocksCopy = swap(blocksCopy, freeBlockIndex, lastFileBlockIndex)
			break
		}

		// need to split free space
		freeSpaceLeftover := blocksCopy[freeBlockIndex].size - lastFileBlock.size
		freeSpaceToMoveToEnd := blocksCopy[freeBlockIndex].size - (blocksCopy[freeBlockIndex].size - lastFileBlock.size)

		freeSpaceLeftoverBlock := block{
			isFile: false,
			size:   freeSpaceLeftover,
		}
		freeSpaceToMoveToEndBlock := block{
			isFile: false,
			size:   freeSpaceToMoveToEnd,
		}
		fmt.Printf("moving %d\n", lastFileBlock.fileID)
		blocksCopy[freeBlockIndex] = blocksCopy[lastFileBlockIndex]
		blocksCopy = insert(blocksCopy, freeBlockIndex+1, freeSpaceLeftoverBlock)
		_, blocksCopy = remove(blocksCopy, lastFileBlockIndex+1)
		blocksCopy = append(blocksCopy, freeSpaceToMoveToEndBlock)
		break
	}

	return blocksCopy, moved
}

func getFreeBlockIndices(blocks []block) []int {
	freeBlockIndices := make([]int, 0)

	for i, block := range blocks {
		if block.isFile {
			continue
		}
		freeBlockIndices = append(freeBlockIndices, i)
	}
	return freeBlockIndices
}

func getLastFileBlockIndices(blocks []block) []int {
	lastFileBlockIndices := make([]int, 0)

	// get last fileBlock
	for i := len(blocks) - 1; i >= 0; i-- {
		if !blocks[i].isFile {
			continue
		}
		lastFileBlockIndices = append(lastFileBlockIndices, i)
	}

	return lastFileBlockIndices
}

// replaces the given index of the block slice with the input block
// returns the block no longer in the slice (that was replaced) and the new block slice
func replace(s []block, i int, b block) []block {
	_, nb0 := remove(s, i)
	nb1 := insert(nb0, i, b)
	return nb1
}

func insert(s []block, i int, b block) []block {
	s = append(s[:i+1], s[i:]...)
	s[i] = b
	return s
}

func remove(s []block, i int) (block, []block) {
	item := s[i]
	return item, append(s[:i], s[i+1:]...)
}

func blocksToString(bs []block) string {
	var export string
	for _, b := range bs {
		if b.isFile {
			for i := 0; i < b.size; i++ {
				export += fmt.Sprintf("%d", b.fileID)
			}
			continue
		}
		for i := 0; i < b.size; i++ {
			export += "."
		}
	}
	return export
}

func swap(blocks []block, a, b int) []block {
	blocks[a], blocks[b] = blocks[b], blocks[a]
	return blocks
}

func printValues(blocks []block) {
	for _, b := range blocks {
		if !b.isFile {
			fmt.Print("[ ]")
		}
		fmt.Printf("[%d]", b.fileID)
	}
}
