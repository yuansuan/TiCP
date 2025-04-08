package middleware

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/common/openapi-go/utils/compress/lz"
)

const MaxLzCompressLength = 20 * 1024 * 1024 // 20MB

var ErrMaxLzCompressLength = errors.New("lz compress max length is 20MB")

func Lz() gin.HandlerFunc {
	return NewLzHandler().Handle
}

type lzWriter struct {
	gin.ResponseWriter
	buf    *bytes.Buffer
	length int
	exceed bool
}

func newLzWriter(w gin.ResponseWriter) *lzWriter {
	return &lzWriter{
		ResponseWriter: w,
		buf:            &bytes.Buffer{},
	}
}

func (l *lzWriter) WriteString(s string) (int, error) {
	if l.buf.Len() > MaxLzCompressLength {
		l.exceed = true
		return 0, ErrMaxLzCompressLength
	}

	l.Header().Del("Content-Length")
	return l.buf.Write([]byte(s))
}

func (l *lzWriter) Write(data []byte) (int, error) {
	if l.buf.Len() > MaxLzCompressLength {
		l.exceed = true
		return 0, ErrMaxLzCompressLength
	}

	l.Header().Del("Content-Length")
	return l.buf.Write(data)
}

func (l *lzWriter) WriteHeader(code int) {
	l.Header().Del("Content-Length")
	l.ResponseWriter.WriteHeader(code)
}

func (l *lzWriter) Close() error {
	if l.exceed {
		return errors.New("exceed lz max size(20MB)")
	}

	bs, err := lz.CompressInDouble(l.buf.Bytes())
	if err != nil {
		return err
	}
	l.length = len(bs)
	_, err = l.ResponseWriter.Write(bs)
	return err
}
