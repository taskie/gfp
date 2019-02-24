package gfp

import (
	"fmt"

	"github.com/taskie/gfp/codecs"
)

type PrinterConfig struct {
	Type   string
	Format string
	Preset string
}

type Printer interface {
	Print(v interface{}) (string, error)
}

func NewPrinter(config *PrinterConfig) (Printer, error) {
	switch config.Type {
	case "fmt", "":
		return codecs.NewFmtPrinter(&codecs.FmtPrinterConfig{Format: config.Format})
	case "json":
		return codecs.NewJsonPrinter(&codecs.JsonPrinterConfig{})
	case "time":
		return codecs.NewTimePrinter(&codecs.TimePrinterConfig{Layout: config.Format, Preset: config.Preset})
	case "duration":
		return codecs.NewDurationPrinter(&codecs.DurationPrinterConfig{})
	default:
		return nil, fmt.Errorf("unknown type: %s", config.Type)
	}
}
