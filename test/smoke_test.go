package main

import (
	server "github.com/TimoKats/emmer/server"

	"bytes"
	"log"
	"net/http"
	"testing"
	"time"
)

func serve() {
	http.HandleFunc("/ping", server.PingHandler)
	http.HandleFunc("/api/", server.ApiHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// send http request and returns the status code and error
func request(cfg RequestConfig) int {
	url := "http://localhost:8080" + cfg.Endpoint
	req, err := http.NewRequest(cfg.Method, url, bytes.NewBuffer([]byte(cfg.Body)))
	if err != nil {
		return 0
	}
	client := &http.Client{}
	resp, _ := client.Do(req) //nolint:errcheck
	return resp.StatusCode
}

func run(tests []RequestConfig, results []map[string]any, t *testing.T) {
	filename := filename()
	for index, test := range tests {
		StatusCode := request(test)
		if StatusCode != test.ExpectedStatus {
			t.Errorf("Expected status %d, got %d", test.ExpectedStatus, StatusCode)
		}
		if results[index] != nil {
			result := readJson(filename)
			if !jsonEqual(results[index], result) {
				t.Errorf("Expected result %s, got %s", results[index], result)
			}
		}
	}
}

func TestApi(t *testing.T) {
	go serve()
	time.Sleep(500 * time.Millisecond)
	// test scenario's + expected results
	tests := []RequestConfig{
		{"PUT", "/api/test", `{"timo":1}`, http.StatusOK},
		{"PUT", "/api/test/pipo", `5`, http.StatusOK},
		{"DELETE", "/api/test", `{}`, http.StatusOK},
	}
	results := []map[string]any{
		{"timo": 1},
		{"timo": 1, "pipo": 5},
		nil,
	}
	// run tests
	run(tests, results, t)
}
