package log

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
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
	apply(c *Config)
}

type optionFunc func(c *Config)

func (f optionFunc) apply(c *Config) {
	f(c)
}

func withDefaultOption(conf config.Log) Option {
	return optionFunc(func(c *Config) {
		c.logLevel = LogLevel(conf.Level)
		c.releaseLevel = ReleaseLevel(conf.ReleaseLevel)
		c.useConsole = conf.UseConsole
		c.maxSize = conf.MaxSize
		c.maxAge = conf.MaxAge
		c.maxBackups = conf.MaxBackups
		c.path = conf.Path
	})
}

func WithLogLevel(logLevel LogLevel) Option {
	return optionFunc(func(c *Config) {
		if string(logLevel) == "" {
			return
		}

		c.logLevel = logLevel
	})
}

func WithReleaseLevel(releaseLevel ReleaseLevel) Option {
	return optionFunc(func(c *Config) {
		if string(releaseLevel) == "" {
			return
		}

		c.releaseLevel = releaseLevel
	})
}

func WithUseConsole(useConsole bool) Option {
	return optionFunc(func(c *Config) {
		if !useConsole {
			return
		}

		c.useConsole = useConsole
	})
}

func WithPath(path string) Option {
	return optionFunc(func(c *Config) {
		if path == "" {
			return
		}

		c.path = path
	})
}

func WithMaxSize(maxSize int) Option {
	return optionFunc(func(c *Config) {
		if maxSize == 0 {
			return
		}

		c.maxSize = maxSize
	})
}

func WithMaxAge(maxAge int) Option {
	return optionFunc(func(c *Config) {
		if maxAge == 0 {
			return
		}

		c.maxAge = maxAge
	})
}

func WithMaxBackups(maxBackups int) Option {
	return optionFunc(func(c *Config) {
		if maxBackups == 0 {
			return
		}

		c.maxBackups = maxBackups
	})
}

type Config struct {
	logLevel     LogLevel
	releaseLevel ReleaseLevel

	useConsole bool
	path       string
	maxSize    int // MB
	maxAge     int // Day
	maxBackups int
}

func InitLogger(conf config.Log, opts ...Option) error {
	c := &Config{}
	withDefaultOption(conf).apply(c)

	for _, opt := range opts {
		opt.apply(c)
	}

	logger, err := logging.NewLogger(
		logging.WithLogLevel(logging.LogLevel(c.logLevel)),
		logging.WithReleaseLevel(logging.ReleaseLevel(c.releaseLevel)),
		logging.WithUseConsole(c.useConsole),
		logging.WithMaxAge(c.maxAge),
		logging.WithMaxSize(c.maxSize),
		logging.WithLogPath(c.path),
	)
	if err != nil {
		return err
	}

	logging.SetDefault(logging.WithLogger(logger))
	_logger = logger.WithOptions(zap.AddCallerSkip(1))

	return nil
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
