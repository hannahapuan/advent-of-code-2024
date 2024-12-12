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
	isFile bool
	fileID int // only set if it is not free
	size   int
}

func main() {
	blocks, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(blocksToString(blocks))

	firstFreeBlockIndex := getFirstFreeBlockIndex(blocks)
	lastFileBlockIndex := getLastFileBlockIndex(blocks)

	fmt.Println("firstFreeBlockIndex: ", firstFreeBlockIndex)
	fmt.Println("lastFileBlockIndex: ", lastFileBlockIndex)
	afterMove := blocks
	count := 0
	fmt.Println(blocksToString(blocks))
	for count < 10 {
		fmt.Println("move")
		afterMove = move(afterMove, firstFreeBlockIndex, lastFileBlockIndex)
		fmt.Println(blocksToString(afterMove))
		if afterMove == nil {
			blocks = afterMove
			break
		}
		count++
		fmt.Println()
	}

	fmt.Println("after")
	fmt.Println(blocks)
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

// // returns nil if no changes were made
// func move(blocks []block, firstFreeBlockIndex, lastFileBlockIndex int) []block {
// 	blocksCopy := append([]block{}, blocks...)

// 	// all file blocks are before the free blocks: finished
// 	if firstFreeBlockIndex > lastFileBlockIndex {
// 		fmt.Println("here")
// 		return nil
// 	}

// 	firstFreeBlock := blocks[firstFreeBlockIndex]
// 	lastFileBlock := blocks[lastFileBlockIndex]

// 	fileBlocksToMove := make([]block, 0)

// 	// check if the file block needs to be split
// 	if firstFreeBlock.size < lastFileBlock.size {
// 		fmt.Println("split case")
// 		aBlock := lastFileBlock
// 		bBlock := lastFileBlock
// 		aBlock.size = lastFileBlock.size - firstFreeBlock.size
// 		bBlock.size = lastFileBlock.size - (lastFileBlock.size - firstFreeBlock.size)
// 		fmt.Println("a size:", aBlock.size)
// 		fmt.Println("b size:", bBlock.size)
// 		fileBlocksToMove = append(fileBlocksToMove, []block{bBlock, aBlock}...)
// 	} else {
// 		fileBlocksToMove = append(fileBlocksToMove, lastFileBlock)
// 	}

// 	// replace first file block to free block location
// 	blocksCopy[firstFreeBlockIndex] = fileBlocksToMove[0]

// 	// if the block was split, update the existing last file block size
// 	if len(fileBlocksToMove) > 1 {
// 		blocksCopy[lastFileBlockIndex].size = fileBlocksToMove[1].size
// 	}

// 	// append the free space block to the end of the blocks
// 	blocksCopy = append(blocksCopy, firstFreeBlock)

// 	return blocksCopy
// }

// func getFirstFreeBlockIndex(blocks []block) int {
// 	firstFreeBlockIndex := -1

// 	// get first free block
// 	for i, block := range blocks {
// 		if block.isFile {
// 			continue
// 		}
// 		firstFreeBlockIndex = i
// 		break
// 	}
// 	return firstFreeBlockIndex
// }
// func getLastFileBlockIndex(blocks []block) int {
// 	lastFileBlockIndex := -1

// 	// get last fileBlock
// 	for i := len(blocks) - 1; i >= 0; i-- {
// 		if !blocks[i].isFile {
// 			continue
// 		}
// 		lastFileBlockIndex = i
// 		break
// 	}

// 	return lastFileBlockIndex
// }

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
