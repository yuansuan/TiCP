/*
 * // Copyright (C) 2018 LambdaCal Inc.
 *
 */

package logging

import (
	"context"
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	LoggerName        = "logger"
	DefaultMaxSize    = 500 // MB
	DefaultMaxAge     = 180 // DAY
	DefaultMaxBackups = 30  // COUNT
	DefaultLogPath    = "./custom.log"
)

func init() {
	logger, err := SetDefault()
	if err != nil {
		panic(fmt.Sprintf("set default logger failed, %v", err))
	}

	defaultLogger = logger
}

// Logger Logger
type Logger = zap.SugaredLogger

var atomLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
var originLevel zapcore.Level

// GetLogger GetLogger
func GetLogger(ctx context.Context) *Logger {
	logger, ok := ctx.Value(LoggerName).(*Logger)
	if ok && logger != nil {
		return logger
	}
	return defaultLogger
}

// AppendWith AppendWith
func AppendWith(ctx context.Context, kvs ...interface{}) context.Context {
	return context.WithValue(ctx, LoggerName, GetLogger(ctx).With(kvs...))
}

type LogConfig struct {
	logLevel     LogLevel
	releaseLevel ReleaseLevel

	useConsole bool
	path       string
	maxSize    int // MB
	maxAge     int // Day
	maxBackups int
}

func (c *LogConfig) createCore(lumberjackLogger *lumberjack.Logger) zapcore.Core {
	if c.releaseLevel == ProductionLevel {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(productionEncoder()),
			zapcore.AddSync(lumberjackLogger),
			c.logLevel.ToZapLevel(),
		)
	}

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(developEncoder()),
		zapcore.AddSync(lumberjackLogger),
		c.logLevel.ToZapLevel(),
	)
}

type LogConfigOption interface {
	apply(c *LogConfig)
}

type logConfigOptionFunc func(c *LogConfig)

func (f logConfigOptionFunc) apply(c *LogConfig) {
	f(c)
}

func WithDefaultLogConfigOption() LogConfigOption {
	return logConfigOptionFunc(func(c *LogConfig) {
		c.logLevel = InfoLevel
		c.releaseLevel = ProductionLevel
		c.useConsole = true
		c.maxSize = DefaultMaxSize
		c.maxAge = DefaultMaxAge
		c.maxBackups = DefaultMaxBackups
		c.path = DefaultLogPath
	})
}

// WithLogLevel 日志等级 [ info | debug ]
func WithLogLevel(logLevel LogLevel) LogConfigOption {
	return logConfigOptionFunc(func(c *LogConfig) {
		if logLevel.String() != "" {
			c.logLevel = logLevel
		}
	})
}

// WithReleaseLevel 日志发布等级 [ development | production ]
func WithReleaseLevel(releaseLevel ReleaseLevel) LogConfigOption {
	return logConfigOptionFunc(func(c *LogConfig) {
		if releaseLevel.String() != "" {
			c.releaseLevel = releaseLevel
		}
	})
}

// WithUseConsole 是否打印在console（即使用标准stdout/stderr），true: 打印在console上，false：打印在文件中
func WithUseConsole(useConsole bool) LogConfigOption {
	return logConfigOptionFunc(func(c *LogConfig) {
		if !useConsole {
			c.useConsole = useConsole
		}
	})
}

// WithMaxSize 单个日志文件最大限制 MB
func WithMaxSize(maxSize int) LogConfigOption {
	return logConfigOptionFunc(func(c *LogConfig) {
		if maxSize != 0 {
			c.maxSize = maxSize
		}
	})
}

// WithMaxAge 日志文件存档时间上限
func WithMaxAge(maxAge int) LogConfigOption {
	return logConfigOptionFunc(func(c *LogConfig) {
		if maxAge != 0 {
			c.maxAge = maxAge
		}
	})
}

// WithMaxBackups 日志文件存档最多个数
func WithMaxBackups(maxBackups int) LogConfigOption {
	return logConfigOptionFunc(func(c *LogConfig) {
		if maxBackups != 0 {
			c.maxBackups = maxBackups
		}
	})
}

// WithLogPath 日志文件路径
func WithLogPath(path string) LogConfigOption {
	return logConfigOptionFunc(func(c *LogConfig) {
		if path != "" {
			c.path = path
		}
	})
}

// NewLogger NewLogger
func NewLogger(opts ...LogConfigOption) (*Logger, error) {
	c := &LogConfig{}

	WithDefaultLogConfigOption().apply(c)
	for _, opt := range opts {
		opt.apply(c)
	}

	zapOpts := []zap.Option{
		zap.AddCaller(),
		zap.Development(),
		zap.AddStacktrace(zap.WarnLevel),
	}

	var logger *zap.Logger
	var err error
	if c.useConsole {
		var zc zap.Config
		if c.releaseLevel == DevelopmentLevel {
			zc = zap.NewDevelopmentConfig()
		} else {
			zc = zap.NewProductionConfig()
		}
		zc.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zc.Level = zap.NewAtomicLevelAt(c.logLevel.ToZapLevel())
		logger, err = zc.Build(zapOpts...)
		if err != nil {
			return nil, err
		}
	} else {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   c.path,
			MaxSize:    c.maxSize,
			MaxAge:     c.maxAge,
			MaxBackups: c.maxBackups,
			Compress:   true,
		}
		core := c.createCore(lumberjackLogger)
		logger = zap.New(core, zapOpts...)
	}

	return logger.Sugar(), nil
}

func developEncoder() zapcore.EncoderConfig {
	return zap.NewDevelopmentEncoderConfig()
}

func productionEncoder() zapcore.EncoderConfig {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.ISO8601TimeEncoder

	return ec
}

// SetLevel SetLevel
// Deprecated: 使用NewLogger(WithLogLevel())
func SetLevel(level int) {
	atomLevel.SetLevel(zapcore.Level(level))
}

// SetToDebug SetToDebug
// Deprecated: 使用NewLogger(WithLogLevel())
func SetToDebug() {
	atomLevel.SetLevel(zap.DebugLevel)
}

// Reverse Reverse
// Deprecated: 使用NewLogger(WithLogLevel())
func Reverse() {
	atomLevel.SetLevel(originLevel)
}

// IsTerminal IsTerminal
func IsTerminal() bool {
	return os.Getenv("TERM") != "dumb" && terminal.IsTerminal(int(os.Stdout.Fd()))
}

func Sync() error {
	return defaultLogger.Sync()
}

// Stack Stack
func Stack(skip int) []byte {
	return stack(skip)
}
