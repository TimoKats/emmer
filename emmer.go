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

	// api
	http.HandleFunc("/ping", server.PingHandler)
	http.HandleFunc("/logs", server.Auth(server.LogsHandler))
	http.HandleFunc("/commit", server.CommitHandler)
	http.HandleFunc("/api/", server.ApiHandler)

	// start the server
	log.Println("server is running on http://localhost" + flags.port)
	log.Fatal(http.ListenAndServe(flags.port, nil))
}
