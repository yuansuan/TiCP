package collector

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

/*
FlexLMRE 从 lmutil lmstat -a 输出提取数据的正则表达式（Don't use underscores. 好。）

Feature                         Version     #licenses    Expires      Vendor
_______                         _________   _________    __________   ______
acdi_adprepost                  9999.9999    4           30-may-2022  ansyslmd
acpreppost                      9999.9999    4           30-may-2022  ansyslmd
advanced_meshing                9999.9999    8           30-may-2022  ansyslmd
ans_act                         9999.9999    4           30-may-2022  ansyslmd
ansys                           9999.9999    4           30-may-2022  ansyslmd
aqwa_pre                        9999.9999    4           30-may-2022  ansyslmd
aqwa_solve                      9999.9999    4           30-may-2022  ansyslmd
*/
// FlexLMRE 获取license存量和用量的正则表达式
const FlexLMRE string = `Users of (\w+):.*?Total.*?(\d+).*?issued.*?Total.*?(\d+).*?in use`

// FlexLMCollector 采集使用 FlexLM 的软件
type FlexLMCollector struct {
	Collector
}

// NewFlexLMCollector 创建 FlexLMCollector 实例
func NewFlexLMCollector(license *LicenseCollectInfo) (c *FlexLMCollector) {
	c = &FlexLMCollector{}
	c.License = license
	c.Components = make(map[string]Component)
	return
}

func (c *FlexLMCollector) parse(raw string) {
	// 延迟绑定Regexp
	c.Regexp, _ = regexp.Compile(FlexLMRE)
	parsed := c.Regexp.FindAllStringSubmatch(raw, -1)
	for _, item := range parsed {
		// 正则表达式保证了一定是数字
		total, _ := strconv.Atoi(item[2])
		used, _ := strconv.Atoi(item[3])
		c.Components[item[1]] = Component{
			Total: float64(total),
			Used:  float64(used),
		}
	}
}

// Collect 采集
func (c *FlexLMCollector) Collect() error {
	command := fmt.Sprintf("%v lmstat -a -c %v", c.License.LmstatPath, c.License.LicensePath)
	const maxRetries = 3 // 定义最大重试次数
	for i := 0; i < maxRetries; i++ {
		raw, err := ExecHpcCommand(c.License.HpcEndpoint, command)
		if err == nil { // 成功执行并解析，返回无错误
			c.parse(raw)
			return nil
		}
		// 如果不是最后一次尝试，则继续下一次尝试
		if i < maxRetries-1 {
			continue
		}
		logging.Default().Warnf("ExecHpcCommand, Error: %s, Command: %s, HpcEndpoint: %s",
			err.Error(), command, c.License.HpcEndpoint)
		return CollectorRuntimeErr
	}
	return nil // 这行实际上不会被执行
}

// GetComponents ...
func (c *FlexLMCollector) GetComponents() map[string]Component {
	return c.Components
}
