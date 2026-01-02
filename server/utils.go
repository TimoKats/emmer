package server

import (
	"log/slog"

	emmerFs "github.com/TimoKats/emmer/server/fs"

	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// used to check EM_PORT value
func ValidPort(s string) bool {
	port, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return port >= 1 && port <= 65535
}

// returns true if the error message implies a bad request error. Else 500.
func errorBadRequest(message string) bool {
	indicators := []string{"not found", "404", "invalid path", "invalid index"}
	for _, indicator := range indicators {
		if strings.Contains(message, indicator) {
			return true
		}
	}
	return false
}

// get HTTP request and format it into Request object used by server
func parseRequest(r *http.Request) (Request, error) {
	request := Request{Method: r.Method, Mode: r.FormValue("mode")}
	urlPath := r.URL.Path[len("/api/"):]
	urlItems := strings.Split(urlPath, "/")
	if len(urlItems) > 0 {
		request.Table = urlItems[0]
		if len(urlItems) > 1 {
			request.Key = urlItems[1:]
		}
	}
	if strings.Contains(urlPath, "..") {
		return request, errors.New("invalid path: parent directory")
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
		if errorBadRequest(response.Error.Error()) {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return json.NewEncoder(w).Encode(response.Error.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(response.Data)
}

// tries reading data from cache, reads from filesystem as backup
func read(filename string, mode string) (any, error) {
	if data, ok := session.cache.data[filename]; ok && mode != "fs" {
		slog.Debug("reading data from cache")
		return data, nil
	}
	slog.Debug("reading data from filesystem")
	data, err := session.fs.Get(filename) // NOTE: returns only {}
	if err == nil {
		session.cache.data[filename] = data
	}
	return data, err
}

// write to cache, and potentially to filesystem (depending on commit strategy)
func write(table string, data any) error {
	session.cache.data[table] = data
	if session.config.commit == session.commits {
		slog.Debug("writing to filesystem")
		err := session.fs.Put(table, data)
		session.commits = 0
		if err != nil {
			return err
		}
	}
	session.commits += 1
	return nil
}

// creates new value based on parameter and mode (add error return)
func updateValue(current any, new any, mode string) any {
	switch mode {
	case "append":
		// if it's a slice we append, else we create a new slice to append to.
		if tempSlice, ok := current.([]any); ok {
			return append(tempSlice, new)
		} else if current == nil {
			return []any{new}
		}
		return append([]any{current}, new)
	case "increment":
		// if it's an increment, either increase, or replace.
		if new == nil {
			new = float64(1) // default value
		}
		currentInt, currentOk := current.(float64)
		newInt, newOk := new.(float64)
		if currentOk && newOk {
			return currentInt + newInt
		} else if newOk {
			return newInt
		}
		slog.Error("incompatible values", "current", current, "new", new)
		return current
	default:
		return new
	}
}

// add value on nested key (e.g. [1,2,3] > map[1][2][3] = value)
func insert(data any, path []string, value any, mode string) error {
	current := data
	for index, step := range path {
		switch d := current.(type) {
		case map[string]any:
			if index == len(path)-1 { // last step in path
				d[step] = updateValue(d[step], value, mode)
				return nil
			}
			if _, ok := d[step]; !ok {
				d[step] = map[string]any{}
			}
			current = d[step]
		case []any:
			idx, err := strconv.Atoi(step)
			if err != nil || idx < 0 || idx >= len(d) {
				return errors.New("invalid index: " + step)
			}
			if index == len(path)-1 { // last step in path
				d[idx] = updateValue(d[idx], value, mode)
				return nil
			}
			current = d[idx]
		default:
			return errors.New("invalid path or data type")
		}
	}
	return nil
}

// delete value on nested key (e.g. [1,2,3] > map[1][2][3])
func pop(data any, path []string) error {
	current := data
	for index, step := range path {
		switch d := current.(type) {
		case map[string]any:
			if index == len(path)-1 { // last step
				delete(d, step)
				return nil
			}
			next, ok := d[step]
			if !ok {
				return errors.New("invalid path: " + step)
			}
			current = next
		case []any:
			idx, err := strconv.Atoi(step)
			if err != nil || idx < 0 || idx >= len(d) {
				return errors.New("invalid index: " + step)
			}
			if index == len(path)-1 { // last step
				d[idx] = nil // json safe delete
				return nil
			}
			current = d[idx]
		default:
			return errors.New("invalid path or data type")
		}
	}

	return nil
}

// used to query on multi-keys. E.g. [1,2,3] returns map[1,2,3] > value
func query(data any, path []string) (any, error) {
	if path[0] == "" {
		return data, nil
	}
	for _, step := range path {
		switch d := data.(type) {
		case map[string]any:
			if match, ok := d[step]; ok {
				data = match
			} else {
				return nil, errors.New("path not found")
			}
		case []any:
			i, err := strconv.Atoi(step) // convert string to int for slice index
			if err != nil || i < 0 || i >= len(d) {
				return nil, errors.New("path not found")
			}
			data = d[i]
		default:
			return nil, errors.New("path not found")
		}
	}
	return data, nil
}

// if you want to put a folder path before accessing the json, use '--'
func formatFilename(filename string) string {
	if strings.Contains(filename, "--") {
		return strings.ReplaceAll(filename, "--", "/")
	}
	return filename
}

// returns the item to apply CRUD operations on
func setItem(request Request) (Item, error) {
	if len(request.Key) > 0 {
		return EntryItem{}, nil
	}
	return TableItem{}, nil
}

// access level used to authenticate requests based on method
func setAccess(method string) int {
	level := session.config.access
	if method != "GET" {
		level++
	}
	slog.Debug("request auth:", "level", level)
	return level
}

// generates (or) selects a username and password
func initCredentials() (string, string) {
	username := os.Getenv("EM_USERNAME")
	if username == "" {
		username = "admin"
		slog.Info("set credentials:", "username", username)
	}
	password := os.Getenv("EM_PASSWORD")
	if password == "" {
		b := make([]byte, 12)
		rand.Read(b) //nolint:errcheck
		password = base64.URLEncoding.EncodeToString(b)
		slog.Info("set credentials:", "password", password)
	}
	return username, password
}

// selects the connector (fs interface) based on env variable
func initConnector() emmerFs.FileSystem {
	if os.Getenv("EM_CONNECTOR") == "S3" {
		return emmerFs.SetupS3()
	}
	return emmerFs.SetupLocal()
}

// selects the number of operations needed before a write action to fs
func initCache() int {
	commit := 1
	commitEnv := os.Getenv("EM_COMMITS")
	if commitEnv != "" {
		commitInt, err := strconv.Atoi(commitEnv)
		if err != nil {
			slog.Error("illegal commit strategy:", "EM_COMMITS", commitEnv)
			return 1
		}
		commit = commitInt
	}
	slog.Info("cache strategy set:", "commits", commit)
	return commit
}

// access 2 (default) means auth on put/del/get, 1 only on get, 0 on no methods
func initAccess() int {
	access := 2
	accessEnv := os.Getenv("EM_ACCESS")
	if accessEnv != "" {
		accessInt, err := strconv.Atoi(accessEnv)
		if err != nil {
			slog.Error("illegal access strategy:", "EM_ACCESS", accessEnv)
			return 2
		}
		access = accessInt
	}
	slog.Info("access set:", "level", access)
	return access
}
