package server

import (
	"errors"
	"path/filepath"
)

func (query *QueryPayload) queryMap(data map[string]any) map[string]any {
	if _, ok := data[query.Key]; ok {
		filteredData := map[string]any{
			query.Key: data[query.Key],
		}
		data = filteredData
	} else { // return empty result
		data = make(map[string]any)
	}
	return data
}

func (query *QueryPayload) queryTable(data [][]string) [][]string {
	if len(data) == 0 {
		return [][]string{}
	}
	// columns := data[0]
	// for // IM HERE MAKING A MAPPING BETWEEN colindex and colname!
	// data = data[:1]

	return data
}

func (query *QueryPayload) execute() (Response, error) {
	var err error
	var response Response
	path, err := fs.GetFileByName("data/", query.TableName)
	if err != nil {
		return response, err
	}

	switch filepath.Ext(path) {
	case ".csv":
		var data [][]string
		data, err = fs.ReadCSV(path)
		if err == nil {
			response.TabularResult = query.queryTable(data)
		}
	case ".json":
		var data map[string]any
		data, err = fs.ReadJSON(path)
		if err == nil {
			response.MapResult = query.queryMap(data)
		}
	default:
		err = errors.New("file extension not supported")
	}

	return response, err
}
