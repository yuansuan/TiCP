package trace

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

const (
	RequestIdKey      = "x-ys-request-id"
	PathKey           = "path"
	RequestHeaderKey  = "request-header"
	RequestBodyKey    = "request-body"
	ResponseHeaderKey = "response-header"
	ResponseBodyKey   = "response-body"
)

type Logger struct {
	base *logging.Logger
}

func (l *Logger) Base() *logging.Logger {
	return l.base.WithOptions(zap.AddCallerSkip(-1))
}

func (l *Logger) Debug(args ...interface{}) {
	l.base.Debug(args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.base.Debugf(template, args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.base.Debugw(msg, keysAndValues...)
}

func (l *Logger) Info(args ...interface{}) {
	l.base.Info(args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.base.Infof(template, args...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.base.Infow(msg, keysAndValues...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.base.Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.base.Warnf(template, args...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.base.Warnw(msg, keysAndValues...)
}

func (l *Logger) Error(args ...interface{}) {
	l.base.Error(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.base.Errorf(template, args...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.base.Errorw(msg, keysAndValues...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.base.Panic(args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.base.Panicf(template, args...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.base.Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.base.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.base.Fatalf(template, args...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.base.Fatalw(msg, keysAndValues...)
}

func GetLogger(c context.Context) *Logger {
	logger := logging.GetLogger(c)

	return &Logger{
		base: logger.With(RequestIdKey, GetRequestId(c)).WithOptions(zap.AddCallerSkip(1)),
	}
}

func WithTraceLoggerAndId(ctx context.Context, traceLogger *Logger, requestId string) context.Context {
	return context.WithValue(
		context.WithValue(ctx, logging.LoggerName, traceLogger),
		RequestIdKey,
		requestId,
	)
}

func GetRequestId(c context.Context) string {
	switch c.(type) {
	case *gin.Context:
		ginCtx := c.(*gin.Context)
		requestId := ginCtx.Request.Header.Get(RequestIdKey)
		if requestId != "" {
			return requestId
		}

		requestIdV, exist := ginCtx.Get(RequestIdKey)
		if exist {
			ok := false
			requestId, ok = requestIdV.(string)
			if ok {
				return requestId
			}
		}

		return ""
	default:
		requestId, ok := c.Value(RequestIdKey).(string)
		if !ok {
			return ""
		}

		return requestId
	}
}
