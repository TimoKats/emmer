package server

import (
	"errors"
	"log"
)

type TableItem struct{}

// check if table exists, if yes, remove table.
func (TableItem) Del(request Request) Response {
	log.Printf("deleting table: %s", request.Table)
	if _, err := config.fs.Fetch(request.Table); err != nil {
		return Response{Data: nil, Error: err}
	}
	err := config.fs.DeleteFile(request.Table)
	return Response{Data: "deleted " + request.Table, Error: err}
}

// parses payload of table, and creates it if it doesn't exist.
func (TableItem) Add(request Request) Response {
	log.Printf("creating table: %s", request.Table)
	if _, err := config.fs.Fetch(request.Table); err == nil {
		return Response{Data: nil, Error: errors.New("table already exists")}
	}
	err := config.fs.CreateJSON(request.Table, request.Value)
	return Response{Data: "added " + request.Table, Error: err}
}

// queries tables (so not table contents)
func (TableItem) Get(request Request) Response {
	log.Printf("querying table meta-data: %s", request.Table)
	// fetch all tables
	result := []string{}
	files, err := config.fs.List()
	if err != nil {
		return Response{Data: nil, Error: err}
	}
	// iterate over json files (have to do type assertion because it's any)
	for _, filename := range files {
		if request.Table+".json" == filename || request.Table == "" {
			result = append(result, filename)
		}
	}
	// handle not found error
	if len(result) == 0 && request.Table != "" {
		err = errors.New("table " + request.Table + " not found")
	}
	return Response{Data: result, Error: err}
}
