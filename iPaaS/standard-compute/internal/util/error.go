package util

import (
	"github.com/pkg/errors"
	"syscall"
)

var ClientReadError = errors.New("error reading in client")
var RollingHashReadError = errors.New("error reading in comparison")

type WrappedError struct {
	Msg string
	Err error
}

func (e WrappedError) Error() string {
	return e.Msg
}

func (e WrappedError) Unwrap() error {
	return e.Err
}

func IsFileMissingCausedError(err error) bool {
	return errors.Is(err, syscall.ENOENT) || // no such file or directory, when locale fs read, or openapi writeAt read source data which throws a http error from net/http, while file is missing
		errors.Is(err, RollingHashReadError) || // when locale fs read specific block while file is missing
		errors.Is(err, ClientReadError)
}
