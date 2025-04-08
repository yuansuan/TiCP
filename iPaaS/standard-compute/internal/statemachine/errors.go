package statemachine

import (
	"context"
	_errors "errors"

	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/util"
)

var (
	// ErrJobCanceled 表示当前任务已经被取消
	ErrJobCanceled = _errors.New("fsm: job canceled")
	// ErrJobUserFailed 表示当前任务在用户自定义规则中失败
	ErrJobUserFailed = _errors.New("fsm: job failed in custom rules")
)

// isJobCanceled 判断任务是否是已取消导致的
func isJobCanceled(err error) bool {
	return errors.Is(err, ErrJobCanceled) || util.IsCanceledError(err)
}

// isJobUserFailed 判断任务是否是用户自定义规则导致的失败
func isJobUserFailed(err error) bool {
	return errors.Is(err, ErrJobUserFailed)
}

// errWrap 将一些特定类型的错误转换为合适的错误类型
func errWrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) {
		return ErrJobCanceled
	}

	return errors.Wrap(err, msg)
}
