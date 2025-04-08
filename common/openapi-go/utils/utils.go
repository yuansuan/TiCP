package utils

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/openapi-go/common/compress"
	"io"
	"reflect"
	"sort"
	"time"
	"unsafe"
)

func Bytes(s string) (bs []byte) {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	hdr.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	hdr.Len = (*reflect.StringHeader)(unsafe.Pointer(&s)).Len
	hdr.Cap = hdr.Len
	return
}

func String(bs []byte) (s string) {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = (*reflect.SliceHeader)(unsafe.Pointer(&bs)).Data
	hdr.Len = (*reflect.SliceHeader)(unsafe.Pointer(&bs)).Len
	return
}

func Keys(values map[string]interface{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func SortEach(values map[string]interface{}, fn func(key string, val interface{})) {
	keys := Keys(values)
	sort.Strings(keys)

	for _, key := range keys {
		fn(key, values[key])
	}
}

func Stringify(data interface{}) string {
	if data == nil {
		return ""
	}

	fv := reflect.ValueOf(data)
	if fv.Kind() == reflect.Ptr && fv.IsNil() {
		return ""
	}

	if fv.Type().String() == "*time.Time" {
		return data.(*time.Time).Format(time.RFC3339)
	}

	return fmt.Sprintf("%v", data)
}

func CompressData(input interface{}, compressorType string) (io.Reader, error) {
	var data []byte
	switch v := input.(type) {
	case []byte:
		data = v
	case io.Reader:
		var err error
		data, err = io.ReadAll(v)
		if err != nil {
			return nil, errors.Errorf("read input error: %v", err)
		}
	default:
		return nil, errors.Errorf("unsupported input type: %T", input)
	}

	if compressorType == "" {
		return bytes.NewReader(data), nil
	}

	compressor, err := compress.GetCompressor(compressorType)
	if err != nil {
		return nil, errors.Errorf("get compressor error, err: %v", err)
	}

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		compressWriter, err := compressor.Compress(pw)
		if err != nil {
			pw.CloseWithError(errors.Errorf("compress error, err: %v", err))
			return
		}
		defer compressWriter.Close()

		if _, err := compressWriter.Write(data); err != nil {
			pw.CloseWithError(errors.Errorf("write error, err: %v", err))
			return
		}
	}()

	return pr, nil
}
