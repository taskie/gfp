package codecs

import (
	"fmt"
	"time"
)

type DurationScannerConfig struct{}

type DurationScanner struct{}

func NewDurationScanner(config *DurationScannerConfig) (*DurationScanner, error) {
	return &DurationScanner{}, nil
}

func (sc *DurationScanner) Scan(s string) (interface{}, error) {
	return time.ParseDuration(s)
}

type DurationPrinterConfig struct{}

type DurationPrinter struct{}

func NewDurationPrinter(config *DurationPrinterConfig) (*DurationPrinter, error) {
	return &DurationPrinter{}, nil
}

func (p *DurationPrinter) Print(v interface{}) (string, error) {
	switch v2 := v.(type) {
	case time.Duration:
		return v2.String(), nil
	case int64:
		return time.Duration(v2).String(), nil
	case float64:
		return time.Duration(v2).String(), nil
	default:
		return "", fmt.Errorf("not duration: %+v (%T)", v, v)
	}
}
