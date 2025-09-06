package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"

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

func adminPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/index.tpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Title   string
		Message string
	}{
		Title:   "Hello, World!",
		Message: "Welcome to my simple Go web server!",
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// basics
	flags := getFlags()
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", adminPage)

	// api
	http.HandleFunc("/api/ping", server.PingHandler)
	http.HandleFunc("/api/{item}/add", server.Auth(server.AddHandler))
	http.HandleFunc("/api/{item}/del", server.Auth(server.DelHandler))
	http.HandleFunc("/api/{item}/query", server.Auth(server.QueryHandler))

	// start the server
	log.Println("server is running on http://localhost" + flags.port)
	log.Fatal(http.ListenAndServe(flags.port, nil))
}
