package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestApi(t *testing.T) {
	go serve()
	time.Sleep(500 * time.Millisecond)

	tests := []struct {
		endpoint       string
		body           string
		expectedStatus int
	}{
		{"/ping", "", http.StatusOK},
		{"/entry/add", `{"table":"test"}`, http.StatusOK},
		{"/table/add", `{"name":"test"}`, http.StatusInternalServerError},
		{"/entry/add", `{"table":"test","key":["a"],"value":"b"}`, http.StatusOK},
		{"/entry/query", `{"table":"test","key":["a"]}`, http.StatusOK},
		{"/table/del", `{"name":"test"}`, http.StatusOK},
	}

	for _, test := range tests {
		resp, err := http.Post("http://localhost:8080"+test.endpoint, "application/json", strings.NewReader(test.body))
		if err != nil {
			t.Fatalf("Failed to send request to %s: %v", test.endpoint, err)
		}
		defer resp.Body.Close() //nolint:errcheck

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		if resp.StatusCode != test.expectedStatus {
			t.Errorf("Expected status %d, got %d '%s'", test.expectedStatus, resp.StatusCode, respBody)
		}
	}
}
