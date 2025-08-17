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
func query(body []byte, path Path) (Response, error) { // add page param?
	var query QueryPayload
	if err := json.Unmarshal(body, &query); err != nil {
		return Response{}, err
	}
	return query.execute()
}

// Upon init, selects which filesystem to use based on env variable.
func init() {
	env := os.Getenv("EMMER_FS")
	switch env {
	case "aws":
		log.Println("aws not implemented yet")
	default:
		fs = io.LocalIO{Folder: "data"}
	}
	if fs == nil {
		log.Fatal("EMMER_FS '" + env + "' invalid")
	}
	log.Println("selected " + fs.Info())
}
