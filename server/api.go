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

// get HTTP request and format it into Request object used by server
func parseRequest(r *http.Request) (Request, error) {
	// parse URL path (parameters, path)
	request := Request{Method: r.Method, Mode: r.FormValue("mode")}
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
		return request, err
	}
	if len(payload) > 0 {
		err = json.Unmarshal(payload, &request.Value)
	}
	return request, err
}

// takes response object and writes the HTTP response object
func parseResponse(w http.ResponseWriter, response Response) error {
	if response.Error != nil {
		w.Header().Set("Content-Type", "text/plain")
		if strings.Contains(response.Error.Error(), "not found") {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(500)
		}
		return json.NewEncoder(w).Encode(response.Error.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(response.Data)
}

// returns the item to apply CRUD operations on
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
	// set up
	var response Response
	request, parseErr := parseRequest(r)
	item, itemErr := selectItem(request)
	if err := errors.Join(parseErr, itemErr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // to response
		return
	}
	// select function
	switch request.Method {
	case "PUT":
		response = item.Add(request)
	case "DELETE":
		response = item.Del(request)
	case "GET":
		response = item.Get(request)
	default:
		http.Error(w, "please use put/del/get", http.StatusMethodNotAllowed) // to response
		return
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

// shows last n (20) logs from server
func LogsHandler(w http.ResponseWriter, r *http.Request) {
	for _, entry := range config.logBuffer.GetLogs() {
		fmt.Fprint(w, entry) //nolint:errcheck
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
	// auth settings
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
	// logs settings
	buffer := NewLogBuffer(20)
	log.SetOutput(io.MultiWriter(os.Stdout, buffer))
	// create config object
	config = Config{
		logBuffer: buffer,
		autoTable: os.Getenv("EM_AUTOTABLE") != "false",
		username:  username,
		password:  password,
		fs:        emmerFs.SetupLocal(),
	}
}
