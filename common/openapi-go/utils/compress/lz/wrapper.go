package lz

import (
	"bytes"
	"io"
)

type Reader struct {
	r   io.ReadCloser
	buf *bytes.Buffer
}

func NewReader(r io.ReadCloser) (*Reader, error) {
	bs, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	data, err := DecompressInDouble(bs)
	if err != nil {
		return nil, err
	}

	return &Reader{
		r:   r,
		buf: bytes.NewBuffer(data),
	}, nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.buf.Read(p)
}

func (r *Reader) Close() error {
	return r.Close()
}
