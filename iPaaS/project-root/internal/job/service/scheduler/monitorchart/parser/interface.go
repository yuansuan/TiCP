package parser

import "io"

// Parser parses the monitor chart
type Parser interface {
	Parse(file string, r io.Reader) (resMap map[string]*Result, err error)
}

// Result defines the Result
type Result struct {
	Key   string
	Items []*ParseItem
}

// ParseItem defines the ParseItem
type ParseItem struct {
	Key       string
	Iteration float64
	Value     float64
}
