package compress

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/project-root-api/common"
)

// Compressor : compress to writer, Decompress from reader
type Compressor interface {
	// Compress compress to writer, call Close() to flush content
	Compress(io.Writer) (io.WriteCloser, error)
	// Decompress Decompress from reader
	Decompress(io.Reader) (io.Reader, error)
}

func GetCompressor(name string) (Compressor, error) {
	switch name {
	case "", common.NONE:
		return &None{}, nil
	case common.GZIP:
		return &Gzip{}, nil
	case common.ZSTD:
		return &Zstd{}, nil
	}
	msg := fmt.Sprintf("no such compressor, name=%v", name)
	return nil, errors.Errorf(msg)
}
