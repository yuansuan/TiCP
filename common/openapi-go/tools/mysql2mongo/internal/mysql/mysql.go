package mysql

import (
	"context"
	"log"

	"github.com/yuansuan/ticp/common/openapi-go/tools/mysql2mongo/internal/model"

	"github.com/ory/viper"
	"xorm.io/xorm"
)

type Config struct {
	Uri string `yaml:"uri"`
}

func init() {
	viper.AutomaticEnv()
	_ = viper.BindEnv("mysql.uri", "MYSQL_URI")
}

func FindResidualsFromMysql(ctx context.Context, uri string) []model.Residual {
	// 连接MySQL
	engine, err := xorm.NewEngine("mysql", uri)
	if err != nil {
		log.Fatal(err)
	}

	var residuals []model.Residual
	err = engine.Find(&residuals)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(len(residuals), "residuals found in MySQL")

	return residuals
}
