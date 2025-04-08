package daemon

import (
	"context"
	"fmt"

	"xorm.io/xorm"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/database"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/server"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/state"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/statemachine"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/systemuser"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/taskgroup"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/migrations"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/with"
)

type Option interface {
	apply(dc *daemonConfig)
}

type optionFunc func(dc *daemonConfig)

func (f optionFunc) apply(dc *daemonConfig) {
	f(dc)
}

func withDefaultOption() Option {
	return optionFunc(func(dc *daemonConfig) {
		dc.logConfig.path = "./log/standard-compute.log"
	})
}

func WithConfigPath(path string) Option {
	return optionFunc(func(dc *daemonConfig) {
		dc.confConfig.path = path
	})
}

func WithLogLevel(logLevel string) Option {
	return optionFunc(func(dc *daemonConfig) {
		dc.logConfig.logLevel = log.LogLevel(logLevel)
	})
}

func WithReleaseLevel(releaseLevel string) Option {
	return optionFunc(func(dc *daemonConfig) {
		dc.logConfig.releaseLevel = log.ReleaseLevel(releaseLevel)
	})
}

func WithLogUseConsole(useConsole bool) Option {
	return optionFunc(func(dc *daemonConfig) {
		dc.logConfig.useConsole = useConsole
	})
}

func WithLogPath(path string) Option {
	return optionFunc(func(dc *daemonConfig) {
		dc.logConfig.path = path
	})
}

func WithLogMaxSize(maxSize int) Option {
	return optionFunc(func(dc *daemonConfig) {
		dc.logConfig.maxSize = maxSize
	})
}

func WithLogMaxAge(maxAge int) Option {
	return optionFunc(func(dc *daemonConfig) {
		dc.logConfig.maxAge = maxAge
	})
}

func WithLogMaxBackups(maxBackups int) Option {
	return optionFunc(func(dc *daemonConfig) {
		dc.logConfig.maxBackups = maxBackups
	})
}

type daemonConfig struct {
	confConfig
	logConfig
}

type confConfig struct {
	path string
}

type logConfig struct {
	logLevel     log.LogLevel
	useConsole   bool
	releaseLevel log.ReleaseLevel
	path         string
	maxSize      int // MB
	maxAge       int // Day
	maxBackups   int
}

type Daemon struct {
	ctx context.Context
	dc  *daemonConfig

	conf      *config.Config
	db        *xorm.Engine
	taskGroup *taskgroup.TaskGroup

	httpserver   *server.Server
	jobScheduler backend.Provider
	factory      *statemachine.Factory
}

func New(ctx context.Context, opts ...Option) *Daemon {
	d := &Daemon{
		ctx: ctx,
		dc:  &daemonConfig{},
	}

	withDefaultOption().apply(d.dc)
	for _, opt := range opts {
		opt.apply(d.dc)
	}

	return d
}

func (d *Daemon) Wait() {
	d.taskGroup.Wait()
}

func (d *Daemon) Init() error {
	if err := d.init(); err != nil {
		log.Fatalf("init daemon failed, %v", err)
		return err
	}

	return nil
}

func (d *Daemon) init() error {
	var err error

	if err = d.initConfig(); err != nil {
		return fmt.Errorf("init config failed, %w", err)
	}

	if err = d.initLogger(); err != nil {
		return fmt.Errorf("init logger failed, %w", err)
	}

	if err = d.initDB(); err != nil {
		return fmt.Errorf("init db failed, %w", err)
	}

	if err = d.initSysUser(); err != nil {
		return fmt.Errorf("init sys user failed, %w", err)
	}

	if err = d.initTaskGroup(); err != nil {
		return fmt.Errorf("init task group failed, %w", err)
	}

	if err = d.initJobScheduler(); err != nil {
		return fmt.Errorf("init job scheduler failed, %w", err)
	}

	if err = d.initStateMachineFactory(); err != nil {
		return fmt.Errorf("init statemachine factory, %w", err)
	}

	if err = d.initHTTPServer(); err != nil {
		return fmt.Errorf("init http server failed, %w", err)
	}

	d.taskGroup.StartAll()

	return nil
}

