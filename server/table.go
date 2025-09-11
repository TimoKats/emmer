package server

import (
	"errors"
	"log"
)

type TableItem struct{}

// check if table exists, if yes, remove table.
func (TableItem) Del(request Request) error {
	log.Printf("deleting table: %s", request.Table)
	if _, err := config.fs.Fetch(request.Table); err != nil {
		return errors.New("table " + request.Table + " doesn't exist")
	}
	return config.fs.DeleteFile(request.Table)
}

// parses payload of table, and creates it if it doesn't exist.
func (TableItem) Add(request Request) error {
	log.Printf("creating table: %s", request.Table)
	if _, err := config.fs.Fetch(request.Table); err == nil {
		return errors.New("table " + request.Table + " already exists")
	}
	return config.fs.CreateJSON(request.Table)
}

// queries tables (so not table contents)
func (TableItem) Query(request Request) (Response, error) {
	var response Response
	response.Result = make(map[string]any)
	files, err := config.fs.List()
	if err != nil {
		return response, err
	}
	// iterate over json files (have to do type assertion because it's any)
	if result, ok := response.Result.(map[string]any); ok {
		for filename, contents := range files {
			if request.Table+".json" == filename {
				result[filename] = contents
			}
		}
	}
	return response, nil
}
