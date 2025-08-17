package server

// Applies query object to data (json) and returns the result.
func (query *QueryPayload) apply(data map[string]any) map[string]any {
	if len(query.Key) == 0 {
		return data
	}
	if _, ok := data[query.Key]; ok {
		filteredData := map[string]any{
			query.Key: data[query.Key],
		}
		data = filteredData
	} else { // return empty result if key not found
		data = make(map[string]any)
	}
	return data
}

// Gets JSON table data, and applies the query.
func (query *QueryPayload) execute() (Response, error) { // check this
	var err error
	var path string
	var response Response
	if path, err = fs.Fetch(query.TableName); err != nil {
		return response, err
	}
	if data, err := fs.ReadJSON(path); err == nil {
		response.Result = query.apply(data)
	}
	return response, err
}
