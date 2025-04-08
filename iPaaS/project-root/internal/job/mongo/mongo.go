package mongo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ory/viper"
	"github.com/pkg/errors"
	"github.com/qiniu/qmgo"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

var (
	_client *qmgo.Client
)

type AllConfig struct {
	// ... other config
	Mongo Config `yaml:"mongo"`
}

type Config struct {
	Enable   bool   `yaml:"enable"`
	username string `yaml:"username"`
	password string `yaml:"password"`
	host     string `yaml:"host"`
	port     string `yaml:"port"`
	database string `yaml:"database"`
}

func init() {
	viper.AutomaticEnv()
	_ = viper.BindEnv("mongo.username", "MONGO_USERNAME")
	_ = viper.BindEnv("mongo.password", "MONGO_PASSWORD")
	_ = viper.BindEnv("mongo.host", "MONGO_HOST")
	_ = viper.BindEnv("mongo.port", "MONGO_PORT")
	_ = viper.BindEnv("mongo.database", "MONGO_DATABASE")
}

func (c *Config) Username() string {
	user := os.Getenv("MONGO_USERNAME")
	if user != "" {
		return user
	}
	return c.username
}

func (c *Config) Password() string {
	password := os.Getenv("MONGO_PASSWORD")
	if password != "" {
		return password
	}
	return c.password
}

func (c *Config) Host() string {
	host := os.Getenv("MONGO_HOST")
	if host != "" {
		return host
	}
	return c.host
}

func (c *Config) Port() string {
	port := os.Getenv("MONGO_PORT")
	if port != "" {
		return port
	}
	return c.port
}

func (c *Config) Database() string {
	db := os.Getenv("MONGO_DATABASE")
	if db != "" {
		return db
	}
	return c.database
}

func (c *Config) URI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", c.Username(), c.Password(), c.Host(), c.Port())
}

func Init(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	_client, err = qmgo.NewClient(ctx, &qmgo.Config{Uri: uri})
	if err != nil {
		logging.GetLogger(ctx).Errorf("Failed to connect to Mongo database: %v", err)
		return errors.Wrap(err, "Failed to connect to Mongo database")
	}

	logging.GetLogger(ctx).Info("Mongo Database connected")
	return nil
}

func Client() *qmgo.Client {
	return _client
}

func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := _client.Close(ctx); err != nil {
		logging.GetLogger(ctx).Errorf("Failed to close database client: %v", err)
		return
	}
	logging.GetLogger(ctx).Info("Database client closed")
}
