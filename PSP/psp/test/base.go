package test

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"xorm.io/xorm"

	mainconfig "github.com/yuansuan/ticp/PSP/psp/cmd/config"
	"github.com/yuansuan/ticp/PSP/psp/test/consts"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

var (
	ConfigFilePath string
)

func InitTestEnv() context.Context {
	// 设置配置文件路径
	SetConfigFilePath()

	// 注册服务
	config.InitConfig(ConfigFilePath)
	middleware.Init(config.Conf, logging.Default())
	grpc_boot.InitClient(&config.Conf.App.Middleware.GRPC.Client)

	db, err := xorm.NewEngine(consts.MySQL, config.Conf.App.Middleware.Mysql[consts.Default].Dsn)
	if err != nil {
		panic(err)
	}

	session := db.NewSession().MustLogSQL(true)
	testContext := context.WithValue(context.TODO(), consts.SessionKey{}, session)

	// 加载自定义配置
	_ = mainconfig.InitConfig()

	return testContext
}

func SetConfigFilePath() {
	_ = os.Setenv("YS_TOP", GetCurrentPath())
	ConfigFilePath = filepath.Join(GetCurrentPath(), mainconfig.PSPName, config.ConfigDir, "local.yml")
}

func GetCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
