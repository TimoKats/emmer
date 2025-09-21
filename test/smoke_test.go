package main

import (
	server "github.com/TimoKats/emmer/server"

	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

type RequestConfig struct {
	Method         string
	Endpoint       string
	Body           string
	ExpectedStatus int
}

func serve() {
	// api
	http.HandleFunc("/logs", server.LogsHandler)
	http.HandleFunc("/ping", server.PingHandler)
	http.HandleFunc("/api/", server.ApiHandler)

	// start the server
	log.Println("server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// send http request and returns the status code and error
func request(cfg RequestConfig) int {
	// format request from object
	req, err := http.NewRequest(cfg.Method, "http://localhost:8080"+cfg.Endpoint, bytes.NewBuffer([]byte(cfg.Body)))
	if err != nil {
		log.Println(err)
		return 0
	}
	// set headers / send request
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer resp.Body.Close() //nolint:errcheck
	return resp.StatusCode
}

// read json file at test location
func readJson(filename string) map[string]any {
	// get raw data
	data := make(map[string]any)
	file, err := os.ReadFile(filename)
	if err != nil {
		return data
	}
	// put raw data into map object
	json.Unmarshal(file, &data) //nolint:errcheck
	return data
}

// test if two maps are equal (loosly/stringified)
func jsonEqual(a, b any) bool {
	aByte, _ := json.Marshal(a) //nolint:errcheck
	bByte, _ := json.Marshal(b) //nolint:errcheck
	return string(aByte) == string(bByte)
}

// get location of test file
func getTestFile() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Panic("can't setup emmer folder")
	}
	return dirname + "/emmer/test.json"
}

func TestApi(t *testing.T) {
	// start server
	go serve()
	time.Sleep(500 * time.Millisecond)
	// setup tests, files, etc
	testFile := getTestFile()
	tests := []RequestConfig{
		{"PUT", "/api/test", `{"timo":1}`, http.StatusOK},
		{"PUT", "/api/test/pipo", `5`, http.StatusOK},
		{"DELETE", "/api/test", `{}`, http.StatusOK},
	}
	expectedData := []map[string]any{
		{"timo": 1},
		{"timo": 1, "pipo": 5},
		nil,
	}
	// run tests
	for index, test := range tests {
		StatusCode := request(test)
		if StatusCode != test.ExpectedStatus {
			t.Errorf("Expected status %d, got %d", test.ExpectedStatus, StatusCode)
		}
		if expectedData[index] != nil {
			result := readJson(testFile)
			if !jsonEqual(expectedData[index], result) {
				t.Errorf("Expected result %s, got %s", expectedData[index], result)
			}
		}
	}
}
