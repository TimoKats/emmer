package server

import (
	"fmt"

	emmerFs "github.com/TimoKats/emmer/server/fs"

	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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
		if strings.Contains(response.Error.Error(), "not found") || strings.Contains(response.Error.Error(), "404") {
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

// tries reading data from cache, reads from filesystem as backup
func read(filename string, mode string) (map[string]any, error) {
	if data, ok := session.cache.data[filename]; ok && mode != "fs" {
		log.Println("reading data from cache")
		return data, nil
	}
	log.Println("reading data from filesystem")
	data, err := session.fs.Get(filename)
	if err == nil {
		session.cache.data[filename] = data
	}
	return data, err
}

// write to cache, and potentially to filesystem (depending on commit strategy)
func write(table string, data map[string]any) error {
	session.cache.data[table] = data
	if session.config.commit == session.commits {
		log.Println("writing to filesystem")
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
		log.Printf("%s, %s both not numeric", current, new)
		return current
	default:
		return new
	}
}

// add value on nested key (e.g. [1,2,3] > map[1][2][3] = value)
func insert(data map[string]any, keys []string, value any, mode string) error {
	current := data
	for i, key := range keys {
		if i == len(keys)-1 {
			current[key] = updateValue(current[key], value, mode)
		} else {
			if _, ok := current[key]; !ok {
				current[key] = make(map[string]any)
			}
			next, ok := current[key].(map[string]any)
			if !ok {
				return errors.New("conflict at key: " + key)
			}
			current = next
		}
	}
	return nil
}

// delete value on nested key (e.g. [1,2,3] > map[1][2][3])
func pop(data map[string]any, key []string) error {
	keyFound := true
	current := data
	for index, step := range key {
		next, ok := current[step].(map[string]any)
		if !ok {
			if _, ok = current[step]; !ok {
				keyFound = false
			}
			break
		}
		if index < len(key)-1 {
			current = next
		}
	}
	if keyFound {
		delete(current, key[len(key)-1])
		return nil
	}
	return errors.New("key not found in table")
}

// used to query on multi-keys. E.g. [1,2,3] returns map[1,2,3] > value
func query(data map[string]any, key []string) (any, error) {
	var current any = data
	if len(key) == 0 || key[0] == "" {
		return data, nil
	}
	for _, step := range key {
		switch typed := current.(type) {
		case map[string]any:
			val, ok := typed[step]
			if !ok {
				return nil, errors.New("key " + step + " not found in map")
			}
			current = val
		case []any:
			index, err := strconv.Atoi(step)
			if err != nil {
				return nil, errors.New("invalid index " + step + " for list")
			}
			if index < 0 || index >= len(typed) {
				return nil, errors.New("index " + step + " out of bounds")
			}
			current = typed[index]
		default:
			return nil, errors.New("cannot descend into type")
		}
	}
	return current, nil
}

func formatFilename(filename string) string {
	if strings.Contains(filename, "--") {
		return strings.ReplaceAll(filename, "--", "/")
	}
	return filename
}

func initCredentials() (string, string) {
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
	return username, password
}

func initConnector() emmerFs.FileSystem {
	if os.Getenv("EM_CONNECTOR") == "S3" {
		return emmerFs.SetupS3()
	}
	return emmerFs.SetupLocal()
}

func initCache() int {
	commit := 1
	commitEnv := os.Getenv("EM_COMMIT")
	if commitEnv != "" {
		commitInt, err := strconv.Atoi(commitEnv)
		if err != nil {
			fmt.Printf("Error converting commit strategy to int: %v", err)
			return 1
		}
		commit = commitInt
	}
	return commit
}
