package server

import (
	"sync"

	emmerFs "github.com/TimoKats/emmer/server/fs"
)

// session

type LogBuffer struct {
	mu    sync.Mutex
	logs  []string
	limit int
}

type Config struct {
	username string
	password string
	commit   int
}

type Backup struct {
	table string
	value map[string]any
}

type Cache struct {
	tables []string
	data   map[string]map[string]any
	backup Backup
}

type Session struct {
	commits   int
	config    Config
	fs        emmerFs.FileSystem
	logBuffer *LogBuffer
	cache     Cache
}

// api

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
	Get(request Request) Response
}
