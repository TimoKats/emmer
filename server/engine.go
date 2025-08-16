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

func query(body []byte, path Path) (Response, error) { // add page param?
	var query QueryPayload
	if err := json.Unmarshal(body, &query); err != nil {
		return Response{}, err
	}
	return query.execute()
}

func (r Response) format() ([]byte, error) {
	output := make(map[string]any)
	output["message"] = r.Message
	output["page"] = r.Page

	// Add only one result field
	if r.TabularResult != nil && r.MapResult != nil {
		return nil, fmt.Errorf("both TabularResult and MapResult are set")
	}

	if r.TabularResult != nil {
		output["result"] = r.TabularResult
	} else if r.MapResult != nil {
		output["result"] = r.MapResult
	}

	log.Println(output)

	return json.Marshal(output)
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
		log.Fatal("EMMER_FS '" + env + "' invalid")
	}
	log.Println("selected " + fs.Info())
}
