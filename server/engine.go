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

func del(body []byte, path Path) error {
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

func query() error {
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
