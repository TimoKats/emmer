package server

import (
	"log"
	"os"

	io "github.com/TimoKats/emmer/server/io"

	"encoding/json"
	"errors"
	"fmt"
)

var fs io.IO

// Adds table or entry to table based on API path.
func add(body []byte, path Path) error {
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

// Deletes table or entry to table based on API path.
func del(body []byte, path Path) error {
	var err error
	switch path {
	case Table:
		var table TablePayload
		if err = json.Unmarshal(body, &table); err == nil {
			err = table.del()
		}
	case Entry:
		var entry EntryPayload
		if err = json.Unmarshal(body, &entry); err == nil {
			err = entry.del()
		}
	default:
		err = errors.New(fmt.Sprint(path) + ", unknown add option")
	}
	return err
}

// Queries table based on request body (path for querying meta data).
func query(body []byte, path Path) (Response, error) {
	var query QueryPayload
	if err := json.Unmarshal(body, &query); err != nil {
		return Response{}, err
	}
	switch path {
	case Table:
		return query.table()
	case Entry:
		return query.entry()
	default:
		return Response{}, errors.New(fmt.Sprint(path) + ", unknown add option")
	}
}

// Upon init, selects which filesystem to use based on env variable.
func init() {
	switch os.Getenv("EMMER_FS") {
	// case "aws": < this will be the pattern
	// 	log.Println("aws not implemented yet")
	default:
		fs = io.LocalIO{Folder: "data"}
	}
	log.Println("selected " + fs.Info())
}
