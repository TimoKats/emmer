package server

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Csv struct {
	Folder string
}

// takes filename and returns full path + extension to csv file
func (fs Csv) getPath(filename string) string {
	return filepath.Join(fs.Folder, filename) + ".csv"
}

// creates empty (or prefilled) CSV file at path
func (fs Csv) Put(filename string, value any) error {
	path := fs.getPath(filename)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck

	writer := csv.NewWriter(f)
	writer.Comma = ';'
	defer writer.Flush()

	if value == nil {
		return nil
	}

	rows, ok := value.([]map[string]any)
	if !ok {
		return errors.New("value must be []map[string]any")
	}

	if len(rows) == 0 {
		return nil
	}

	// collect headers from first row
	headers := []string{}
	for k := range rows[0] {
		headers = append(headers, k)
	}

	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, row := range rows {
		record := make([]string, len(headers))
		for i, h := range headers {
			if v, ok := row[h]; ok {
				record[i] = fmt.Sprint(v)
			}
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// reads CSV file into []map[string]any
func (fs Csv) Get(filename string) (any, error) {
	path := fs.getPath(filename)

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New("table " + filename + " not found")
	}
	defer f.Close() //nolint:errcheck

	reader := csv.NewReader(f)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("error reading csv file")
	}

	if len(records) == 0 {
		return []map[string]any{}, nil
	}

	headers := records[0]
	result := []map[string]any{}

	for _, row := range records[1:] {
		entry := map[string]any{}
		for i, h := range headers {
			if i < len(row) {
				entry[h] = row[i]
			}
		}
		result = append(result, entry)
	}

	return result, nil
}

// removes entire CSV file
func (fs Csv) Del(filename string) error {
	path := fs.getPath(filename)
	return os.Remove(path)
}

// list csv files in fs folder
func (fs Csv) Ls() ([]string, error) {
	files, err := os.ReadDir(fs.Folder)
	result := []string{}
	if err != nil {
		return result, err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".csv" {
			filename := strings.TrimSuffix(f.Name(), ".csv")
			result = append(result, filename)
		}
	}

	return result, nil
}

// creates new localFS instance with settings applied
func SetupCSV() *Csv {
	folder := selectFolder()
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		slog.Debug("created folder", "folder", folder)
		if err := os.Mkdir(folder, 0755); err != nil {
			slog.Error("can't create emmer folder:", "path", folder)
			os.Exit(1)
		}
	}
	slog.Info("selected CSV fs:", "folder", folder)
	return &Csv{
		Folder: folder,
	}
}
