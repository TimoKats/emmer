package main

import (
	"log"
	"net/http"

	server "github.com/TimoKats/emmer/server"
)

func serve() {
	// api
	http.HandleFunc("/ping", server.PingHandler)
	http.HandleFunc("/api/", server.ApiHandler)

	// start the server
	log.Println("server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	log.Println("starting server...")
	serve()
}
