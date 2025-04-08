package boot

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/fs"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/spf13/viper"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/file"
	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
	_http "github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/internal/cmd"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware/tracing/otelgin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/migration"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/monitor"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/registry"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
	_ "github.com/yuansuan/ticp/common/go-kit/gin-boot/validation"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/version"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"go.elastic.co/apm"
	"go.elastic.co/apm/transport"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/reflection"
)

func init() {
	// check go version
	if versionAllow, err := util.GoVersionCompare(runtime.Version(), "go1.12"); !versionAllow {
		if err != nil {
			util.ChkErr(err)
		}
		util.ChkErr(errors.New("check that the go version is greater than 1.12"))
	}
}

var (
	server       *Server
	logger       *logging.Logger
	Env          env.EnvType
	Config       config.Config
	Registry     *registry.Registry
	configLoaded bool
	Version      string
)

// Recovery Recovery
func Recovery() {
	if r := recover(); r != nil {
		logger.Fatalw(string(logging.Stack(5)), "panic", r)
	}
}

var mut sync.Mutex

var grpcServerReadyWG sync.WaitGroup

var configPath string

var customLogLevel string

// DefaultServer ...
func DefaultServer(_configPath, _customLogLevel string) *Server {
	if _configPath != "" {
		configPath = _configPath
	}

	if _customLogLevel != "" {
		customLogLevel = _customLogLevel
	}

	return Default()
}

