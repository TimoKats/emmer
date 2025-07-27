package server

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strings"
)

// sep

type Separator string

const (
	Comma     Separator = ","
	Semicolon Separator = ";"
	Tab       Separator = "\t"
	Pipe      Separator = "|"
)

func (sep Separator) Rune() rune {
	runes := []rune(sep)
	if len(runes) == 0 {
		return ';'
	}
	return runes[0]
}

// file handling

func getFirstLine(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimRight(line, "\r\n"), nil
}

func GetCSVInfo(path string) (Separator, int) {
	line, err := getFirstLine(path)
	if err != nil {
		log.Println("when getting seperator: ", err)
	}
	separators := []Separator{Comma, Semicolon, Tab, Pipe}
	maxCols := 0
	var bestSep Separator = Comma // default

	for _, sep := range separators {
		parts := strings.Split(line, string(sep))
		if len(parts) > maxCols {
			maxCols = len(parts)
			bestSep = sep
		}
	}

	return bestSep, maxCols
}

func CreateCSV(path string, columns []string, sep Separator) error {
	data := [][]string{columns}
	f, _ := os.Create(path)
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = sep.Rune()
	return w.WriteAll(data)
}

func AppendCSV(path string, values []string, sep Separator) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = sep.Rune()
	err = w.Write(values)
	if err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}
