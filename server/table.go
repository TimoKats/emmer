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
	if _, err := session.fs.Fetch(request.Table); err != nil {
		return Response{Data: nil, Error: err}
	}
	// delete file (and reset cache)
	err := session.fs.DeleteFile(request.Table)
	session.cache.tables = nil
	return Response{Data: "deleted " + request.Table, Error: err}
}

// parses payload of table, and creates it if it doesn't exist.
func (TableItem) Add(request Request) Response {
	log.Printf("creating table: %s", request.Table)
	// check if table exists
	if _, err := session.fs.Fetch(request.Table); err == nil {
		return Response{Data: nil, Error: errors.New("table already exists")}
	}
	// create table and add to cache
	err := session.fs.CreateJSON(request.Table, request.Value)
	if err == nil {
		session.cache.tables = append(session.cache.tables, request.Table)
	}
	return Response{Data: "added " + request.Table, Error: err}
}

// queries tables (so not table contents)
func (TableItem) Get(request Request) Response {
	log.Printf("querying table meta-data: %s", request.Table)
	// fetch all tables
	if len(session.cache.tables) != 0 {
		return Response{Data: session.cache.tables, Error: nil}
	}
	// if no cache, read directly
	files, err := session.fs.List()
	session.cache.tables = files
	return Response{Data: files, Error: err}
}
