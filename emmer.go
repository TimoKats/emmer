package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	server "github.com/TimoKats/emmer/server"
)

func getPort() string {
	port := os.Getenv("EM_PORT")
	if !server.ValidPort(port) {
		port = "8080"
	}
	return ":" + port
}

func main() {
	// basics
	port := getPort()

	// api
	server.Configure()
	http.HandleFunc("/ping", server.PingHandler)
	http.HandleFunc("/commit", server.Auth(server.CommitHandler))
	http.HandleFunc("/api/", server.Auth(server.ApiHandler))

	// start the server
	slog.Info("started emmer:", "port", "http://localhost"+port+"/api")
	log.Fatal(http.ListenAndServe(port, nil))
}
