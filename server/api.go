package server

import (
	"encoding/json"
	"errors"

	emmerFs "github.com/TimoKats/emmer/server/fs"

	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var fs emmerFs.FileSystem

// helper function to parse the payload of a post request.
func parsePost(w http.ResponseWriter, r *http.Request) []byte {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return nil
	}
	payload, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return payload
}

// this function selects the interface based on the URL path.
func parsePathValue(value string) (Item, error) {
	switch value {
	case "table":
		return TableItem{}, nil
	case "entry":
		return EntryItem{}, nil
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
	payload := parsePost(w, r)
	if payload == nil {
		http.Error(w, "no payload", http.StatusBadRequest)
		return
	}
	// switch paths for add (are we adding table or entry?)
	data, err := parsePathValue(r.PathValue("item"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// execute add on item
	if err := data.Add(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// used for creating tables or adding key/values to table.
func DelHandler(w http.ResponseWriter, r *http.Request) {
	// parse request
	payload := parsePost(w, r)
	if payload == nil {
		http.Error(w, "no payload", http.StatusBadRequest)
		return
	}
	// switch paths for add (are we adding table or entry?)
	data, err := parsePathValue(r.PathValue("item"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// execute del on item
	if err := data.Del(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// used for querying tables or adding key/values to table.
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	// parse request
	w.Header().Set("Content-Type", "application/json")
	payload := parsePost(w, r)
	if payload == nil {
		http.Error(w, "no payload", http.StatusBadRequest)
		return
	}
	// switch paths for add (are we adding table or entry?)
	data, err := parsePathValue(r.PathValue("item"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// execute and return query
	response, err := data.Query(payload)
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
	switch os.Getenv("EM_FILESYSTEM") {
	// case "aws": < this will be the pattern
	// 	log.Println("aws not implemented yet")
	default:
		fs = emmerFs.SetupLocal()
	}
	log.Println("selected " + fs.Info())
}
