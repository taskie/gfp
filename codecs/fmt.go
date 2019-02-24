package codecs

import (
	"fmt"
	"strings"
)

func countFormatPercents(format string) int {
	return strings.Count(strings.Replace(format, "%%", "", -1), "%")
}

type FmtScannerConfig struct {
	Format string
}

type FmtScanner struct {
	format   string
	scanFunc func(s string) (interface{}, error)
	c        string
	n        int
}

func NewFmtScanner(config *FmtScannerConfig) (*FmtScanner, error) {
	format := config.Format
	n := countFormatPercents(format)
	if n != 1 {
		return nil, fmt.Errorf("invalid format: %d values", n)
	}
	c := format[len(format)-1:]
	var scanFunc func(s string) (interface{}, error)
	switch c {
	case "s", "c", "q":
		scanFunc = func(s string) (interface{}, error) {
			var v string
			fmt.Sscanf(s, format, &v)
			return v, nil
		}
	case "d", "b", "o", "X", "x", "U":
		scanFunc = func(s string) (interface{}, error) {
			var v int64
			fmt.Sscanf(s, format, &v)
			return v, nil
		}
	case "e", "E", "f", "g", "G":
		scanFunc = func(s string) (interface{}, error) {
			var v float64
			fmt.Sscanf(s, format, &v)
			return v, nil
		}
	case "t":
		scanFunc = func(s string) (interface{}, error) {
			var v bool
			fmt.Sscanf(s, format, &v)
			return v, nil
		}
	default:
		return nil, fmt.Errorf("invalid format: %s is unknown", c)
	}
	return &FmtScanner{
		format:   format,
		scanFunc: scanFunc,
		c:        c,
		n:        n,
	}, nil
}

func (sc *FmtScanner) Scan(s string) (interface{}, error) {
	return sc.scanFunc(s)
}

type FmtPrinterConfig struct {
	Format string
}

type FmtPrinter struct {
	format string
	c      string
	n      int
}

func NewFmtPrinter(config *FmtPrinterConfig) (*FmtPrinter, error) {
	format := config.Format
	n := countFormatPercents(format)
	if n != 1 {
		return nil, fmt.Errorf("invalid format: %d values", n)
	}
	c := format[len(format)-1:]
	switch c {
	case "s", "c", "q":
		break
	case "d", "b", "o", "X", "x", "U":
		break
	case "e", "E", "f", "g", "G":
		break
	case "t":
		break
	case "v": // allow %v
		break
	default:
		return nil, fmt.Errorf("invalid format: %s is unknown", c)
	}
	return &FmtPrinter{
		format: format,
		c:      c,
		n:      n,
	}, nil
}

func (p *FmtPrinter) Print(v interface{}) (string, error) {
	return fmt.Sprintf(p.format, v), nil
}
