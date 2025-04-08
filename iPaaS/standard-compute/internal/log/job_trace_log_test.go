package log

import (
	"context"
	"testing"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

func TestJobTraceLogger(t *testing.T) {
	ctx := WithTraceLoggerAndJobId(context.Background(), logging.Default(), 123)
	logger := GetJobTraceLogger(ctx)
	logger.Info("info")
	logger.Error("error")
}
