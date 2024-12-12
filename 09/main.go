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

type fileBlock struct {
	id  int
	val int
}

func main() {
	blocks, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(blocksToString(blocks))
}

func readInput(fname string) ([]int, error) {
	// blocks := make([]block, 0)
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

		lenFree, err = strconv.Atoi(string(lenFreeRu))
		if err != nil {
			return nil, fmt.Errorf("error converting free length input to int: %w", err)
		}

		for i := 0; i < lenFree; i++ {
			ids = append(ids, -1)
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
