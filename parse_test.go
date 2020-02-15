package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	// mock out log.Fatalf
	origLogFatalf := logFatalf
	defer func() { logFatalf = origLogFatalf }()
	errors := []string{}
	logFatalf = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}
	testTable := []struct {
		name                   string
		filePath               string
		expectedNumberOfErrors int
	}{
		{"file does not exist", "./sdfsdf.txt", 1},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			errors = []string{}
			Parse(tc.filePath)
			assert.Equal(t, tc.expectedNumberOfErrors, len(errors))
		})
	}
}

func TestIndexWords(t *testing.T) {
	simpleFile, err := os.Open("./data/simple.txt")
	if err != nil {
		t.Errorf("Could not open simple file: %v", err)
	}

	testTable := []struct {
		Name          string
		file          *os.File
		expectedError string
	}{
		{"reads all words in correctly", simpleFile, ""},
	}
	for _, tc := range testTable {
		err := indexWords(tc.file)
		if err != nil {
			assert.Equal(t, tc.expectedError, err.Error())
		} else {
			assert.Equal(t, tc.expectedError, "")
		}
	}
}

func TestCleanWord(t *testing.T) {
	for _, tc := range []struct {
		name           string
		word           string
		expectedOutput string
	}{
		{"returns lowercase", "HELPME", "helpme"},
		{"removes all non [a-b][0-9] chars", ".'42fjae2h", "42fjae2h"},
		{"removes all non [a-b][0-9] chars (2)", ".'4@%*$)(@(#) ....2fjae2h", "42fjae2h"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedOutput, cleanWord(tc.word))
		})
	}
}
