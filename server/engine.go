package server

import (
	"encoding/json"
	"errors"
	"log"
	"path/filepath"
	"slices"
)

// table

func (table *TablePayload) path() string {
	return "data/" + table.Name + "." + string(table.FileFormat)
}

func (table *TablePayload) exists() bool {
	filePath, err := getFile("data/", table.Name)
	return err == nil && len(filePath) != 0
}

func (table *TablePayload) valid() bool {
	var validFormats = []Format{Json, Jsonl, Csv}
	return slices.Contains(validFormats, table.FileFormat)
}

func (table *TablePayload) create() error {
	switch table.FileFormat {
	case Json, Jsonl:
		return createJSON(table.path())
	case Csv:
		return createCSV(table.path(), table.Columns, ';')
	default:
		return errors.New("unsupported file format")
	}
}

func (table *TablePayload) add() error {
	if !table.valid() {
		return errors.New("invalid file format, check docs")
	}
	if table.exists() {
		return errors.New("table '" + table.Name + "' already exists")
	}
	return table.create()
}

// row

func (row *RowPayload) add() error {
	path, err := getFile("data/", row.TableName)
	if err != nil {
		return err
	}
	switch filepath.Ext(path) {
	case ".csv":
		return appendCSVRow(path, row.Values, ';')
	case ".json":
		log.Println("json!")
		return nil
	default:
		return errors.New("file extension not supported")
	}
}

// generics

func AddSwitch(body []byte, path Path) error {
	var err error

	switch path {
	case Table:
		var table TablePayload
		if err = json.Unmarshal(body, &table); err == nil {
			err = table.add()
		}
	case Row:
		var row RowPayload
		if err = json.Unmarshal(body, &row); err == nil {
			err = row.add()
		}
	default:
		err = errors.New("how did you get this error? hacker?")
	}

	return err
}
