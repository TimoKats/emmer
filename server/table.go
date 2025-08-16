package server

import (
	"errors"
)

func (table *TablePayload) exists() bool {
	filePath, err := fs.Fetch(table.Name)
	return err == nil && len(filePath) != 0
}

func (table *TablePayload) create() error {
	return fs.CreateJSON(table.Name)
}

func (table *TablePayload) add() error {
	if table.exists() {
		return errors.New("table '" + table.Name + "' already exists")
	}
	return table.create()
}

func (table *TablePayload) del() error {
	if !table.exists() {
		return errors.New("table '" + table.Name + "' doesn't exist")
	}
	return fs.DeleteFile(table.Name)
}
