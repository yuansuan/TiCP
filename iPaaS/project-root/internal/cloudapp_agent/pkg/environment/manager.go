package environment

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/log"
)

const (
	lineBreaker      = "\n"
	readFileInterval = 1 * time.Second
	timeout          = 600 * time.Second
)

type Option interface {
	apply(c *config)
}

type optionFunc func(c *config)

func (f optionFunc) apply(c *config) {
	f(c)
}

func withDefaultOption() Option {
	return optionFunc(func(c *config) {
		//c.file = defaultCustomEnvPath
	})
}

func WithCustomEnvFile(file string) Option {
	return optionFunc(func(c *config) {
		if file == "" {
			return
		}

		c.file = file
	})
}

type config struct {
	file string
}

type CustomEnv struct {
	file string

	m map[string]string
}

func NewCustomEnv(opts ...Option) (*CustomEnv, error) {
	conf := &config{}
	withDefaultOption().apply(conf)
	for _, opt := range opts {
		opt.apply(conf)
	}

	customEnv := &CustomEnv{
		file: conf.file,
		m:    make(map[string]string),
	}

	if err := customEnv.init(); err != nil {
		return nil, fmt.Errorf("init custom env failed, %w", err)
	}

	return customEnv, nil
}

func (m *CustomEnv) init() error {
	content, err := m.readFileWithRetry()
	if err != nil {
		err = fmt.Errorf("read file %s with retry failed, %w", m.file, err)
		log.Error(err)
		return err
	}

	lines := strings.Split(string(content), lineBreaker)
	for _, line := range lines {
		kv := strings.Split(line, "=")
		if len(kv) < 2 {
			log.Warnf("parse %s to key=value failed", line)
			continue
		}
		key := kv[0]
		m.m[key] = line[len(key+"="):]
	}

	return nil
}

func (m *CustomEnv) Get(key string) string {
	return strings.Trim(m.m[key], "\r")
}

func (m *CustomEnv) readFileWithRetry() ([]byte, error) {
	tc := time.NewTicker(readFileInterval)
	to := time.NewTimer(timeout)

	for {
		select {
		case <-tc.C:
			content, err := os.ReadFile(m.file)
			if err != nil {
				log.Warnf("read file %s failed, %v", m.file, err)
				continue
			}

			return content, nil
		case <-to.C:
			return nil, fmt.Errorf("timout")
		}
	}
}
