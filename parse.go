package main

import (
	"bufio"
	"github.com/dgoldstein1/crawler/db"
	"os"
	"regexp"
	"strings"
)

var nodeToNeighbors = map[string][]string{}
var maxLenBeforeDump = 100000

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
// logs errors on failures
func indexWords(file *os.File) {
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	currWord, nextWord := "", ""
	for scanner.Scan() {
		nextWord = cleanWord(scanner.Text())
		// add edge if there is one
		if currWord != "" {
			// add to big map
			nodeToNeighbors[currWord] = append(nodeToNeighbors[currWord], nextWord)
		}
		// update words
		currWord = nextWord
		// dump map if too large
		if len(nodeToNeighbors) > maxLenBeforeDump {
			dumpMap()
		}
	}
	// dump map at end to remove remaining
	dumpMap()
}

// dumps all nodeToNeighbors into back end
func dumpMap() {
	// add each edge for each neighbor
	for n, neighbors := range nodeToNeighbors {
		if err := addEdge(n, neighbors); err != nil {
			logErrorf("could not add edge %s to %v: %v", n, neighbors, err)
		}
	}
	// reset in local memory
	nodeToNeighbors = map[string][]string{}
}

var reg, _ = regexp.Compile("[^a-zA-Z0-9]+")

// cleans punctuation, capitalization, etc out of words
func cleanWord(w string) string {
	return strings.ToLower(reg.ReplaceAllString(w, ""))
}

// adds neccesary nodes and edges to e
func addEdge(currWord string, neighbors []string) error {
	_, err := db.AddEdgesIfDoNotExist(
		currWord,
		neighbors,
		cleanWord,
		"",
	)
	return err
}
