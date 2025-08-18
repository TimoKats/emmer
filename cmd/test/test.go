package main

import (
	"fmt"
	"net/http"
)

func request(path string) int {
	resp, err := http.Get("http://localhost:8080" + path)
	if err != nil {
		return 500
	}
	return resp.StatusCode
}

func main() {
	fmt.Println("starting tests\n---")
	paths := []string{"/api/ping", "/api/table/add", "/api/entry/get"}
	for _, path := range paths {
		fmt.Printf("%d\t%s\n", request(path), path)
	}
	fmt.Println("---\nfinished tests")
}
