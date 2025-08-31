package server

import (
	"encoding/json"
	"errors"
)

type TableItem struct{}

// check if table exists, if yes, remove table.
func (TableItem) Del(payload []byte) error {
	var table TablePayload
	if err := json.Unmarshal(payload, &table); err != nil {
		return err
	}
	return fs.DeleteFile(table.Name)
}

// parses payload of table, and creates it if it doesn't exist.
func (TableItem) Add(payload []byte) error {
	var table TablePayload
	if err := json.Unmarshal(payload, &table); err != nil {
		return err
	}
	// fetching table contents
	if _, err := fs.Fetch(table.Name); err == nil {
		return errors.New("table " + table.Name + " already exists")
	}
	return fs.CreateJSON(table.Name)
}

// queries tables (so not table contents)
func (TableItem) Query(payload []byte) (Response, error) {
	// parse query payload into object
	var response Response
	var query QueryPayload
	response.Result = make(map[string]any)
	if err := json.Unmarshal(payload, &query); err != nil {
		return Response{}, err
	}
	// fetch store contents
	files, err := fs.List()
	if err != nil {
		return response, err
	}
	// iterate over json files (have to do type assertion because it's any)
	if result, ok := response.Result.(map[string]any); ok {
		for filename, contents := range files {
			if len(query.TableName) == 0 || query.TableName+".json" == filename {
				result[filename] = contents
			}
		}
	}
	return response, nil
}
