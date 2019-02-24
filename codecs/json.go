package codecs

import (
	"bytes"
	"encoding/json"
	"strings"
)

type JsonScannerConfig struct{}

type JsonScanner struct{}

func NewJsonScanner(config *JsonScannerConfig) (*JsonScanner, error) {
	return &JsonScanner{}, nil
}

func (sc *JsonScanner) Scan(s string) (interface{}, error) {
	buf := bytes.NewBufferString(s)
	var data interface{}
	err := json.NewDecoder(buf).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type JsonPrinterConfig struct{}

type JsonPrinter struct{}

func NewJsonPrinter(config *JsonPrinterConfig) (*JsonPrinter, error) {
	return &JsonPrinter{}, nil
}

func (p *JsonPrinter) Print(v interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}
