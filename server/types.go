package server

// enums

type Format string
type Operator string
type Path int

const (
	Equals  Operator = "equals"
	Differs          = "differs"
	Before           = "before"
	After            = "after"
	Above            = "above"
	Below            = "below"
)

const (
	Json Format = "json"
	Csv         = "csv"
)

const (
	Table Path = iota
	Entry
)

// structs

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`

	Result map[string]any // This will contain query results.
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
	Operator Operator `json:"op"`
	Args     any      `json:"args"`
}

type QueryPayload struct {
	Filters   map[string][]Filter `json:"filter"`
	TableName string              `json:"table"`
}
