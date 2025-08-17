package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Helper function to parse the body of a post request.
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

// Does nothing. Only used for health checks.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

// Used for creating tables or adding key/values to table.
func AddHandler(path Path) http.Handler {
	var body []byte
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if body = parsePost(w, r); body == nil {
			return
		}
		if err := add(body, path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// Used for removing tables or key/values from tables.
func DelHandler(path Path) http.Handler {
	var body []byte
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if body = parsePost(w, r); body == nil {
			return
		}
		if err := del(body, path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// Query handler is used to filter/fetch data from jsons.
func QueryHandler(path Path) http.Handler {
	var body []byte
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if body = parsePost(w, r); body == nil {
			return
		}
		response, err := query(body, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
