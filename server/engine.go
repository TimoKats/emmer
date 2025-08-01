package server

import (
	"log"
	"os"

	io "github.com/TimoKats/emmer/server/io"

	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
)

var fs io.IO

// table

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

// entry

func (entry *EntryPayload) add() error {
	path, err := fs.GetFileByName("data/", entry.TableName)

	if err != nil {
		return err
	}
	switch filepath.Ext(path) {
	case ".csv":
		return fs.AppendCSV(path, entry.Values)
	case ".json":
		return fs.WriteJSON(path, entry.Key, entry.Value)
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

func Del(body []byte, path Path) error {
	var err error

	switch path {
	case Table:
		var table TablePayload
		if err = json.Unmarshal(body, &table); err == nil {
			err = table.del()
		}
	default:
		err = errors.New(fmt.Sprint(path) + ", unknown add option")
	}

	return err
}

func Query() error {
	return nil
}

func init() {
	env := os.Getenv("EMMER_FS")
	switch env {
	case "aws":
		log.Println("aws not implemented yet")
	default:
		fs = io.LocalIO{Folder: "/home/timokats/.emmer/"}
	}
	if fs == nil {
		log.Println("EMMER_FS '" + env + "' invalid")
		os.Exit(1)
	}
	log.Println("selected " + fs.Info())
}
