package consts

import _errors "errors"

var (
	// ErrMalformedInstanceID 表示InstanceID格式不正确
	ErrMalformedInstanceID = _errors.New("malformed instance id")
	// ErrInstanceNotFound 表示找不到实例对象
	ErrInstanceNotFound = _errors.New("instance not found")
	// ErrInsufficientMachine 机器资源不足
	ErrInsufficientMachine = _errors.New("insufficient machine")
)
