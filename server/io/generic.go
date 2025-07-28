package server

import (
	"errors"
)

type Separator string
type Format string
type Fs int

const (
	Comma     Separator = ","
	Semicolon Separator = ";"
	Tab       Separator = "\t"
	Pipe      Separator = "|"
)

const (
	Json Format = "json"
	Csv         = "csv"
)

const (
	Local Fs = iota
	Aws      // for example...
)

func (sep Separator) Rune() rune {
	runes := []rune(sep)
	if len(runes) == 0 {
		return ';'
	}
	return runes[0]
}

func CreateJSON(path string, fs Fs) error {
	switch fs {
	case Local:
		return createLocalJSON(path)
	case Aws:
		return errors.New("AWS not implemented")
	default:
		return errors.New("Unknown fs type")
	}
}

func CreateCSV(path string, columns []string, sep Separator, fs Fs) error {
	switch fs {
	case Local:
		return createLocalCSV(path, columns, sep)
	case Aws:
		return errors.New("AWS not implemented")
	default:
		return errors.New("Unknown fs type")
	}
}

func UpdateCSV(path string, values []string, fs Fs) error {
	switch fs {
	case Local:
		if sep, cols := getLocalCSVInfo(path); cols == len(values) {
			return appendLocalCSV(path, values, sep)
		}
		return errors.New("number of values incompatible with table")
	default:
		return errors.New("Unknown fs type")
	}
}

func UpdateJSON(path string, key string, value any, fs Fs) error {
	switch fs {
	case Local:
		if data, err := getLocalJson(path); err == nil {
			data[key] = value
			return writeLocalJson(path, data)
		}
		return nil // TO FIX!
	default:
		return errors.New("Unknown fs type")
	}
}
