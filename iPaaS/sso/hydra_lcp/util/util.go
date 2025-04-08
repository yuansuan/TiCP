package util

import (
	"os"
	"strconv"
	"sync"

	"github.com/ory/hydra/sdk/go/hydra/client"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/config"
)

// HydraConfig hydra config
type HydraConfig struct {
	// 对外服务登录token过期时间，单位秒
	TokenExpireTime int64
	HydraClient     *client.OryHydra

	// 内部服务登录token过期时间，单位秒
	TokenExpireTimePortal int64
	HydraClientPortal     *client.OryHydra

	// 每个客户端IP每小时最大可发送数
	PerIPHourSendPhoneCodeMax int64
	// 验证失败最大次数
	VerifyFailMax int64
	// 验证错误次数超过最大值，手机号冻结登录时间; 单位秒
	VerifyFailOverMaxFreezeLoginTime int64
}

var (
	hydraConfig     = new(HydraConfig)
	hydraConfigOnce sync.Once
)

// GetHydraConfig get hydra config
func GetHydraConfig() *HydraConfig {
	conf := config.Custom
	hydraConfigOnce.Do(func() {
		hydraConfig.HydraClient = client.NewHTTPClientWithConfig(nil,
			&client.TransportConfig{
				Schemes:  []string{conf.HydraAdminConf.Scheme},
				Host:     conf.HydraAdminConf.Host,
				BasePath: conf.HydraAdminConf.Path,
			})
		hydraConfig.HydraClientPortal = client.NewHTTPClientWithConfig(nil,
			&client.TransportConfig{
				Schemes:  []string{conf.HydraPortalAdminConf.Scheme},
				Host:     conf.HydraPortalAdminConf.Host,
				BasePath: conf.HydraPortalAdminConf.Path,
			})

		configMap := os.Getenv("TOKEN_EXPIRE_TIME")
		if configMap != "" {
			hydraConfig.TokenExpireTime, _ = strconv.ParseInt(configMap, 10, 64)
		}

		configMap = os.Getenv("TOKEN_EXPIRE_TIME_PORTAL")
		if configMap != "" {
			hydraConfig.TokenExpireTimePortal, _ = strconv.ParseInt(configMap, 10, 64)
		}

		configMap = os.Getenv("PER_IP_HOUR_SEND_SMS_CODE_MAX")
		if configMap != "" {
			hydraConfig.PerIPHourSendPhoneCodeMax, _ = strconv.ParseInt(configMap, 10, 64)
		} else {
			hydraConfig.PerIPHourSendPhoneCodeMax = 100
		}

		configMap = os.Getenv("VERIFY_FAIL_MAX")
		if configMap != "" {
			hydraConfig.VerifyFailMax, _ = strconv.ParseInt(configMap, 10, 64)
		} else {
			hydraConfig.VerifyFailMax = 5
		}

		configMap = os.Getenv("VERIFY_FAIL_OVER_MAX_FREEZE_LOGIN_TIME")
		if configMap != "" {
			hydraConfig.VerifyFailOverMaxFreezeLoginTime, _ = strconv.ParseInt(configMap, 10, 64)
		} else {
			hydraConfig.VerifyFailOverMaxFreezeLoginTime = 3600
		}
	})

	return hydraConfig
}
