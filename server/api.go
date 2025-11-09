package server

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strconv"

	emmerFs "github.com/TimoKats/emmer/server/fs"

	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var session Session

// helper function that selects the interface based on the URL path
func ApiHandler(w http.ResponseWriter, r *http.Request) {

	// returns the item to apply CRUD operations on
	toggle := func(request Request) (Item, error) {
		if len(request.Key) > 0 {
			return EntryItem{}, nil
		}
		return TableItem{}, nil
	}

	// apply CRUD to correct item and return response
	var response Response
	request, parseErr := parseRequest(r)
	item, itemErr := toggle(request)
	if err := errors.Join(parseErr, itemErr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // to response
		return
	}
	switch request.Method {
	case "PUT":
		response = item.Add(request)
	case "DELETE":
		response = item.Del(request)
	case "GET":
		response = item.Get(request)
	default:
		http.Error(w, "please use put/del/get", http.StatusMethodNotAllowed)
		return
	}
	if err := parseResponse(w, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// does nothing, only used for health checks
func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong") //nolint:errcheck
}

// shows last n (20) logs from server
func LogsHandler(w http.ResponseWriter, r *http.Request) {
	for _, entry := range session.logBuffer.GetLogs() {
		fmt.Fprint(w, entry) //nolint:errcheck
	}
}

// write all cache to filesystem
func CommitHandler(w http.ResponseWriter, r *http.Request) {
	for filename, data := range session.cache.data {
		err := session.fs.Put(filename, data)
		if err != nil {
			log.Printf("error writing cache of %s", filename)
		} else {
			fmt.Fprint(w, "cache written to filesystem") //nolint:errcheck
		}
	}
}

// writes backup cache to filesystem, and clears backup
func UndoHandler(w http.ResponseWriter, r *http.Request) {
	if session.cache.backup.table != "" {
		write(session.cache.backup.table, session.cache.backup.value)
		session.cache.backup = Backup{}
		return
	}
	http.Error(w, "no undo available", http.StatusBadRequest)
}

// basic auth that uses public username/password for check
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != session.config.username || pass != session.config.password {
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
	// cache settings
	commit := 1
	commitEnv := os.Getenv("EM_COMMIT")
	if commitEnv != "" {
		commitInt, err := strconv.Atoi(commitEnv)
		if err != nil {
			fmt.Printf("Error converting commit strategy to int: %v", err)
			return
		}
		commit = commitInt
	}
	// logs settings
	buffer := NewLogBuffer(20)
	log.SetOutput(io.MultiWriter(os.Stdout, buffer))
	// create config object
	session.config = Config{
		username: username,
		password: password,
		commit:   commit,
	}
	session.logBuffer = buffer
	session.cache.data = make(map[string]map[string]any)
	session.cache.data = make(map[string]map[string]any)
	session.fs = emmerFs.SetupLocal()
	session.commits = 1
}
