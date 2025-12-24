package main

import (
	server "github.com/TimoKats/emmer/server"

	"log"
	"net/http"
	"testing"
)

func serve() {
	server.Configure()
	http.HandleFunc("/ping", server.PingHandler)
	http.HandleFunc("/api/", server.ApiHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func TestBasicPut(t *testing.T) {
	server.ClearCache()
	file := testFile()
	request("PUT", "/api/test", `{"foo":1}`)
	expected := map[string]any{"foo": 1}
	result := readJson(file)
	if !jsonEqual(result, expected) {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestNestedPut(t *testing.T) {
	server.ClearCache()
	file := testFile()
	request("PUT", "/api/test", `{"foo":1}`)
	request("PUT", "/api/test/1", `2`)
	expected := map[string]any{"foo": 1, "1": 2}
	result := readJson(file)
	if !jsonEqual(result, expected) {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestAppend(t *testing.T) {
	server.ClearCache()
	file := testFile()
	request("PUT", "/api/test", `{"list":1}`)
	request("PUT", "/api/test/list?mode=append", `2`)
	expected := map[string]any{"list": []int{1, 2}}
	result := readJson(file)
	if !jsonEqual(result, expected) {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestIncrement(t *testing.T) {
	server.ClearCache()
	file := testFile()
	request("PUT", "/api/test", `{"list":1}`)
	request("PUT", "/api/test/list?mode=increment", `2`)
	expected := map[string]any{"list": 3}
	result := readJson(file)
	if !jsonEqual(result, expected) {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestDelete(t *testing.T) {
	server.ClearCache()
	file := testFile()
	request("PUT", "/api/test", `{"foo":"test"}`)
	request("PUT", "/api/test/bar/something", `"else"`)
	result1 := readJson(file)
	request("DELETE", "/api/test/bar", ``)
	expected1 := map[string]any{"foo": "test", "bar": map[string]any{"something": "else"}}
	expected2 := map[string]any{"foo": "test"}
	result2 := readJson(file)
	if !jsonEqual(result1, expected1) || !jsonEqual(result2, expected2) {
		t.Errorf("Failed comparison when deleting recently added data.")
	}
}

func TestCache(t *testing.T) {
	t.Setenv("EM_COMMITS", "2")
	server.Configure()
	file := testFile()
	request("PUT", "/api/test", `{"foo":"test"}`)
	result1 := readJson(file)
	request("PUT", "/api/test/bar/something", `"else"`)
	result2 := readJson(file)
	expected := map[string]any{"foo": "test", "bar": map[string]any{"something": "else"}}
	if result1 != nil || !jsonEqual(result2, expected) {
		t.Errorf("Failed comparison when deleting recently added data.")
	}
}

func init() {
	go serve()
}
