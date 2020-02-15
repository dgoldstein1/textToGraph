package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
			assert.Equal(t, len(errors), tc.expectedNumberOfErrors)
		})
	}
}
