package common

import (
	"fmt"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
)

func GetRemoteConfigPath() string {
	return fmt.Sprintf("project_root/%s/%s_custom.yaml", boot.Config.App.Name, env.ModeName(boot.Env.Mode))
}
