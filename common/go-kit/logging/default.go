package logging

var defaultLogger *Logger

type setDefaultConfig struct {
	logger *Logger
}

type SetDefaultOption interface {
	apply(c *setDefaultConfig)
}

type setDefaultOptionFunc func(c *setDefaultConfig)

func (f setDefaultOptionFunc) apply(c *setDefaultConfig) {
	f(c)
}

func WithLogger(logger *Logger) SetDefaultOption {
	return setDefaultOptionFunc(func(c *setDefaultConfig) {
		if logger == nil {
			return
		}

		c.logger = logger
	})
}

// SetDefault SetDefault
func SetDefault(opts ...SetDefaultOption) (*Logger, error) {
	c := &setDefaultConfig{}
	for _, opt := range opts {
		opt.apply(c)
	}

	var logger *Logger
	var err error
	if c.logger == nil {
		logger, err = NewLogger()
		if err != nil {
			panic(err)
		}
	} else {
		logger = c.logger
	}

	defaultLogger = logger
	return defaultLogger, err
}

// GetDefault GetDefault
func Default() *Logger {
	return defaultLogger
}
