package server

// enums

type Format string
type Path int

const (
	Json Format = "json"
	Csv  Format = "csv"
)

const (
	Table Path = iota
	Entry
)

// structs

type Response struct {
	Message string

	// for queries
	Page          int
	TabularResult [][]string
	MapResult     map[string]any
}

type TablePayload struct {
	Name       string   `json:"name"`
	Columns    []string `json:"columns"`
	FileFormat Format   `json:"format"`
	Sep        string   `json:"sep"`
}

type EntryPayload struct {
	TableName string `json:"table"`
	// csv
	Values []string `json:"values"`
	// json
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type Filter struct {
	With    []string `json:"with"`
	Without []string `json:"without"`
}

type QueryPayload struct {
	Filters   map[string][]Filter `json:"filter"`
	Key       string              `json:"key"` // for jsons
	TableName string              `json:"table"`
}
