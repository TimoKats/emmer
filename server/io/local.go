package server

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
)

type LocalIO struct {
	path string
}

func getFirstLine(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimRight(line, "\r\n"), nil
}

func getCSVInfo(path string) (Separator, int) {
	line, err := getFirstLine(path)
	if err != nil {
		log.Println("when getting seperator: ", err)
	}
	maxCols := 0
	separators := []Separator{Comma, Semicolon, Tab, Pipe}
	var bestSep Separator = Semicolon // default

	for _, sep := range separators {
		parts := strings.Split(line, string(sep))
		if len(parts) > maxCols {
			maxCols = len(parts)
			bestSep = sep
		}
	}

	return bestSep, maxCols
}

func (io LocalIO) CreateCSV(path string, columns []string, sep Separator) error {
	data := [][]string{columns}
	f, _ := os.Create(path)
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = sep.Rune()
	return w.WriteAll(data)
}

func (io LocalIO) AppendCSV(path string, values []string) error {
	sep, cols := getCSVInfo(path)
	if cols != len(values) {
		return errors.New("number of values and columns incompatible")
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = sep.Rune()
	err = w.Write(values)
	if err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}

func (io LocalIO) CreateJSON(path string) error {
	f, err := os.Create(path)
	defer f.Close()
	f.WriteString("{}")
	return err
}

func (io LocalIO) ReadJSON(path string) (map[string]any, error) {
	var data map[string]any
	file, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(file, &data)
	return data, err
}

func (io LocalIO) WriteJSON(path string, key string, value any) error {
	data, err := io.ReadJSON(path)
	if err != nil {
		return err
	}
	data[key] = value
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}
