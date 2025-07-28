package server

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func getLocalFirstLine(path string) (string, error) {
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

func getLocalCSVInfo(path string) (Separator, int) {
	line, err := getLocalFirstLine(path)
	if err != nil {
		log.Println("when getting seperator: ", err)
	}
	maxCols := 0
	separators := []Separator{Comma, Semicolon, Tab, Pipe}
	var bestSep Separator = Semicolon // default

	for _, sep := range separators {
		parts := strings.Split(line, string(sep))
		if len(parts) > maxCols {
			maxCols = len(parts)
			bestSep = sep
		}
	}

	return bestSep, maxCols
}

func createLocalCSV(path string, columns []string, sep Separator) error {
	data := [][]string{columns}
	f, _ := os.Create(path)
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = sep.Rune()
	return w.WriteAll(data)
}

func appendLocalCSV(path string, values []string, sep Separator) error {
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
