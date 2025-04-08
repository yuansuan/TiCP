package collector

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

/*
LSDynaRE 从 lstc_run -q 输出提取数据的正则表达式

PROGRAM          EXPIRATION CPUS  USED   FREE    MAX | QUEUE
---------------- ----------      ----- ------ ------ | -----
LS-DYNA          01/15/2038          0    256    256 |     0
MPPDYNA          01/15/2038          0    256    256 |     0
LS-DYNA_971      01/15/2038          0    256    256 |     0
MPPDYNA_971      01/15/2038          0    256    256 |     0

	LICENSE GROUP     0    256    256 |     0
*/
const LSDynaRE string = `(\d+\S?\d+\S?\d+)\s+\d+\s+\d+\s+\d+\s+\|\s+\d+\s+\w+\s+\w+\s+(\d+)\s+\d+\s+(\d+)`

// LSDynaCollector 采集 LS-DYNA 数据
// LS-DYNA 自行实现了一套
type LSDynaCollector struct {
	Collector
}

// NewLSDynaCollector 创建 LSDynaCollector 实例
func NewLSDynaCollector(license *LicenseCollectInfo) (c *LSDynaCollector) {
	// 测试过的正则表达式，可以忽略
	re, _ := regexp.Compile(LSDynaRE)
	c = &LSDynaCollector{}
	c.License = license
	c.Regexp = re
	c.Components = make(map[string]Component)
	return
}

func (c *LSDynaCollector) parse(raw string) {
	parsed := c.Regexp.FindStringSubmatch(raw)
	if len(parsed) < 4 {
		logging.Default().Warnf("Warning: parsed result does not contain enough elements.Input was: %s", raw)
		return
	}
	// 正则表达式保证了一定是数字
	total, errTotal := strconv.Atoi(parsed[3])
	if errTotal != nil {
		logging.Default().Error("Error converting total '%s' to int: %v", parsed[3], errTotal)
		return
	}
	used, errUsed := strconv.Atoi(parsed[2])
	if errUsed != nil {
		logging.Default().Error("Error converting used '%s' to int: %v", parsed[2], errUsed)
		return
	}
	c.Components["all"] = Component{
		Total: float64(total),
		Used:  float64(used),
	}

}

// Collect 采集
func (c *LSDynaCollector) Collect() error {
	command := fmt.Sprintf("%v -r -s %v", c.License.LmstatPath, c.License.LicensePath)
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
func (c *LSDynaCollector) GetComponents() map[string]Component {
	return c.Components
}
