package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
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
	currWord, nextWord := "", ""
	for scanner.Scan() {
		nextWord = cleanWord(scanner.Text())
		// add edge
		if currWord != "" {
			if err := addEdge(currWord, nextWord); err != nil {
				return err
			}
		}
		// update words
		currWord = nextWord
	}
	return nil
}

var reg, _ = regexp.Compile("[^a-zA-Z0-9]+")

// cleans punctuation, capitalization, etc out of words
func cleanWord(w string) string {
	return strings.ToLower(reg.ReplaceAllString(w, ""))
}

// adds neccesary nodes and edges to e
func addEdge(currWord string, nextWord string) error {
	return nil
}
