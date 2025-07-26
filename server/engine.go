package server

import (
	"encoding/json"
	"errors"
	"fmt"
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
		return createCSV(table.path(), table.Columns, table.Sep)
	default:
		return errors.New("unsupported file format")
	}
}

func (table *TablePayload) add() error {
	if !table.valid() {
		return errors.New("invalid payload")
	}
	if table.exists() {
		return errors.New("table '" + table.Name + "' already exists")
	}
	return table.create()
}

// entry

func (entry *EntryPayload) add() error {
	path, err := getFile("data/", entry.TableName)
	if err != nil {
		return err
	}
	switch filepath.Ext(path) {
	case ".csv":
		if sep, cols := getCSVInfo(path); cols == len(entry.Values) {
			return appendCSVEntry(path, entry.Values, sep)
		}
		return errors.New("number of values incompatible with table")
	case ".json":
		if data, err := getJson(path); err == nil {
			data[entry.Key] = entry.Value
			return updateJson(path, data)
		}
		return err
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
	case Entry:
		var entry EntryPayload
		if err = json.Unmarshal(body, &entry); err == nil {
			err = entry.add()
		}
	default:
		err = errors.New(fmt.Sprint(path) + ", unknown add option")
	}

	return err
}
