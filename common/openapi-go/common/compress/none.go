package compress

import (
	"io"
)

// None not compress
type None struct {
}

type nopCloser struct {
	io.Writer
}

func (p *nopCloser) Close() error {
	return nil
}

// Compress not compress
func (p None) Compress(writer io.Writer) (io.WriteCloser, error) {
	return &nopCloser{writer}, nil
}

// Decompress not Decompress
func (p None) Decompress(reader io.Reader) (io.Reader, error) {
	return reader, nil
}
