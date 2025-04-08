package util

import (
	"io"
)

// DevZero /dev/zero
var DevZero = io.Reader(devZero(0))

type devZero int

func (p devZero) Read(b []byte) (n int, err error) {
	for i := range b {
		b[i] = 0
	}

	return len(b), nil
}
