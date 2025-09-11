package server

import (
	"log"
)

type TableItem struct{}

// check if table exists, if yes, remove table.
func (TableItem) Del(request Request) Response {
	log.Printf("deleting table: %s", request.Table)
	if _, err := config.fs.Fetch(request.Table); err != nil {
		return formatResponse(err, "", nil)
	}
	err := config.fs.DeleteFile(request.Table)
	return formatResponse(err, "deleted "+request.Table, nil)
}

// parses payload of table, and creates it if it doesn't exist.
func (TableItem) Add(request Request) Response {
	log.Printf("creating table: %s", request.Table)
	if _, err := config.fs.Fetch(request.Table); err == nil {
		return formatResponse(err, "", nil)
	}
	err := config.fs.CreateJSON(request.Table)
	return formatResponse(err, "added "+request.Table, nil)
}

// queries tables (so not table contents)
func (TableItem) Query(request Request) Response {
	result := make(map[string]any)
	files, err := config.fs.List()
	if err != nil {
		return formatResponse(err, "", nil)
	}
	// iterate over json files (have to do type assertion because it's any)
	for filename, contents := range files {
		if request.Table+".json" == filename {
			result[filename] = contents
		}
	}
	return formatResponse(nil, "queried tables", result)
}
