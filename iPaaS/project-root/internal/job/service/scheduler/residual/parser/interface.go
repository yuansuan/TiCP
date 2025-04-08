package parser

import (
	"io"
)

// Parser ...
type Parser interface {
	Parse(r io.Reader) (res *Result, err error)
}

// Result defines the Result
type Result struct {
	Series map[string][]float64
	XVar   []string
}
