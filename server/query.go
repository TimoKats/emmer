package server

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

func (query *QueryPayload) execute() (Response, error) { // check this
	var err error
	var response Response
	path, err := fs.Fetch(query.TableName)
	if err != nil {
		return response, err
	}

	data, err := fs.ReadJSON(path)
	if err == nil {
		response.Result = query.apply(data)
	}

	return response, err
}