func (d *Daemon) initLogger() error {
	if err := log.InitLogger(d.conf.Log,
		log.WithPath(d.dc.logConfig.path),
		log.WithUseConsole(d.dc.logConfig.useConsole),
		log.WithLogLevel(d.dc.logConfig.logLevel),
		log.WithReleaseLevel(d.dc.logConfig.releaseLevel),
		log.WithMaxAge(d.dc.logConfig.maxAge),
		log.WithMaxBackups(d.dc.logConfig.maxBackups),
		log.WithMaxSize(d.dc.logConfig.maxSize)); err != nil {
		return fmt.Errorf("init logger failed, %w", err)
	}

	return nil
}

func (d *Daemon) initConfig() error {
	conf, err := config.NewConfig(config.WithPath(d.dc.confConfig.path))
	if err != nil {
		return fmt.Errorf("new config failed, %w", err)
	}
	d.conf = conf

	return nil
}

func (d *Daemon) initTaskGroup() error {
	d.taskGroup = taskgroup.New(context.WithValue(d.ctx, with.OrmKey, d.db))
	return nil
}

func (d *Daemon) initDB() error {
	db, err := database.NewOrm(d.conf)
	if err != nil {
		return fmt.Errorf("new database orm failed, %w", err)
	}
	showSQL := !d.conf.Database.HiddenSQL
	db.ShowSQL(showSQL)
	db.SetLogger(log.NewXormLogger(log.GetLogger().Sugar(), showSQL))
	d.db = db

	migrationSource, err := migrations.NewSource(d.conf)
	if err != nil {
		return fmt.Errorf("new migration source failed, %w", err)
	}

	migration, err := database.NewMigration(d.db, d.conf, migrationSource)
	if err != nil {
		return fmt.Errorf("new migration failed, %w", err)
	}

	if err = migration.AutoMigrate(d.conf.Migrations.MigrationVersion); err != nil {
		return fmt.Errorf("auto migrate failed, %w", err)
	}

	daoInstance, err := dao.NewDao(d.conf.Snowflake)
	if err != nil {
		return err
	}
	dao.Default = daoInstance

	return nil
}

func (d *Daemon) initJobScheduler() error {
	jobScheduler, err := backend.NewProvider(d.conf.BackendProvider)
	if err != nil {
		return fmt.Errorf("backend new provider failed, %w", err)
	}

	d.jobScheduler = jobScheduler
	return nil
}

func (d *Daemon) initStateMachineFactory() error {
	factory, err := statemachine.NewFactory(d.conf, d.db, d.jobScheduler)
	if err != nil {
		return fmt.Errorf("new statemachine factory failed, %w", err)
	}

	d.factory = factory

	log.Info("starting recovery jobs...")
	if err = d.factory.RecoveryJobs(context.WithValue(d.ctx, with.OrmKey, d.db)); err != nil {
		return fmt.Errorf("recovery jobs failed, %w", err)
	}
	log.Info("jobs recovery done")

	if err = d.taskGroup.Add(d.factory); err != nil {
		return fmt.Errorf("add consume-jobs task to taskGroup failed, %w", err)
	}

	return nil
}

func (d *Daemon) initHTTPServer() error {
	hs := server.New(d.State())

	err := hs.Init()
	if err != nil {
		return fmt.Errorf("init http-server failed, %w", err)
	}

	if err = d.taskGroup.Add(hs); err != nil {
		return fmt.Errorf("add http-server to taskgroup failed, %w", err)
	}

	return nil
}

func (d *Daemon) State() *state.State {
	return &state.State{
		Conf:         d.conf,
		DB:           d.db,
		JobScheduler: d.jobScheduler,
		Factory:      d.factory,
	}
}

func (d *Daemon) initSysUser() error {
	return systemuser.Init()
}
