package leader

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/redsync.v1"
)

func newPools(name string) []redsync.Pool {
	pool := &redis.Pool{}

	opt := boot.Config.App.Middleware.Redis[name]

	// Other pool configuration not shown in this example.
	pool.IdleTimeout = opt.IdleTimeout
	pool.Wait = false
	pool.Dial = func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", opt.Addr)
		if err != nil {
			return nil, err
		}
		if opt.Password != "" {
			if _, err := c.Do("AUTH", opt.Password); err != nil {
				c.Close()
				return nil, err
			}
		}

		if _, err := c.Do("SELECT", opt.DB); err != nil {
			c.Close()
			return nil, err
		}
		return c, nil
	}
	return []redsync.Pool{pool}

}
