package main

import (
	"flag"
	"log"
	"log/slog"
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
	server.Configure()

	// api
	http.HandleFunc("/ping", server.PingHandler)
	http.HandleFunc("/commit", server.Auth(server.CommitHandler))
	http.HandleFunc("/api/", server.Auth(server.ApiHandler))

	// start the server
	slog.Info("started emmer:", "port", "http://localhost"+flags.port+"/")
	log.Fatal(http.ListenAndServe(flags.port, nil))
}
