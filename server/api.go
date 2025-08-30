package server

import (
	"encoding/json"
	"errors"

	. "github.com/TimoKats/emmer/server/fs"

	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var fs IFileSystem

// helper function to parse the body of a post request.
func parsePost(w http.ResponseWriter, r *http.Request) []byte {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return nil
	}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return body
}

// this function selects the interface based on the URL path.
func parsePathValue(value string) (IData, error) {
	switch value {
	case "table":
		return TableData{}, nil
	case "entry":
		return EntryData{}, nil
	default:
		return nil, errors.New("path " + value + " invalid")
	}
}

// does nothing. Only used for health checks.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

// used for creating tables or adding key/values to table.
func AddHandler(w http.ResponseWriter, r *http.Request) {
	// parse request
	body := parsePost(w, r)
	if body == nil {
		http.Error(w, "no body", http.StatusBadRequest)
		return
	}
	// switch paths for add (are we adding table or entry?)
	data, err := parsePathValue(r.PathValue("item"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// execute add on item
	if err := data.Add(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// used for creating tables or adding key/values to table.
func DelHandler(w http.ResponseWriter, r *http.Request) {
	// parse request
	body := parsePost(w, r)
	if body == nil {
		http.Error(w, "no body", http.StatusBadRequest)
		return
	}
	// switch paths for add (are we adding table or entry?)
	data, err := parsePathValue(r.PathValue("item"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// execute del on item
	if err := data.Del(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// used for querying tables or adding key/values to table.
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	// parse request
	w.Header().Set("Content-Type", "application/json")
	body := parsePost(w, r)
	if body == nil {
		http.Error(w, "no body", http.StatusBadRequest)
		return
	}
	// switch paths for add (are we adding table or entry?)
	data, err := parsePathValue(r.PathValue("item"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// execute and return query
	response, err := data.Query(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// upon init, selects which filesystem to use based on env variable.
func init() {
	switch os.Getenv("EMMER_FS") {
	// case "aws": < this will be the pattern
	// 	log.Println("aws not implemented yet")
	default:
		fs = LocalFS{Folder: "data"}
	}
	log.Println("selected " + fs.Info())
}
