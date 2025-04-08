package collector

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	openapigo "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/config"
)

var (
	CollectorRuntimeErr = errors.New("collector runtime error")
)

// Collector 采集器基类，Collect的实现继承此类
type Collector struct {
	ICollector
	License        *LicenseCollectInfo  //某一款app具体版本：License server信息+超算服务地址
	*regexp.Regexp                      //每款收集器的正则
	Components     map[string]Component //key:模块，val:（剩余数量、总数、过期时间）
}

// ICollector 采集器接口
type ICollector interface {
	Collect() error
	GetComponents() map[string]Component
}

// Component 组件数据
type Component struct {
	Total, Used float64
}

// LicenseCollectInfo model
type LicenseCollectInfo struct {
	LicensePath   string //ip@port
	LmstatPath    string //pwd 路径
	HpcEndpoint   string
	CollectorType string
}

// NewCollector 靠appType区分new的类型
func NewCollector(license *LicenseCollectInfo) (collector ICollector, err error) {
	switch license.CollectorType {
	case "flex":
		collector = NewFlexLMCollector(license)
	case "lsdyna":
		collector = NewLSDynaCollector(license)
	case "altair":
		collector = NewAltairCollector(license)
	case "dsli":
		collector = NewDSLiCollector(license)
	default:
		err = errors.New("不支持的收集器类型: " + license.CollectorType)
	}
	return
}

func ExecHpcCommand(hpcEndpoint string, command string) (string, error) {
	logging.Default().Info("StartExecHpcCommand, HpcEndPoint: %s, Command: %s", command)
	akId, akSecret := config.GetCustom().AccessKeyId, config.GetCustom().AccessKeySecret
	c, err := openapigo.NewClient(credential.NewCredential(akId, akSecret), openapigo.WithBaseURL(hpcEndpoint))
	if err != nil {
		return "", err
	}
	resp, err := c.HPC.Command.System.Execute(
		c.HPC.Command.System.Execute.Command(command),
		c.HPC.Command.System.Execute.Timeout(30),
	)
	if err != nil {
		return "", err
	}
	if resp.ErrorCode != "" {
		return "", errors.New(fmt.Sprintf("ErrorCode: %s, ErrorMsg: %s, RequestId: %s",
			resp.ErrorCode, resp.ErrorMsg, resp.RequestID))
	}
	if resp.Data.IsTimeout {
		return "", errors.New(fmt.Sprintf("ExecCommandTimeout"))
	}
	if resp.Data.ExitCode != 0 {
		return "", errors.New(fmt.Sprintf("ExitCode: %d, StdOut: %s, StdErr: %s",
			resp.Data.ExitCode, resp.Data.Stdout, resp.Data.Stderr))
	}
	return resp.Data.Stdout, nil
}
