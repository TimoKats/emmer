package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// read json file at test location
func readJson(filename string) any {
	data := make(map[string]any)
	list := []any{}
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(file, &data); err != nil {
		json.Unmarshal(file, &list) //nolint:errcheck
		return list
	}
	return data
}

// read csv file at testCsv location into list
func readCsv(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

// test if two maps are equal (loosly/stringified)
func jsonEqual(a, b any) bool {
	aByte, _ := json.Marshal(a) //nolint:errcheck
	bByte, _ := json.Marshal(b) //nolint:errcheck
	return string(aByte) == string(bByte)
}

// send http request and returns the status code and error
func request(method string, endpoint string, body string) int {
	url := "http://localhost:8080" + endpoint
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return 0
	}
	client := &http.Client{}
	resp, _ := client.Do(req) //nolint:errcheck
	return resp.StatusCode
}

// create path to test file (GitHub actions uses EM_FOLDER) and delete current.
func testFile(ext string) string {
	var folder string
	if folder = os.Getenv("EM_FOLDER"); folder == "" {
		folder = filepath.Join(os.Getenv("HOME"), ".local", "share", "emmer")
	}
	path := filepath.Join(folder, "test."+ext)
	if err := os.Remove(path); err != nil {
		log.Println("Error deleting test file:", err)
	}
	return path
}
