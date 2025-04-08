package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/server/api"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/server/middleware"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/state"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/statemachine"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/util"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/with"
)

type Server struct {
	s  *state.State
	r  *gin.Engine
	db *xorm.Engine

	businessSrv    *http.Server
	performanceSrv *http.Server

	stateMachineFactory *statemachine.Factory
}

func (s *Server) Name() string {
	return "http-server"
}

func (s *Server) Start(ctx context.Context) error {
	errCh1 := make(chan error, 1)
	go func() {
		if s.businessSrv != nil {
			if err := s.businessSrv.ListenAndServe(); err != nil {
				log.Errorf("run business server failed, %v", err)
				errCh1 <- err

				return
			}
		}
	}()

	errCh2 := make(chan error, 1)
	go func() {
		if s.performanceSrv != nil {
			if err := s.performanceSrv.ListenAndServe(); err != nil {
				log.Errorf("run performance server failed, %v", err)
				errCh2 <- err

				return
			}
		}
	}()

	select {
	case err := <-errCh1:
		return fmt.Errorf("start business server failed, %w", err)
	case err := <-errCh2:
		return fmt.Errorf("start performance server failed, %w", err)
	case <-ctx.Done():
	}

	// graceful shutdown
	ctxImmediate, cancel := context.WithTimeout(context.Background(), -1)
	defer cancel()

	var err error
	if s.businessSrv != nil {
		if err = s.businessSrv.Shutdown(ctxImmediate); err != nil {
			return fmt.Errorf("shutdown business server failed, %w", err)
		}
	}

	if s.performanceSrv != nil {
		if err = s.performanceSrv.Shutdown(ctxImmediate); err != nil {
			return fmt.Errorf("shutdown performance server failed, %w", err)
		}
	}

	return nil
}

func (s *Server) Run() error {
	return s.businessSrv.ListenAndServe()
}

func (s *Server) injectorForBusinessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// inject db
		ctx := context.WithValue(c.Request.Context(), with.OrmKey, s.s.DB)
		c.Request = c.Request.WithContext(ctx)

		// inject state
		c.Set(util.StateKeyInGinCtx, s.s)
	}
}

func New(s *state.State) *Server {
	srv := &Server{
		s: s,
	}

	return srv
}

func (s *Server) Init() error {
	var err error
	if err = s.initBusinessSrv(); err != nil {
		return fmt.Errorf("init business server failed, %w", err)
	}

	if err = s.initPerformanceSrv(); err != nil {
		return fmt.Errorf("init performance server failed, %w", err)
	}

	return nil
}

func (s *Server) initBusinessSrv() error {
	if s.s.Conf.HttpAddress == "" {
		return fmt.Errorf("http_address is empty")
	}

	s.r = gin.New()
	s.businessSrv = &http.Server{
		Addr:    s.s.Conf.HttpAddress,
		Handler: s.r,
	}

	s.initEngine(s.r)

	s.initSystemAPI()

	return nil
}

func (s *Server) initSystemAPI() {
	systemGroup := s.r.Group("/system")
	s.registerBusinessMiddleware(systemGroup)

	systemGroup.Handle(http.MethodPost, "/jobs", api.PostJobs)
	systemGroup.Handle(http.MethodGet, "/jobs/:JobID", api.GetJob)
	systemGroup.Handle(http.MethodGet, "/jobs", api.GetJobs)
	systemGroup.Handle(http.MethodPost, "/jobs/:JobID/cancel", api.CancelJob)
	systemGroup.Handle(http.MethodDelete, "/jobs/:JobID", api.DeleteJob)
	systemGroup.Handle(http.MethodGet, "/resource", api.GetResource)
	systemGroup.Handle(http.MethodGet, "/jobs/:JobID/cpuusage", api.GetCpuUsage)
	systemGroup.Handle(http.MethodPost, "/command", api.PostCommand)
	systemGroup.Handle(http.MethodGet, "/healthz", api.Health)
}

func (s *Server) registerBusinessMiddleware(group *gin.RouterGroup) {
	group.Use(s.injectorForBusinessHandler())
	group.Use(middleware.SignatureValidator)
	group.Use(middleware.UrlValidator)
	group.Use(middleware.IngressLogger)
}

func (s *Server) initEngine(engine *gin.Engine) {
	if s.s.Conf.AccessLog.UseConsole {
		engine.Use(gin.Logger())
	} else {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   s.s.Conf.AccessLog.Path,
			MaxSize:    s.s.Conf.AccessLog.MaxSize,
			MaxAge:     s.s.Conf.AccessLog.MaxAge,
			MaxBackups: s.s.Conf.AccessLog.MaxBackups,
			Compress:   true,
		}

		engine.Use(gin.LoggerWithWriter(lumberjackLogger))
	}
	engine.Use(gin.Recovery())
}

func (s *Server) initPerformanceSrv() error {
	if s.s.Conf.PerformanceAddress == "" {
		// 隔离pprof handler与业务handler，此处没有必要注册一些例如校验签名的中间件。
		pprofGroup := s.r.Group("/debug/pprof")
		pprof.RouteRegister(pprofGroup, "")

		return nil
	}

	performanceEngine := gin.New()
	s.performanceSrv = &http.Server{
		Addr:    s.s.Conf.PerformanceAddress,
		Handler: performanceEngine,
	}
	s.initEngine(performanceEngine)
	pprof.Register(performanceEngine)

	return nil
}
