package lz

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
)

func goZipInternal(data []byte, level int) ([]byte, error) {
	var buffer bytes.Buffer
	w, err := gzip.NewWriterLevel(&buffer, level)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(data); err != nil {
		w.Close()
		return nil, err
	}

	w.Close()
	return buffer.Bytes(), nil
}

// GoZip Gzip use golang
func GoZip(data []byte) ([]byte, error) {
	return goZipInternal(data, flate.DefaultCompression)
}

func GoUnZip(data []byte) ([]byte, error) {
	up, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer up.Close()

	bt, err := io.ReadAll(up)
	if err != nil {
		return nil, err
	}
	return bt, nil
}