func LoadConfigAndLogger() error {
	if configLoaded {
		return nil
	}

	// import docker secrets as env
	env.ImportDockerSecretAsEnv()

	// init env
	env.InitEnv(`.env`)
	Env = *env.Env
	modeStr := env.ModeName(env.Env.Mode)
	logMode := "mode: " + modeStr

	// init config
	if configPath == "" {
		configPath = config.ConfigDir
	}
	configFilePath := configPath + string(os.PathSeparator) + modeStr + ".yml"
	logConfig := "config file: " + configFilePath

	config.InitConfig(configFilePath)
	{
		rewriteFile := path.Join(configPath, config.RewriteFileName)
		if f, e := os.Open(rewriteFile); e == nil {
			defer f.Close()

			d := yaml.NewDecoder(f)
			util.ChkErr(d.Decode(config.Conf))
		}
	}

	err := config.ReadConfig(modeStr+"_custom", configPath, "yaml")
	if err != nil {
		// _custom file is optional, ignore not found
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	{
		v := viper.New()
		v.SetConfigName(strings.TrimSuffix(config.RewriteFileName, filepath.Ext(config.RewriteFileName)))
		v.AddConfigPath(configPath)
		v.SetConfigType("yaml")
		err := v.ReadInConfig()
		if err == nil {
			viper.MergeConfigMap(v.AllSettings())
		}
	}

	Config = *config.Conf

	//init APM EVN Config
	if config.Conf.App.Middleware.APM.APMStartUp {
		os.Setenv("ELASTIC_APM_SERVER_URL", config.Conf.App.Middleware.APM.APMServerURL)
		apm.DefaultTracer.Transport, _ = transport.InitDefault()
	}

	// init logger
	logLevel := env.Env.LogLevel
	if customLogLevel != "" {
		logLevel = env.LogLevelMap[customLogLevel]
	}

	opts := make([]logging.LogConfigOption, 0)
	if logLevel == int(zapcore.DebugLevel) {
		opts = append(opts, logging.WithLogLevel(logging.DebugLevel))
		opts = append(opts, logging.WithReleaseLevel(logging.DevelopmentLevel))
	}

	if config.Conf.App.Middleware.Logger.UseFile {
		opts = append(opts,
			logging.WithUseConsole(false),
			logging.WithLogPath(ensureLogPath()),
			logging.WithMaxAge(config.Conf.App.Middleware.Logger.MaxAge),
			logging.WithMaxSize(config.Conf.App.Middleware.Logger.MaxSize),
			logging.WithMaxBackups(config.Conf.App.Middleware.Logger.MaxBackups),
		)
	}

	logger, err = logging.NewLogger(opts...)
	if err != nil {
		return fmt.Errorf("new logger failed, %w", err)
	}
	if _, err = logging.SetDefault(logging.WithLogger(logger)); err != nil {
		return fmt.Errorf("set default logger failed, %w", err)
	}

	// the logger.log_dir is optional
	if d := config.Conf.App.Middleware.Logger.LogDir; len(d) != 0 {
		file.TouchDir(d)
		_ = flag.Set("log_dir", d)
	}

	if version.ShouldLogVersion() {
		version.LogVersion()
		os.Exit(0)
	}

	logger.Debug("logLevel: " + env.LogLevelName(env.Env.LogLevel))
	logger.Debug(logMode, logConfig)
	configLoaded = true
	return nil
}

// Default Default
func Default() *Server {
	defer Recovery()
	if server != nil {
		return server
	}
	mut.Lock()
	defer mut.Unlock()

	//doubling-check
	if server != nil {
		return server
	}

	if err := LoadConfigAndLogger(); err != nil {
		util.ChkErr(err)
	}

	// init middlewareType
	logger.Debug("initial middlewareType")
	middleware.Init(config.Conf, logger)

	// initGRPCClient
	initGRPCClient()

	// init registry
	Registry = registry.GetRegistry()

	// init gin mode
	gin.SetMode(env.GinMode(env.Env.Mode))

	// init server
	server = &Server{
		waitGroup: &sync.WaitGroup{},
	}
	if server.RunningInConsole() {
		ctx := kong.Parse(&cmd.Cli)
		// Call the Run() method of the selected parsed command.
		err := ctx.Run(cmd.Context{AppName: config.Conf.App.Name})
		ctx.FatalIfErrorf(err)
		return server
	}
	server.Register(initGRPCServer, initMonitor)
	server.Driver = gin.New()
	server.Driver.Use(
		otelgin.Middleware(config.Conf.App.Name,
			otelgin.WithEnabled(config.Conf.App.Middleware.Tracing.Startup),
			otelgin.WithExcludes(append(config.Conf.App.Middleware.Tracing.Http.Excludes, "/metrics", "/debug")...)),
		middleware.GinLogger(append(config.Conf.App.Middleware.HTTP.Logger.Excludes, "/debug/pprof/cmdline", "/metrics")...),
		logging.GinRecovery(),
	)

	server.OnShutdown(func(server *_http.Driver) {
		util.ChkErr(middleware.Shutdown())
	})
	return server
}

// Server Server
type Server struct {
	_http.IServer
	Driver             *_http.Driver
	handlers           []Handler
	backgroundHandlers []Handler
	onShutdownHandlers []Handler
	waitGroup          *sync.WaitGroup
}

// Handler Handler
type Handler func(server *_http.Driver)

// Register Register
func (s *Server) Register(handlers ...Handler) *Server {
	s.handlers = append(s.handlers, handlers...)
	return s
}

// RegisterRoutine RegisterRoutine
func (s *Server) RegisterRoutine(handlers ...Handler) *Server {
	for _, h := range handlers {
		hh := h
		s.backgroundHandlers = append(s.backgroundHandlers, func(drv *_http.Driver) {
			s.waitGroup.Add(1)
			go func(wg *sync.WaitGroup, handler Handler) {
				defer Recovery()
				defer wg.Done()
				handler(drv)
			}(s.waitGroup, hh)
		})
	}
	return s
}

// OnShutdown OnShutdown
func (s *Server) OnShutdown(handlers ...Handler) *Server {
	s.onShutdownHandlers = append(s.onShutdownHandlers, handlers...)
	return s
}

func (s *Server) DBAutoMigrate(migrationSourceFS fs.FS) *Server {
	if migrationSourceFS == nil {
		panic("migration source fs is nil")
	}

	if err := doMigration(migrationSourceFS); err != nil {
		panic(fmt.Sprintf("migration failed, %v", err))
	}
	logger.Info("database auto migrate success")

	return s
}

// RunningInConsole Determine if the application is running in the console.
func (s *Server) RunningInConsole() bool {
	args := os.Args
	return len(args) >= 2 && args[1] == "artisan"
}

// Run Run
func (s *Server) Run() {
	defer Recovery()
	if s.RunningInConsole() {
		return
	}
	httpAddr := config.Conf.App.Host + ":" + strconv.Itoa(config.Conf.App.Port)
	logger.Debug("http server listening on : " + httpAddr)
	s.Register(grpc_boot.InitGrpcGateway, startGRPCServer)
	for _, h := range s.handlers {
		h(s.Driver)
	}
	for _, h := range s.backgroundHandlers {
		h(s.Driver)
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: s.Driver,
	}

	s.ensurePprofHandlers()

	// Wait for the gRPC server ready
	grpcServerReadyWG.Wait()
	grpc_prometheus.Register(grpc_boot.Server().Driver())

	go func(server *http.Server) {
		defer Recovery()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.ChkErr(err)
		}
	}(server)

	sigQuit := <-shutdown

	logger.Infof("http server gracefully shutdown by signal %v", sigQuit.String())

	//wait for routines finishing
	for _, h := range s.onShutdownHandlers {
		go h(s.Driver)
	}
	util.WaitTimeout(s.waitGroup, 2*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}

	logger.Infof("gRPC server %s gracefully shutdown by signal %v", config.Conf.App.Name, sigQuit.String())
	wg.Add(1)
	go func(s *grpc_boot.ServerType, name string, wg *sync.WaitGroup) {
		defer Recovery()
		defer wg.Done()
		s.Driver().GracefulStop()
	}(grpc_boot.Server(), config.Conf.App.Name, wg)
	wg.Wait()

	server.Shutdown(ctx)
}

