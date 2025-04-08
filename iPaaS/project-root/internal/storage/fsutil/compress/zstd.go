package compress

import (
	"github.com/klauspost/compress/zstd"
	"io"
)

// Zstd compress as zstd
type Zstd struct {
}

// Compress compress as zstd
func (p Zstd) Compress(writer io.Writer) (io.WriteCloser, error) {
	return zstd.NewWriter(writer)
}

// Decompress decompress as zstd
func (p Zstd) Decompress(reader io.Reader) (io.Reader, error) {
	return zstd.NewReader(reader)
}
