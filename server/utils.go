package server

import (
	"errors"
	"log"
	"strconv"
)

// tries reading data from cache, reads from filesystem as backup
func read(filename string) (map[string]any, error) {
	if data, ok := session.cache.data[filename]; ok {
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
func write(request Request, data map[string]any) error {
	session.cache.data[request.Table] = data
	log.Println(session.config.commit, session.commits)
	if session.config.commit == session.commits {
		log.Println("writing to filesystem")
		err := session.fs.Put(request.Table, data)
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
func insertNested(data map[string]any, keys []string, value any, mode string) error {
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
func deleteNested(data map[string]any, key []string) error {
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
func findKey(data map[string]any, key []string) (any, error) {
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
