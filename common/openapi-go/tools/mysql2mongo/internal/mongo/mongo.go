package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/yuansuan/ticp/common/openapi-go/tools/mysql2mongo/internal/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/viper"
	"github.com/qiniu/qmgo"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

func init() {
	viper.AutomaticEnv()
	_ = viper.BindEnv("mongo.username", "MONGO_USERNAME")
	_ = viper.BindEnv("mongo.password", "MONGO_PASSWORD")
	_ = viper.BindEnv("mongo.host", "MONGO_HOST")
	_ = viper.BindEnv("mongo.port", "MONGO_PORT")
	_ = viper.BindEnv("mongo.database", "MONGO_DATABASE")
}

func (c *Config) URI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", c.Username, c.Password, c.Host, c.Port)
}

func InsertResidualsToMongoDB(ctx context.Context, uri, database string, residuals []model.Residual) {

	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: uri, Database: database, Coll: model.Residual{}.TableName()})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = cli.Close(ctx); err != nil {
			panic(err)
		}
	}()

	cnt := 0
	for _, residual := range residuals {
		_, err = cli.InsertOne(ctx, residual)
		if err != nil {
			log.Fatal(err)
		}
		cnt++
	}

	log.Println(cnt, "residuals inserted into MongoDB")
}
