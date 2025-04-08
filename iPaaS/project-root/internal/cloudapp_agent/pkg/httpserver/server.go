package httpserver

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/log"
)

type Option interface {
	apply(c *Config)
}

type optionFunc func(c *Config)

func (f optionFunc) apply(c *Config) {
	f(c)
}

type Config struct {
	useConsole bool
	address    string
}

func withDefaultOption() Option {
	return optionFunc(func(c *Config) {
		c.useConsole = false
		c.address = ":3390"
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

func WithAddress(address string) Option {
	return optionFunc(func(c *Config) {
		if address == "" {
			return
		}

		c.address = address
	})
}

type Server struct {
	config *Config
	base   *http.Server
}

func New(opts ...Option) *Server {
	config := new(Config)

	withDefaultOption().apply(config)
	for _, opt := range opts {
		opt.apply(config)
	}

	s := &Server{
		config: config,
		base: &http.Server{
			Addr: config.address,
		},
	}
	s.init()

	return s
}

func (s *Server) Name() string {
	return "http-server"
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.base.ListenAndServe(); err != nil {
			log.Warnf("listen and serve http server failed, %v", err)
		}
	}()

	<-ctx.Done()
	return s.stop()
}

func (s *Server) stop() error {
	ctxImmediate, cancel := context.WithTimeout(context.Background(), -1)
	defer cancel()

	return s.base.Shutdown(ctxImmediate)
}

func (s *Server) init() {
	handler := gin.New()

	if s.config.useConsole {
		handler.Use(gin.Logger())
	} else {
		lumberjackLogger := &lumberjack.Logger{
			//Filename:   defaultAccessLogPath,
			MaxSize:    10,
			MaxAge:     30,
			MaxBackups: 10,
			Compress:   true,
		}

		handler.Use(gin.LoggerWithWriter(lumberjackLogger))
	}
	handler.Use(gin.Recovery())

	for _, endpoint := range api.GetEndpoints() {
		handler.Handle(endpoint.Method, endpoint.RelativePath, endpoint.Handler)
	}

	s.base.Handler = handler
}
