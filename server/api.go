package server

import (
	"errors"
	"log/slog"

	"fmt"
	"net/http"
)

var session Session

// helper function that selects the interface based on the URL path
func ApiHandler(w http.ResponseWriter, r *http.Request) {
	var response Response
	request, parseErr := parseRequest(r)
	item, itemErr := setItem(request)
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

// write all cache to filesystem
func CommitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "use post", http.StatusMethodNotAllowed)
		return
	}
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
	return func(w http.ResponseWriter, r *http.Request) {
		if setAccess(r.Method) > 1 {
			user, pass, ok := r.BasicAuth()
			if !ok || user != session.config.username ||
				pass != session.config.password {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

// reset cache (used for tests)
func ClearCache() {
	session.cache.data = make(map[string]any)
}

// upon init, set credentials and filesystem to use
func Configure() {
	username, password := initCredentials()
	commits := initCache()
	access := initAccess()
	// create config object
	session.config = Config{
		username: username,
		password: password,
		commit:   commits,
		access:   access,
	}
	// session object
	session.cache.data = make(map[string]any)
	session.fs = initConnector()
	session.commits = 1
}
