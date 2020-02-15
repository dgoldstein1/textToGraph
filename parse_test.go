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

func _mockOutCalls(success bool) {
	httpmock.Activate()
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	os.Setenv("TWO_WAY_KV_ENDPOINT", twoWayEndpoint)
	if success {
		// mock out DB call
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
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"errors": []string{"test"},
					"entries": []db.TwoWayEntry{
						db.TwoWayEntry{words[0], 1},
						db.TwoWayEntry{words[1], 1},
					},
				})
			},
		)
	} else {
		httpmock.RegisterResponder("POST", twoWayEndpoint+"/entries",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(500, map[string]interface{}{
					"errors":  []string{"test"},
					"entries": []db.TwoWayEntry{},
				})
			},
		)
	}
}

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
		Name               string
		file               *os.File
		expectedError      string
		successfulResponse bool
	}{
		{"reads all words in correctly", simpleFile, "", true},
	}
	for _, tc := range testTable {
		_mockOutCalls(tc.successfulResponse)
		defer httpmock.DeactivateAndReset()
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
