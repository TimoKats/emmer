package server

import (
	"errors"
	"slices"
)

func (table *TablePayload) path() string { // should be in local!
	return "data/" + table.Name + "." + string(table.FileFormat)
}

func (table *TablePayload) exists() bool { // should be in local!
	filePath, err := fs.GetFileByName("data/", table.Name)
	return err == nil && len(filePath) != 0
}

func (table *TablePayload) valid() bool { // should be in local!
	var validFormats = []Format{Json, Csv}
	return slices.Contains(validFormats, table.FileFormat)
}

func (table *TablePayload) create() error {
	switch table.FileFormat {
	case Json:
		return fs.CreateJSON(table.path()) // this can be an env
	case Csv:
		return fs.CreateCSV(table.path(), table.Columns, table.Sep)
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

func (table *TablePayload) del() error {
	if !table.exists() {
		return errors.New("table '" + table.Name + "' doesn't exist")
	}
	return fs.Delete(table.path())
}
