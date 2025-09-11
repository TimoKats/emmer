package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	emmerFs "github.com/TimoKats/emmer/server/fs"

	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var config Config

func parseRequest(w http.ResponseWriter, r *http.Request) Request {
	// parse URL path
	request := Request{Method: r.Method}
	urlPath := r.URL.Path[len("/api/"):]
	urlItems := strings.Split(urlPath, "/")
	if len(urlItems) > 0 {
		request.Table = urlItems[0]
		if len(urlItems) > 1 {
			request.Key = urlItems[1:]
		}
	}
	// parse request body
	payload, err := io.ReadAll(r.Body)
	defer r.Body.Close() //nolint:errcheck
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := json.Unmarshal(payload, &request.Value); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return request
}

func parseResponse(w http.ResponseWriter, response Response) error {
	if response.Error != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(500)
		return json.NewEncoder(w).Encode(response.Error.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(response.Data)
}

func selectItem(request Request) (Item, error) {
	if len(request.Key) > 0 {
		return EntryItem{}, nil
	}
	if len(request.Table) > 0 {
		return TableItem{}, nil
	}
	return nil, errors.New("no table / key provided")
}

// helper function that selects the interface based on the URL path
func ApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response Response
	request := parseRequest(w, r)
	item, _ := selectItem(request)
	// select function
	switch request.Method {
	case "PUT":
		response = item.Add(request)
	case "DELETE":
		response = item.Del(request)
	case "GET":
		response = item.Query(request)
	default:
		http.Error(w, "please use put/del/get", http.StatusMethodNotAllowed)
	}
	// check errors and return response
	if err := parseResponse(w, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// does nothing. Only used for health checks
func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong") //nolint:errcheck
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
		autoTable: os.Getenv("EM_AUTOTABLE") != "false",
		username:  username,
		password:  password,
		fs:        emmerFs.SetupLocal(),
	}
}
