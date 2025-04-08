package util

import (
	"context"
	"strings"

	"github.com/pkg/errors"
)

// IsCanceledError 判断是否是取消错误
func IsCanceledError(err error) bool {
	// terminate condition
	if err == nil {
		return false
	}

	// simple error
	if err == context.Canceled || errors.Is(err, context.Canceled) {
		return true
	}

	return IsCanceledError(errors.Unwrap(err)) || containsCanceledMessage(err)
}

// containsCanceledMessage 检测是否包含 "context canceled" 字符串
func containsCanceledMessage(err error) bool {
	return strings.Contains(err.Error(), context.Canceled.Error())
}
