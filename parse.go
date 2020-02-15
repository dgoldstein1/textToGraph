package main

import (
	"bufio"
	"fmt"
	"os"
)

func Parse(filePath string) {
	logMsg("reading in file %s", filePath)
	// open file
	file, err := os.Open(filePath)
	if err != nil {
		logFatalf("Could not read in file %v", err)
	}
	// add edge for each new word
	indexWords(file)
}

// scan each word in file, adding edge for word -> word[i+1]
func indexWords(file *os.File) error {
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}

// adds neccesary nodes and edges to e
func addEdge(currWord string, nextWord string) error {
	return nil
}
