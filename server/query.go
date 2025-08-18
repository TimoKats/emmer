package server

import (
	"errors"
)

// Used to trim the body of the result using a template function.
func (query *QueryPayload) format() error {
	return nil
}

// Used to query on multi-keys. E.g. [1,2,3] returns map[1,2,3]
func (query *QueryPayload) filterEntry(data map[string]any) (map[string]any, error) {
	if len(query.Key) == 0 {
		return data, nil
	}
	current := data
	for _, step := range query.Key {
		match, ok := current[step].(map[string]any)
		if !ok {
			return make(map[string]any), errors.New(step + " not found")
		}
		current = match
	}
	return current, nil

}

// Gets JSON table data, and applies the query.
func (query *QueryPayload) entry() (Response, error) {
	var response Response
	// fetching table contents
	path, err := fs.Fetch(query.TableName)
	if err != nil {
		return response, err
	}
	// filter contents on query
	data, err := fs.ReadJSON(path)
	if err == nil {
		response.Result, err = query.filterEntry(data)
	}
	return response, err
}

// Lists the tables in the store along with some meta data.
func (query *QueryPayload) table() (Response, error) {
	var response Response
	response.Result = make(map[string]any)
	// fetch store contents
	files, err := fs.List()
	if err != nil {
		return response, err
	}
	// iterate over json files
	for filename, contents := range files {
		if len(query.TableName) == 0 || query.TableName+".json" == filename {
			response.Result[filename] = contents
		}
	}
	return response, nil
}
