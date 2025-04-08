package conf_type

import (
	"time"

	"github.com/go-redis/redis"
)

// Redis Redis
type Redis struct {
	Builder func() *redis.Client
	Startup bool `yaml:"_startup"`

	// The network type, either tcp or unix.
	// Default is tcp.
	Network string `yaml:"network"`
	// host:port address.
	Addr string `yaml:"addr"`
	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string `yaml:"password"`
	// Database to be selected after connecting to the server.
	DB int `yaml:"db"`

	// Maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int `yaml:"max_retries"`
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.  8 * time.M
	MinRetryBackoff time.Duration `yaml:"min_retry_backoff"`
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration `yaml:"max_retry_backoff"`

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration `yaml:"dial_timeout"`
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration `yaml:"read_timeout"`
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration `yaml:"write_timeout"`

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int `yaml:"pool_size"`
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int `yaml:"min_idle_conns"`
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration `yaml:"max_conn_age"`
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration `yaml:"pool_timeout"`
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency time.Duration `yaml:"idle_check_frequency"`
}

// Redises Redises
type Redises map[string]Redis
