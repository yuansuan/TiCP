package daemon

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/environment"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/httpserver"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/log"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/taskgroup"
)

type Option interface {
	apply(c *Config)
}

type optionFunc func(c *Config)

func (f optionFunc) apply(c *Config) {
	f(c)
}

type Config struct {
	hsConfig
	logConfig
	customEnvConfig
}

type hsConfig struct {
	useConsole bool
	address    string
}

type logConfig struct {
	useConsole bool
	logLevel   string
	logPath    string
}

type customEnvConfig struct {
	path string
}

func WithUseConsole(useConsole bool) Option {
	return optionFunc(func(c *Config) {
		c.hsConfig.useConsole = useConsole
		c.logConfig.useConsole = useConsole
	})
}

func WithLogLevel(logLevel string) Option {
	return optionFunc(func(c *Config) {
		c.logConfig.logLevel = logLevel
	})
}

func WithLogPath(logPath string) Option {
	return optionFunc(func(c *Config) {
		c.logConfig.logPath = logPath
	})
}

func WithHTTPServerAddress(address string) Option {
	return optionFunc(func(c *Config) {
		c.hsConfig.address = address
	})
}

func WithCustomEnvPath(path string) Option {
	return optionFunc(func(c *Config) {
		c.customEnvConfig.path = path
	})
}

type Daemon struct {
	config *Config
	tg     *taskgroup.TaskGroup

	customEnv  *environment.CustomEnv
	httpServer *httpserver.Server
}

func New(ctx context.Context, opts ...Option) *Daemon {
	config := new(Config)
	for _, opt := range opts {
		opt.apply(config)
	}

	return &Daemon{
		config: config,
		tg:     taskgroup.New(ctx),
	}
}

func (d *Daemon) Wait() {
	d.tg.Wait()
}

func (d *Daemon) Init() error {
	var err error
	if err = d.initLogger(); err != nil {
		return fmt.Errorf("init logger failed, %w", err)
	}

	if err = d.initCustomEnv(); err != nil {
		return fmt.Errorf("init custom env failed, %w", err)
	}
	log.Info("init custom env success")

	if err = d.resetPassword(); err != nil {
		err = fmt.Errorf("reset password failed, %w", err)
		log.Error(err)
		return err
	}
	log.Info("reset password success")

	if err = d.initHTTPServer(); err != nil {
		err = fmt.Errorf("init http server failed, %w", err)
		log.Error(err)
		return err
	}
	log.Info("init http server success")

	d.tg.StartAll()

	return nil
}

func (d *Daemon) initLogger() error {
	if err := log.InitLogger(
		log.WithUseConsole(d.config.logConfig.useConsole),
		log.WithLogLevel(log.LogLevel(d.config.logConfig.logLevel))); err != nil {
		return fmt.Errorf("init logger failed, %w", err)
	}

	return nil
}

func (d *Daemon) initCustomEnv() error {
	customEnv, err := environment.NewCustomEnv(
		environment.WithCustomEnvFile(d.config.customEnvConfig.path))
	if err != nil {
		return fmt.Errorf("new custom env failed, %w", err)
	}
	d.customEnv = customEnv

	return nil
}

func (d *Daemon) resetPassword() error {
	//return password.Reset(d.customEnv)
	return nil
}

func (d *Daemon) initHTTPServer() error {
	d.httpServer = httpserver.New(
		httpserver.WithAddress(d.config.address),
		httpserver.WithUseConsole(d.config.hsConfig.useConsole))
	if err := d.tg.Add(d.httpServer); err != nil {
		return fmt.Errorf("add task %s to task group failed, %w", d.httpServer.Name(), err)
	}

	return nil
}
