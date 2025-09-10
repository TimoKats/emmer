package main

import (
	"log"
	"net/http"

	server "github.com/TimoKats/emmer/server"
)

func serve() {
	// api
	http.HandleFunc("/ping", server.PingHandler)
	http.HandleFunc("/{item}/add", server.AddHandler)
	http.HandleFunc("/{item}/del", server.DelHandler)
	http.HandleFunc("/{item}/query", server.QueryHandler)

	// start the server
	log.Println("server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	log.Println("starting server...")
	serve()
}
