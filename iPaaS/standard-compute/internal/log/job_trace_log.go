package log

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"go.uber.org/zap"
)

const JobIdKey = "job-id"

func WithTraceLoggerAndJobId(ctx context.Context, logger *logging.Logger, jobId int64) context.Context {
	return context.WithValue(context.WithValue(ctx, logging.LoggerName, logger), JobIdKey, jobId)
}

func GetJobTraceLogger(ctx context.Context) *logging.Logger {
	jobId, _ := ctx.Value(JobIdKey).(int64)
	logger, ok := ctx.Value(logging.LoggerName).(*logging.Logger)
	if !ok {
		return logging.Default().With(zap.Int64(JobIdKey, jobId))
	}

	return logger.With(zap.Int64(JobIdKey, jobId))
}
