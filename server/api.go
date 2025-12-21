package server

import (
	"errors"
	"log/slog"

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
		http.Error(w, "use put/del/get", http.StatusMethodNotAllowed)
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
			slog.Error("error writing cache:", "file", filename)
		} else {
			slog.Debug("cache written to filesystem:", "file", filename)
			fmt.Fprint(w, "cache written to filesystem") //nolint:errcheck
		}
	}
}

// basic auth that uses public username/password for check
func Auth(next http.HandlerFunc) http.HandlerFunc {
	access := func(method string) int {
		level := session.config.access
		if method != "GET" {
			level++
		}
		slog.Debug("request auth level:", "level", level)
		return level
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if (!ok || user != session.config.username || pass != session.config.password) && access(r.Method) > 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// upon init, set credentials and filesystem to use
func init() {
	username, password := initCredentials()
	commits := initCache()
	access := initAccess()
	buffer := NewLogBuffer(20)
	log.SetOutput(io.MultiWriter(os.Stdout, buffer))
	// create config object
	session.config = Config{
		username: username,
		password: password,
		commit:   commits,
		access:   access,
	}
	// session object
	session.logBuffer = buffer
	session.cache.data = make(map[string]map[string]any)
	session.cache.data = make(map[string]map[string]any)
	session.fs = initConnector()
	session.commits = 1
}
