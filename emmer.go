package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	server "github.com/TimoKats/emmer/server"
)

type flags struct {
	port string
}

func getFlags() flags {
	port := flag.String("p", "2112", "Port.")
	flag.Parse()
	if envPort := os.Getenv("EM_PORT"); server.ValidPort(envPort) {
		port = &envPort
	}
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
