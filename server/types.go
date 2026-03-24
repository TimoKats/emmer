package server

import (
	emmerFs "github.com/TimoKats/emmer/server/fs"
)

// session

type Config struct {
	username string
	password string
	commit   int
	access   int
}

type Cache struct {
	tables []string
	data   map[string]any
}

type Session struct {
	commits int
	config  Config
	cache   Cache
}

// api

type Response struct {
	Error error
	Data  any
}

type Request struct {
	Fs     emmerFs.File
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
