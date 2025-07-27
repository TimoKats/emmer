package server

import (
	io "github.com/TimoKats/emmer/server/io"
)

// enums

type Path int
type Format string

const (
	Json  Format = "json"
	Jsonl        = "jsonl"
	Csv          = "csv"
)

const (
	Table Path = iota
	Entry
)

// structs

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type TablePayload struct {
	Name       string       `json:"name"`
	Columns    []string     `json:"columns"`
	FileFormat Format       `json:"format"`
	Sep        io.Separator `json:"sep"`
}

type EntryPayload struct {
	TableName string `json:"table"`
	// for csvs
	Values []string `json:"values"`
	// for jsons
	Key   string `json:"key"`
	Value any    `json:"value"`
}
