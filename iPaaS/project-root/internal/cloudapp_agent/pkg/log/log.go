package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger() *zap.Logger {
	if _logger == nil {
		return nil
	}

	return _logger.Desugar()
}

type LogLevel string

const (
	InfoLevel  LogLevel = "info"
	DebugLevel LogLevel = "debug"
)

func (l LogLevel) toZapLevel() zapcore.Level {
	return map[LogLevel]zapcore.Level{
		InfoLevel:  zap.InfoLevel,
		DebugLevel: zap.DebugLevel,
	}[l]
}

type ReleaseLevel string

const (
	DevelopmentLevel ReleaseLevel = "development"
	ProductionLevel  ReleaseLevel = "production"
)

type Logger struct {
	*zap.Logger
}

var _logger *zap.SugaredLogger

type Option interface {
	apply(c *config)
}

type optionFunc func(c *config)

func (f optionFunc) apply(c *config) {
	f(c)
}

func withDefaultOption() Option {
	return optionFunc(func(c *config) {
		c.logLevel = InfoLevel
		c.releaseLevel = ProductionLevel
		c.useConsole = false
		c.maxSize = 10
		c.maxAge = 30
		c.maxBackups = 10
		c.path = "defaultLogPath"
	})
}

func WithLogLevel(logLevel LogLevel) Option {
	return optionFunc(func(c *config) {
		if string(logLevel) == "" {
			return
		}

		c.logLevel = logLevel
	})
}

func WithReleaseLevel(releaseLevel ReleaseLevel) Option {
	return optionFunc(func(c *config) {
		if string(releaseLevel) == "" {
			return
		}

		c.releaseLevel = releaseLevel
	})
}

func WithUseConsole(useConsole bool) Option {
	return optionFunc(func(c *config) {
		if !useConsole {
			return
		}

		c.useConsole = useConsole
	})
}

func WithPath(path string) Option {
	return optionFunc(func(c *config) {
		if path == "" {
			return
		}

		c.path = path
	})
}

func WithMaxSize(maxSize int) Option {
	return optionFunc(func(c *config) {
		if maxSize == 0 {
			return
		}

		c.maxSize = maxSize
	})
}

func WithMaxAge(maxAge int) Option {
	return optionFunc(func(c *config) {
		if maxAge == 0 {
			return
		}

		c.maxAge = maxAge
	})
}

func WithMaxBackups(maxBackups int) Option {
	return optionFunc(func(c *config) {
		if maxBackups == 0 {
			return
		}

		c.maxBackups = maxBackups
	})
}

type config struct {
	logLevel     LogLevel
	releaseLevel ReleaseLevel

	useConsole bool
	path       string
	maxSize    int // MB
	maxAge     int // Day
	maxBackups int
}

func InitLogger(opts ...Option) error {
	c := &config{}
	withDefaultOption().apply(c)

	for _, opt := range opts {
		opt.apply(c)
	}

	zapOpts := []zap.Option{
		zap.WithCaller(true),
		zap.AddCallerSkip(1),
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

		zc.Level = zap.NewAtomicLevelAt(c.logLevel.toZapLevel())

		logger, err = zc.Build(zapOpts...)
		if err != nil {
			return err
		}
	} else {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   c.path,
			MaxSize:    c.maxSize,
			MaxAge:     c.maxAge,
			MaxBackups: c.maxBackups,
			Compress:   true,
		}

		var core zapcore.Core
		if c.releaseLevel == DevelopmentLevel {
			core = zapcore.NewCore(
				zapcore.NewConsoleEncoder(developEncoder()),
				zapcore.AddSync(lumberjackLogger),
				c.logLevel.toZapLevel(),
			)
		} else {
			core = zapcore.NewCore(
				zapcore.NewJSONEncoder(productionEncoder()),
				zapcore.AddSync(lumberjackLogger),
				c.logLevel.toZapLevel(),
			)
		}

		logger = zap.New(core, zapOpts...)
	}

	_logger = logger.Sugar()

	return nil
}

func (c *config) createCore(lumberjackLogger *lumberjack.Logger) zapcore.Core {
	if c.releaseLevel == ProductionLevel {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(productionEncoder()),
			zapcore.AddSync(lumberjackLogger),
			c.logLevel.toZapLevel(),
		)
	}

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(developEncoder()),
		zapcore.AddSync(lumberjackLogger),
		c.logLevel.toZapLevel(),
	)
}

func developEncoder() zapcore.EncoderConfig {
	return zap.NewDevelopmentEncoderConfig()
}

func productionEncoder() zapcore.EncoderConfig {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.ISO8601TimeEncoder

	return ec
}

func NewLogger(cfg config) (*Logger, error) {
	return nil, nil
}

func SetLogger(logger *Logger) {
	if logger != nil {
		_logger = logger.Sugar()
	}
}

func Info(args ...interface{}) {
	_logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	_logger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	_logger.Infow(msg, keysAndValues...)
}

func Debug(args ...interface{}) {
	_logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	_logger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	_logger.Debugw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	_logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	_logger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	_logger.Errorw(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	_logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	_logger.Warnf(template, args...)
}

func Warnw(template string, args ...interface{}) {
	_logger.Warnw(template, args...)
}

func Fatal(args ...interface{}) {
	_logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	_logger.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	_logger.Fatalw(msg, keysAndValues...)
}
