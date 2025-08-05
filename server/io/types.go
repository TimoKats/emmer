package server

type Fs int

type IO interface {
	// generic
	GetFileByName(path string, filename string) (string, error)
	DeleteTable(path string) error
	Info() string

	// json
	ReadJSON(path string) (map[string]any, error)
	AddJSON(path string, key string, value any) error
	DelJSON(path string, key string) error
	CreateJSON(path string) error

	// csv
	AppendCSV(path string, values []string) error
	CreateCSV(path string, columns []string, sep string) error
	ReadCSV(path string) ([][]string, error)
}
