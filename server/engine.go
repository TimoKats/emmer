package server

import (
	io "github.com/TimoKats/emmer/server/io"

	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
)

// table

func (table *TablePayload) path() string { // should be in local!
	return "data/" + table.Name + "." + string(table.FileFormat)
}

func (table *TablePayload) exists() bool { // should be in local!
	filePath, err := io.GetFileByName("data/", table.Name)
	return err == nil && len(filePath) != 0
}

func (table *TablePayload) valid() bool { // should be in local!
	var validFormats = []io.Format{io.Json, io.Csv}
	return slices.Contains(validFormats, table.FileFormat)
}

func (table *TablePayload) create() error {
	switch table.FileFormat {
	case io.Json:
		return io.CreateJSON(table.path(), io.Local) // this can be an env
	case io.Csv:
		return io.CreateCSV(table.path(), table.Columns, table.Sep, io.Local)
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
	path, err := io.GetFileByName("data/", entry.TableName) // put th
	if err != nil {
		return err
	}
	switch filepath.Ext(path) {
	case ".csv":
		return io.UpdateCSV(path, entry.Values, io.Local)
	case ".json":
		return io.UpdateJSON(path, entry.Key, entry.Value, io.Local)
	default:
		return errors.New("file extension not supported")
	}
}

// actions

func Add(body []byte, path Path) error {
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

func Del() error {
	return nil
}

func Query() error {
	return nil
}
