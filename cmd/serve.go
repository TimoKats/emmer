package main

import (
	"log"
	"net/http"

	server "github.com/TimoKats/emmer/server"
)

func main() {
	// basics
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/api/ping", server.PingHandler)

	// adding
	http.Handle("/api/table/add", server.AddHandler(server.Table))
	http.Handle("/api/entry/add", server.AddHandler(server.Entry))

	// deleting
	http.Handle("/api/table/del", server.DelHandler(server.Table))
	http.Handle("/api/entry/del", server.DelHandler(server.Entry))

	// start the server
	log.Println("server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
