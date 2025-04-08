package middleware

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-pg/pg/v10"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	clientv3 "go.etcd.io/etcd/client/v3"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/multierr"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	conf_type "github.com/yuansuan/ticp/common/go-kit/gin-boot/conf-type"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware/tracing/otelxorm"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
)

// MySQL MySQL
type MySQL struct {
	DB  *sql.DB
	Orm *xorm.Engine
}

// Sqlite 记录SQLite类型的数据库
type Sqlite struct {
	Orm *xorm.Engine
}

// Middleware Middleware
type Middleware struct {
	conf *config.Config

	etcd struct {
		client  *clientv3.Client
		polling bool
	}

	mysqlLock *sync.Mutex
	mysqls    map[string]*MySQL

	pgsqlLock *sync.Mutex
	pgsqls    map[string]*pg.DB

	sqliteLock *sync.RWMutex
	sqlites    map[string]*Sqlite

	defaultDBType DatabaseType

	redisLock *sync.Mutex
	redises   map[string]*redis.Client

	cacheLock *sync.Mutex
	caches    map[string]ICache

	kafka *kafkaT

	eslock        *sync.Mutex
	elasticsearch map[string]*elasticsearch.Client

	logger *logging.Logger

	tracerProvider *tracesdk.TracerProvider
}

const defaultName = "default"

// Instance Instance
var Instance Middleware

// Init Init
func Init(config *config.Config, logger *logging.Logger) *Middleware {
	Instance.logger = logger
	Instance.conf = config
	Instance.initMysql()
	Instance.initPgsql()
	Instance.initSqlite()
	Instance.initRedis()
	Instance.initCache()
	Instance.initEtcd()
	Instance.initKafka()
	Instance.initElasticsearch()
	Instance.initOTEL()
	return &Instance
}

// Shutdown shutdown the middleware
func Shutdown() error {
	return Instance.Shutdown()
}

// DefaultSession return xorm session of `"default"` mysql
func DefaultSession(ctx context.Context) *xorm.Session {
	return Instance.DefaultSession(ctx)
}

// Session return xorm session of `name` mysql
func Session(ctx context.Context, name string) *xorm.Session {
	return Instance.Session(ctx, name)
}

// Shutdown flush and shutdown the middleware
func (mw *Middleware) Shutdown() (err error) {
	if mw.tracerProvider != nil {
		err = multierr.Append(err, mw.tracerProvider.Shutdown(context.Background()))
	}

	return err
}

// DefaultMysql DefaultMysql
func (mw *Middleware) DefaultMysql() *sql.DB {
	return mw.Mysql(defaultName)
}

func (mw *Middleware) mysql(name string) *MySQL {
	mw.mysqlLock.Lock()
	defer mw.mysqlLock.Unlock()
	if db, ok := mw.mysqls[name]; ok {
		return db
	} else {
		return nil
	}
}

// sqlite 获取一个Sqlite实例
func (mw *Middleware) sqlite(name string) *Sqlite {
	mw.sqliteLock.RLock()
	sqlite := mw.sqlites[name]
	mw.sqliteLock.RUnlock()

	return sqlite
}

// DefaultEtcd default etcd
func (mw *Middleware) DefaultEtcd() *clientv3.Client {
	return mw.etcd.client
}

// Mysql Mysql
func (mw *Middleware) Mysql(name string) *sql.DB {
	return mw.mysql(name).DB
}

// ORMEngine ORMEngine
func (mw *Middleware) ORMEngine(name string) *xorm.Engine {
	switch mw.defaultDBType {
	case DatabaseMySQL:
		return mw.mysql(name).Orm
	case DatabaseSQLite:
		return mw.sqlite(name).Orm
	default:
		return nil
	}
}

// DefaultORMEngine DefaultORMEngine
func (mw *Middleware) DefaultORMEngine() *xorm.Engine {
	return mw.ORMEngine(defaultName)
}

// SqliteORMEngine 获取Sqlite的XORM实例对象
func (mw *Middleware) SqliteORMEngine() *xorm.Engine {
	return mw.sqlite(defaultName).Orm
}

// Session Session
func (mw *Middleware) Session(ctx context.Context, name string) *xorm.Session {
	switch mw.defaultDBType {
	case DatabaseMySQL:
		return mw.mysql(name).Orm.Context(ctx)
	case DatabaseSQLite:
		return mw.sqlite(name).Orm.Context(ctx)
	default:
		return nil
	}
}

// DefaultSession DefaultSession
func (mw *Middleware) DefaultSession(ctx context.Context) *xorm.Session {
	return mw.Session(ctx, defaultName)
}

