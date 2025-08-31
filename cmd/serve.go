package main

import (
	"log"
	"net/http"

	server "github.com/TimoKats/emmer/server"
)

func serve() {
	// basics
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// api
	http.HandleFunc("/api/ping", server.PingHandler)
	http.HandleFunc("/api/{item}/add", server.AddHandler)
	http.HandleFunc("/api/{item}/del", server.DelHandler)
	http.HandleFunc("/api/{item}/query", server.QueryHandler)

	// start the server
	log.Println("server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	log.Println("starting server...")
	serve()
}
