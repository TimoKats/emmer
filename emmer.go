package main

import (
	"flag"
	"log"
	"net/http"

	server "github.com/TimoKats/emmer/server"
)

type flags struct {
	port string
}

func getFlags() flags {
	port := flag.String("p", "8080", "Port.")
	flag.Parse()
	return flags{port: ":" + *port}
}

func main() {
	// basics
	flags := getFlags()
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// api
	http.HandleFunc("/api/ping", server.PingHandler)
	http.HandleFunc("/api/{item}/add", server.Auth(server.AddHandler))
	http.HandleFunc("/api/{item}/del", server.Auth(server.DelHandler))
	http.HandleFunc("/api/{item}/query", server.Auth(server.QueryHandler))

	// start the server
	log.Println("server is running on http://localhost" + flags.port)
	log.Fatal(http.ListenAndServe(flags.port, nil))
}
