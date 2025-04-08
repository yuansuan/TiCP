package collector

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// AltairRE 从 lmutil lmstat -a 输出提取数据的正则表达式（Don't use underscores. 好。）
const AltairRE string = `Feature: (\w+)(?:.|\n)+?(\d+) of (\d+) license\(s\) used`

// AltairCollector 采集使用 Altair 的软件
type AltairCollector struct {
	Collector
}

// NewAltairCollector 创建 AltairCollector 实例
func NewAltairCollector(license *LicenseCollectInfo) (c *AltairCollector) {
	// 测试过的正则表达式，可以忽略
	re, _ := regexp.Compile(AltairRE)
	c = &AltairCollector{}
	c.License = license
	c.Regexp = re
	c.Components = make(map[string]Component)
	return
}

func (c *AltairCollector) parse(raw string) {
	parsed := c.Regexp.FindAllStringSubmatch(raw, -1)
	for _, item := range parsed {
		// 正则表达式保证了一定是数字
		total, _ := strconv.Atoi(item[3])
		used, _ := strconv.Atoi(item[2])
		c.Components[item[1]] = Component{
			Total: float64(total),
			Used:  float64(used),
		}
	}
}

// Collect 采集
func (c *AltairCollector) Collect() error {
	// parse `port@host` and `host`
	path := strings.Split(c.License.LicensePath, "@")
	host, port := path[0], ""
	if len(path) > 1 {
		port = host
		host = path[1]
	}
	command := fmt.Sprintf("%v -licstat -host %v -port %v", c.License.LmstatPath, host, port)
	raw, err := ExecHpcCommand(c.License.HpcEndpoint, command)
	if err != nil {
		logging.Default().Warnf("ExecHpcCommandFail, Error: %s, Command: %s, HpcEndpoint: %s",
			err.Error(), command, c.License.HpcEndpoint)
		return CollectorRuntimeErr
	}
	c.parse(raw)
	return nil
}

// GetComponents ...
func (c *AltairCollector) GetComponents() map[string]Component {
	return c.Components
}
