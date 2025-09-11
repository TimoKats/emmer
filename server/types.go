package server

import emmerFs "github.com/TimoKats/emmer/server/fs"

type Config struct {
	autoTable bool
	username  string
	password  string
	fs        emmerFs.FileSystem
}

type Response struct {
	Message string `json:"message"`
	Result  any    `json:"result"`
}

type Request struct {
	Method string // to enum
	Table  string
	Key    []string
	Mode   string // increment, append, empty
	Value  any
}

type Item interface {
	Add(request Request) error
	Del(request Request) error
	Query(request Request) (Response, error)
}
