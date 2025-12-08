package server

import (
	"errors"
	"log"
)

type TableItem struct{}

// check if table exists, if yes, remove table.
func (TableItem) Del(request Request) Response {
	log.Printf("deleting table: %s", request.Table)
	// check if table exists
	if _, err := read(request.Table, request.Mode); err != nil {
		return Response{Data: nil, Error: err}
	}
	// delete file (and reset cache)
	err := session.fs.Del(request.Table)
	session.cache.tables = nil
	delete(session.cache.data, request.Table)
	return Response{Data: "deleted " + request.Table, Error: err}
}

// parses payload of table, and creates it if it doesn't exist.
func (TableItem) Add(request Request) Response {
	log.Printf("creating table: %s", request.Table)
	// check if table exists
	if _, err := read(request.Table, request.Mode); err == nil {
		return Response{Data: nil, Error: errors.New("table already exists")}
	}
	// create table and add to cache
	data, ok := request.Value.(map[string]any)
	if !ok {
		return Response{Data: nil, Error: errors.New("value not json")}
	}
	err := write(request.Table, data)
	if err == nil {
		session.cache.tables = append(session.cache.tables, request.Table)
	}
	return Response{Data: "added " + request.Table, Error: err}
}

// queries tables (so not table contents)
func (TableItem) Get(request Request) Response {
	log.Printf("querying tables: %s", request.Table)
	// fetch all tables
	if len(session.cache.tables) == 0 {
		files, err := session.fs.Ls()
		if err != nil {
			return Response{Data: nil, Error: err}
		}
		session.cache.tables = files
	}
	if len(request.Table) != 0 {
		for _, file := range session.cache.tables {
			if file == formatFilename(request.Table) {
				return Response{Data: []string{file}, Error: nil}
			}
		}
		return Response{Data: nil, Error: errors.New(request.Table + " not found")}
	}
	return Response{Data: session.cache.tables, Error: nil}
}
