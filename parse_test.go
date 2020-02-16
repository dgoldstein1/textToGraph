package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgoldstein1/crawler/db"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var dbEndpoint = "http://localhost:17474"
var twoWayEndpoint = "http://localhost:17475"

func _mockOutCalls() {
	httpmock.Reset()
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	os.Setenv("TWO_WAY_KV_ENDPOINT", twoWayEndpoint)

	httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{"neighborsAdded": []int{2, 3, 4}})
		},
	)

	// mock out metadata call
	httpmock.RegisterResponder("POST", twoWayEndpoint+"/entries",
		func(req *http.Request) (*http.Response, error) {
			defer req.Body.Close()
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			words := []string{}
			err = json.Unmarshal(body, &words)
			if err != nil {
				panic(err)
			}
			errors := []string{}
			entries := []db.TwoWayEntry{}
			for _, w := range words {
				entries = append(entries, db.TwoWayEntry{w, 1})
			}
			if words[0] == "badresponse" {
				entries = []db.TwoWayEntry{}
			}

			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"errors":  errors,
				"entries": entries,
			})
		},
	)
}

func TestParse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
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
		_mockOutCalls()
		t.Run(tc.name, func(t *testing.T) {
			errors = []string{}
			Parse(tc.filePath)
			assert.Equal(t, tc.expectedNumberOfErrors, len(errors))
		})
	}
}

func TestIndexWords(t *testing.T) {
	origLogErrorf := logErrorf
	defer func() { logErrorf = origLogErrorf }()
	errors := []string{}
	logErrorf = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	simpleFile, _ := os.Open("./data/simple.txt")
	simpleLongFile, _ := os.Open("./data/simpleFileLong.txt")
	badResponse, _ := os.Open("./data/badResponse.txt")

	testTable := []struct {
		Name             string
		file             *os.File
		lenexpectedError int
	}{
		{"reads all words in correctly", simpleFile, 0},
		{"dumps local cache if reaches maxLenBeforeDump", simpleLongFile, 0},
		{"unsuccessful response from back end", badResponse, 1},
	}
	for _, tc := range testTable {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Name == "dumps local cache if reaches maxLenBeforeDump" {
				maxLenBeforeDump = 5
			}
			_mockOutCalls()
			errors = []string{}
			indexWords(tc.file)
			assert.Equal(t, tc.lenexpectedError, len(errors))

		})
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
