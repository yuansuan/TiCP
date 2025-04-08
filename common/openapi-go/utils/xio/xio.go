package xio

import (
	"bytes"
	"io"
)

// DupReader reads all of r to memory and then returns two equivalent
// ReadClosers yielding the same bytes
func DupReader(r io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	if r == nil {
		return r, r, nil
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, nil, err
	}

	if buf.Len() == 0 {
		return nil, nil, nil
	}

	if err := r.Close(); err != nil {
		return nil, nil, err
	}

	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// SeekLen returns the size of the io.Seeker
func SeekLen(s io.Seeker) (int64, error) {
	curr, err := s.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	end, err := s.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = s.Seek(curr, io.SeekStart)
	return end, err
}
