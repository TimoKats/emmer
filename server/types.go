package server

import emmerFs "github.com/TimoKats/emmer/server/fs"

type Config struct {
	autoTable bool
	username  string
	password  string
	fs        emmerFs.FileSystem
}

type Response struct {
	StatusCode int
	Message    string `json:"message"`
	Result     any    `json:"result"`
}

type Request struct {
	Method string // to enum
	Table  string
	Key    []string
	Mode   string // increment, append, empty
	Value  any
}

type Item interface {
	Add(request Request) Response
	Del(request Request) Response
	Query(request Request) Response
}

// generic helper function that formats the response object
func formatResponse(err error, message string, result any) Response {
	if err != nil {
		return Response{
			StatusCode: 500,
			Message:    err.Error(),
		}
	}
	if len(message) == 0 {
		message = "success"
	}
	return Response{
		StatusCode: 200,
		Message:    message,
		Result:     result,
	}
}