// Transaction Transaction
func (mw *Middleware) Transaction(ctx context.Context, name string, f func(*xorm.Session) (interface{}, error)) (interface{}, error) {
	return mw.mysql(name).Orm.Transaction(func(s *xorm.Session) (interface{}, error) {
		return f(s.Context(ctx))
	})
}

// DatabaseType 数据库类型
type DatabaseType int

const (
	DatabaseMySQL DatabaseType = iota
	DatabasePostgres
	DatabaseSQLite
)

// UseDefaultDatabase 设置使用哪个数据库, 默认为mysql
func (mw *Middleware) UseDefaultDatabase(dbt DatabaseType) {
	mw.defaultDBType = dbt
}

// DefaultTransaction DefaultTransaction
func (mw *Middleware) DefaultTransaction(ctx context.Context, f func(*xorm.Session) (interface{}, error)) (interface{}, error) {
	return mw.Transaction(ctx, defaultName, f)
}

func (mw *Middleware) pgsql(name string) *pg.DB {
	mw.pgsqlLock.Lock()
	defer mw.pgsqlLock.Unlock()
	if db, ok := mw.pgsqls[name]; ok {
		return db
	}
	return nil
}

// DefaultPgsql DefaultPgsql
func (mw *Middleware) DefaultPgsql() *pg.DB {
	return mw.pgsql(defaultName)
}

// DefaultRedis DefaultRedis
func (mw *Middleware) DefaultRedis() *redis.Client {
	return mw.Redis(defaultName)
}

// Redis Redis
func (mw *Middleware) Redis(name string) *redis.Client {
	mw.redisLock.Lock()
	defer mw.redisLock.Unlock()
	if db, ok := mw.redises[name]; ok {
		return db
	} else {
		return nil
	}
}

// DefaultCache DefaultCache
func (mw *Middleware) DefaultCache() ICache {
	return mw.Cache(defaultName)
}

// Cache Cache
func (mw *Middleware) Cache(name string) ICache {
	mw.cacheLock.Lock()
	defer mw.cacheLock.Unlock()
	if cache, ok := mw.caches[name]; ok {
		return cache
	} else {
		return nil
	}
}

// 初始化需要设置为启动的redis
func (mw *Middleware) initRedis() {
	mw.redisLock = &sync.Mutex{}
	mw.redises = make(map[string]*redis.Client)
	for dbName, dbConfig := range mw.conf.App.Middleware.Redis {
		if dbConfig.Startup == true {
			opt := mw.conf.App.Middleware.Redis[dbName]
			redisOptions := mw.initRedisOptions(&opt)
			mw.redises[dbName] = redis.NewClient(redisOptions)
		}
	}
}

// 初始化需要设置为启动的mysql
func (mw *Middleware) initMysql() {
	mw.mysqlLock = &sync.Mutex{}
	mw.mysqls = make(map[string]*MySQL)
	for dbName, dbConfig := range mw.conf.App.Middleware.Mysql {
		if dbConfig.Startup == true {
			var connect *sql.DB
			var err error
			var engine *xorm.Engine

			connect, err = sql.Open("mysql", dbConfig.Dsn)
			util.ChkErr(err)
			engine, err = xorm.NewEngine("mysql", dbConfig.Dsn)
			util.ChkErr(err)

			if dbConfig.MaxIdleConnection != 0 {
				connect.SetMaxIdleConns(dbConfig.MaxIdleConnection)
				engine.SetMaxIdleConns(dbConfig.MaxIdleConnection)
			}
			if dbConfig.MaxOpenConnection != 0 {
				connect.SetMaxOpenConns(dbConfig.MaxOpenConnection)
				engine.SetMaxOpenConns(dbConfig.MaxOpenConnection)
			}
			if dbConfig.MaxIdleTime != 0 {
				connect.SetConnMaxIdleTime(dbConfig.MaxIdleTime)
				engine.SetConnMaxLifetime(dbConfig.MaxIdleTime)
			}

			engine.ShowSQL(!dbConfig.HiddenSQL)
			engine.SetLogger(NewXormLogger(mw.logger, !dbConfig.HiddenSQL))
			mw.mysqls[dbName] = &MySQL{connect, engine}

			if mw.conf.App.Middleware.Tracing.Database.Enabled {
				engine.AddHook(otelxorm.NewTracingHook(&mw.conf.App.Middleware.Tracing))
			}
		}
	}
}

