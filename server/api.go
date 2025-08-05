package server

import (
	"encoding/json"
	"io"
	"net/http"
)

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

func PingHandler(w http.ResponseWriter, r *http.Request) {
	return
}

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

func QueryHandler(path Path) http.Handler {
	var body []byte
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if body = parsePost(w, r); body == nil {
			return
		}
		result, err := query(body, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "failed to encode result", http.StatusInternalServerError)
			return
		}
	})
}
