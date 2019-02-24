package codecs

import (
	"fmt"
	"strings"
	"time"
)

func presetToLayout(preset string) (string, error) {
	switch strings.ToLower(preset) {
	case "ansic":
		return time.ANSIC, nil
	case "unixdate":
		return time.UnixDate, nil
	case "rfc3339":
		return time.RFC3339, nil
	case "rfc3339nano":
		return time.RFC3339Nano, nil
	default:
		return "", fmt.Errorf("unknown preset: %s", preset)
	}
}

type TimeScannerConfig struct {
	Layout string
	Preset string
}

type TimeScanner struct {
	layout string
}

func NewTimeScanner(config *TimeScannerConfig) (*TimeScanner, error) {
	var err error
	layout := config.Layout
	if config.Layout == "" {
		layout, err = presetToLayout(config.Preset)
		if err != nil {
			return nil, err
		}
	}
	return &TimeScanner{
		layout: layout,
	}, nil
}

func (sc *TimeScanner) Scan(s string) (interface{}, error) {
	return time.Parse(sc.layout, s)
}

type TimePrinterConfig struct {
	Layout string
	Preset string
}

type TimePrinter struct {
	layout string
}

func NewTimePrinter(config *TimePrinterConfig) (*TimePrinter, error) {
	var err error
	layout := config.Layout
	if config.Layout == "" {
		layout, err = presetToLayout(config.Preset)
		if err != nil {
			return nil, err
		}
	}
	return &TimePrinter{
		layout: layout,
	}, nil
}

func (p *TimePrinter) Print(v interface{}) (string, error) {
	switch v2 := v.(type) {
	case time.Time:
		return v2.Format(p.layout), nil
	case int64:
		return time.Unix(v2, 0).Format(p.layout), nil
	case float64:
		return time.Unix(int64(v2), 0).Format(p.layout), nil
	default:
		return "", fmt.Errorf("not time: %+v (%T)", v, v)
	}
}