// initSqlite 初始化Sqlite数据库
func (mw *Middleware) initSqlite() {
	mw.sqliteLock = &sync.RWMutex{}
	mw.sqlites = make(map[string]*Sqlite)
	for name, cfg := range mw.conf.App.Middleware.Sqlite {
		if cfg.Startup {
			schema, rest, err := parseDSN(cfg.Dsn)
			util.ChkErr(err)

			engine, err := xorm.NewEngine(schema, rest)
			util.ChkErr(err)

			if cfg.MaxIdleConnection != 0 {
				engine.SetMaxIdleConns(cfg.MaxIdleConnection)
			}
			if cfg.MaxOpenConnection != 0 {
				engine.SetMaxOpenConns(cfg.MaxOpenConnection)
			}

			engine.ShowSQL(!cfg.HiddenSQL)
			engine.SetLogger(NewXormLogger(mw.logger, !cfg.HiddenSQL))
			mw.sqlites[name] = &Sqlite{Orm: engine}
		}
	}
}

// 初始化需要设置为启动的pgsql
func (mw *Middleware) initPgsql() {
	mw.pgsqlLock = &sync.Mutex{}
	mw.pgsqls = make(map[string]*pg.DB)
	for dbName, dbConfig := range mw.conf.App.Middleware.Pgsql {
		if dbConfig.Startup == true {
			opt, err := pg.ParseURL(dbConfig.URL)
			util.ChkErr(err)

			mw.pgsqls[dbName] = pg.Connect(opt)
		}
	}
}

func (mw *Middleware) initCache() {
	mw.cacheLock = &sync.Mutex{}
	mw.caches = make(map[string]ICache)
	for cacheName, cacheConfig := range mw.conf.App.Middleware.Cache {
		if cacheConfig.BackendType == conf_type.TypeCacheRedis {
			if mw.Redis(cacheConfig.Name) == nil {
				util.ChkErr(fmt.Errorf("cache not exist in redis, %v", cacheConfig.Name))
			}
			mw.caches[cacheName] = NewRedisCache(mw.Redis(cacheConfig.Name))
		}
	}
}

// 初始化redis.Options
func (mw *Middleware) initRedisOptions(opt *conf_type.Redis) *redis.Options {
	redisOptions := &redis.Options{
		Network:      opt.Network,
		Addr:         opt.Addr,
		Password:     opt.Password,
		DB:           opt.DB,
		PoolSize:     opt.PoolSize,
		MaxRetries:   opt.MaxRetries,
		MinIdleConns: opt.MinIdleConns,
	}
	if opt.MinRetryBackoff > 0 {
		redisOptions.MinRetryBackoff = opt.MinRetryBackoff * time.Millisecond
	}

	if opt.MaxRetryBackoff > 0 {
		redisOptions.MaxRetryBackoff = opt.MaxRetryBackoff * time.Millisecond
	}

	if opt.DialTimeout > 0 {
		redisOptions.DialTimeout = opt.DialTimeout * time.Millisecond
	}
	if opt.ReadTimeout > 0 {
		redisOptions.ReadTimeout = opt.ReadTimeout * time.Millisecond
	}
	if opt.WriteTimeout > 0 {
		redisOptions.WriteTimeout = opt.WriteTimeout * time.Millisecond
	}

	if opt.PoolTimeout > 0 {
		redisOptions.PoolTimeout = opt.PoolTimeout * time.Millisecond
	}

	if opt.MaxConnAge > 0 {
		redisOptions.MaxConnAge = opt.MaxConnAge * time.Millisecond
	}

	if opt.IdleTimeout > 0 {
		redisOptions.IdleTimeout = opt.IdleTimeout * time.Millisecond
	}

	if opt.IdleCheckFrequency > 0 {
		redisOptions.IdleCheckFrequency = opt.IdleCheckFrequency * time.Millisecond
	}

	return redisOptions
}

// parseDSN 解析DSN得到schema和剩余的字符
func parseDSN(raw string) (schema string, rest string, err error) {
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		switch {
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
		// do nothing
		case '0' <= c && c <= '9' || c == '+' || c == '-' || c == '.':
			if i == 0 {
				return "", raw, nil
			}
		case c == ':':
			if i == 0 {
				return "", "", errors.New("missing protocol scheme")
			}

			if strings.HasPrefix(raw[i+1:], "//") {
				return raw[:i], raw[i+3:], nil // trim "//"
			}

			return raw[:i], raw[i+1:], nil
		default:
			// we have encountered an invalid character,
			// so there is no valid scheme
			return "", raw, nil
		}
	}

	return "", raw, nil
}
