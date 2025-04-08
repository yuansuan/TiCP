package trace

import (
	"context"
	"testing"
)

func TestTraceLogger(t *testing.T) {
	ctx := WithTraceLoggerAndId(context.Background(), GetLogger(context.TODO()), "requestId-1")
	traceLogger := GetLogger(ctx)
	traceLogger.Info("info")
}
