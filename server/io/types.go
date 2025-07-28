package server

type Separator string
type Fs int

const (
	Comma     Separator = ","
	Semicolon Separator = ";"
	Tab       Separator = "\t"
	Pipe      Separator = "|"
)

func (sep Separator) Rune() rune {
	runes := []rune(sep)
	if len(runes) == 0 {
		return ';'
	}
	return runes[0]
}

type IO interface {
	// generic
	GetFileByName(path string, filename string) (string, error)

	// json
	ReadJSON(path string) (map[string]any, error)
	WriteJSON(path string, key string, value any) error
	CreateJSON(path string) error

	// csv
	AppendCSV(path string, values []string) error
	CreateCSV(path string, columns []string, sep Separator) error
}
