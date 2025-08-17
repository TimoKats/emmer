package server

import "errors"

func (query *QueryPayload) apply(data map[string]any) (map[string]any, error) {
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
func (query *QueryPayload) execute() (Response, error) { // check this
	var response Response
	// fetching table contents
	path, err := fs.Fetch(query.TableName)
	if err != nil {
		return response, err
	}
	// filter contents on query
	data, err := fs.ReadJSON(path)
	if err == nil {
		response.Result, err = query.apply(data)
	}
	return response, err
}
