package server

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
	Row
)

// structs

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type TablePayload struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`

	FileFormat Format `json:"format"`
}

type RowPayload struct {
	TableName string   `json:"table"`
	Values    []string `json:"values"`
}
