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
		port = "2112"
	}
	return ":" + port
}

func getVersion() string {
	version := os.Getenv("EM_VERSION")
	if len(version) == 0 {
		return "unknown"
	}
	return version
}

func main() {
	// basics
	port := getPort()
	version := getVersion()
	url := "http://localhost" + port + "/api"

	// api
	server.Configure()
	http.HandleFunc("/ping", server.PingHandler)
	http.HandleFunc("/commit", server.Auth(server.CommitHandler))
	http.HandleFunc("/api/", server.Auth(server.ApiHandler))

	// start the server
	slog.Info("started emmer:", "url", url, "version", version)
	log.Fatal(http.ListenAndServe(port, nil))
}
