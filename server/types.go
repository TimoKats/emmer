package server

import (
	"sync"

	emmerFs "github.com/TimoKats/emmer/server/fs"
)

type LogBuffer struct {
	mu    sync.Mutex
	logs  []string
	limit int
}

type Config struct {
	logBuffer *LogBuffer
	autoTable bool
	username  string
	password  string
	fs        emmerFs.FileSystem
}

type Response struct {
	Error error
	Data  any
}

type Request struct {
	Method string // get, put, delete
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
