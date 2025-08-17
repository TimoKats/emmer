package server

import (
	"errors"
)

// Return true if table exists in filesystem.
func (table *TablePayload) exists() bool {
	filePath, err := fs.Fetch(table.Name)
	return err == nil && len(filePath) != 0
}

// Creates empty JSON with table name.
func (table *TablePayload) create() error {
	return fs.CreateJSON(table.Name)
}

// Called by engine. Check if table exists, if not, create table.
func (table *TablePayload) add() error {
	if table.exists() {
		return errors.New("table '" + table.Name + "' already exists")
	}
	return table.create()
}

// Calles by engine. Check if table exists, if yes, remove table.
func (table *TablePayload) del() error {
	if !table.exists() {
		return errors.New("table '" + table.Name + "' doesn't exist")
	}
	return fs.DeleteFile(table.Name)
}
