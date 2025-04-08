package compress

import (
	"compress/gzip"
	"io"
)

// Gzip compress as gzip
type Gzip struct {
}

// Compress compress as gzip
func (p Gzip) Compress(writer io.Writer) (io.WriteCloser, error) {
	return gzip.NewWriterLevel(writer, gzip.BestSpeed)
}

// Decompress decompress as gzip
func (p Gzip) Decompress(reader io.Reader) (io.Reader, error) {
	return gzip.NewReader(reader)
}