func (s *Server) ensurePprofHandlers() {
	if shouldPprofRegisterOnBusinessHandler() {
		pprof.Register(s.Driver)
	} else {
		go func() {
			defer Recovery()

			pprofHandler := gin.Default()
			pprof.Register(pprofHandler)
			pprofSrv := &http.Server{
				Addr:    fmt.Sprintf("%s:%d", config.Conf.App.PprofHost, config.Conf.App.PprofPort),
				Handler: pprofHandler,
			}

			if err := pprofSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				util.ChkErr(err)
			}
		}()
	}
}

func shouldPprofRegisterOnBusinessHandler() bool {
	// pporf host port 与业务 host port 相同 或者 pprof相关配置均为空，注册到业务的handler中
	return (config.Conf.App.PprofHost == config.Conf.App.Host && config.Conf.App.PprofPort == config.Conf.App.Port) ||
		(config.Conf.App.PprofHost == "" && config.Conf.App.PprofPort == 0)
}

func initMonitor(server *_http.Driver) {
	if config.Conf.App.Middleware.Monitor.StartUp != true {
		return
	}

	monitor.Use(&monitor.Config{
		Server:     server,
		ListenAddr: config.Conf.App.Middleware.Monitor.ListenAddr,
		MetricPath: config.Conf.App.Middleware.Monitor.MetricPath,
	})
}

func initGRPCServer(server *_http.Driver) {
	if config.Conf.App.Middleware.GRPC.Server.Default.Startup_ {
		addr := config.Conf.App.Middleware.GRPC.Server.Default.Addr
		log.Printf("Init GRPC listener: %s\n", addr)
		err := grpc_boot.InitListener(addr)
		util.ChkErr(err)
	}
	grpc_boot.InitServer(&config.Conf.App.Middleware.GRPC.Server.Default)
}

func initGRPCClient() {
	grpc_boot.InitClient(&config.Conf.App.Middleware.GRPC.Client)
}

func startGRPCServer(server *_http.Driver) {
	if config.Conf.App.Middleware.GRPC.Server.Default.Startup_ {
		grpcServerReadyWG.Add(1)

		go func(name string, s *grpc_boot.ServerType) {
			defer Recovery()
			d := s.Driver()
			reflection.Register(d)
			// After registering, invoke wait group done
			grpcServerReadyWG.Done()

			logger.Infof("starting %v\n", name)
			if err := d.Serve(s.Listener()); err != nil {
				log.Fatalf("booting of grpc server is failed, name is %v, %v", name, err)
			}
		}(config.Conf.App.Name, grpc_boot.Server())
	}
}

const (
	defaultLogDir      = "log"
	defaultLogFilename = "custom.log"
)

func ensureLogPath() string {
	logDir := defaultLogDir
	if config.Conf.App.Middleware.Logger.LogDir != "" {
		logDir = config.Conf.App.Middleware.Logger.LogDir
	}

	logFilename := defaultLogFilename
	if config.Conf.App.Name != "" {
		logFilename = fmt.Sprintf("%s.log", config.Conf.App.Name)
	}

	return filepath.Join(logDir, logFilename)
}

func doMigration(migrationSourceFS fs.FS) error {
	// skip
	if !config.Conf.App.DBMigration.AutoMigrate {
		logger.Info("database auto migration skipped")
		return nil
	}
	logger.Infof("starting database auto migration")

	if migrationSourceFS == nil {
		return fmt.Errorf("auto migration enabled but migration source fs is nil")
	}

	defaultMysql, exist := config.Conf.App.Middleware.Mysql["default"]
	if !exist {
		return fmt.Errorf("cannot find default mysql in config")
	}

	if config.Conf.App.Name == "" {
		return fmt.Errorf("app.name cannot be empty")
	}

	forceMigrate := true
	if config.Conf.App.DBMigration.ForceMigrate != nil {
		forceMigrate = *config.Conf.App.DBMigration.ForceMigrate
	}

	return migration.Migrate(defaultMysql.Dsn, config.Conf.App.Name, config.Conf.App.DBMigration.Version,
		migrationSourceFS, forceMigrate)
}
