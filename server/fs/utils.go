package server

import (
	"errors"
	"log"
)

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
		currentInt, currentOk := current.(float64)
		newInt, newOk := new.(float64)
		if currentOk && newOk {
			return currentInt + newInt
		} else if newOk {
			return newInt
		}
		log.Printf("%s not numeric", current)
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
