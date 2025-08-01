package server

type Fs int

type IO interface {
	// generic
	GetFileByName(path string, filename string) (string, error)

	// json
	ReadJSON(path string) (map[string]any, error)
	WriteJSON(path string, key string, value any) error
	CreateJSON(path string) error

	// csv
	AppendCSV(path string, values []string) error
	CreateCSV(path string, columns []string, sep string) error
}
