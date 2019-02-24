package gfp

import (
	"fmt"

	"github.com/taskie/gfp/codecs"
)

type ScannerConfig struct {
	Type   string
	Format string
	Preset string
}

type Scanner interface {
	Scan(s string) (interface{}, error)
}

func NewScanner(config *ScannerConfig) (Scanner, error) {
	switch config.Type {
	case "fmt", "":
		return codecs.NewFmtScanner(&codecs.FmtScannerConfig{Format: config.Format})
	case "json":
		return codecs.NewJsonScanner(&codecs.JsonScannerConfig{})
	case "time":
		return codecs.NewTimeScanner(&codecs.TimeScannerConfig{Layout: config.Format, Preset: config.Preset})
	case "duration":
		return codecs.NewDurationScanner(&codecs.DurationScannerConfig{})
	default:
		return nil, fmt.Errorf("unknown typoe: %s", config.Type)
	}
}
