package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"

	emmerFs "github.com/TimoKats/emmer/server/fs"

	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var config Config

// helper function to parse the payload of a post request
func parsePost(w http.ResponseWriter, r *http.Request) []byte {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return nil
	}
	payload, err := io.ReadAll(r.Body)
	defer r.Body.Close() //nolint:errcheck
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return payload
}

// this function selects the interface based on the URL path
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

// does nothing. Only used for health checks
func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong") //nolint:errcheck
}

// used for creating tables or adding key/values to table
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

// used for querying tables or adding key/values to table
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

// basic auth that uses public username/password for check
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != config.username || pass != config.password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// upon init, set credentials and filesystem to use
func init() {
	username := os.Getenv("EM_USERNAME")
	if username == "" {
		username = "admin"
		log.Printf("set username to: %s", username)
	}
	password := os.Getenv("EM_PASSWORD")
	if password == "" {
		b := make([]byte, 12)
		rand.Read(b) //nolint:errcheck
		password = base64.URLEncoding.EncodeToString(b)
		log.Printf("set password to: %s", password)
	}
	config = Config{
		autoTable: !(os.Getenv("EM_AUTOTABLE") == "false"),
		username:  username,
		password:  password,
		fs:        emmerFs.SetupLocal(),
	}
}
