package collector

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// DSLi 从 -admin -run 命令输出提取数据的正则表达式
const DSLi string = `(\w+)\s+maxReleaseNumber:\s*(\d+).*?count:\s*(\d+).*?inuse:\s*(\d+)`

// DSLiCollector 采集使用 DSLi 的软件
type DSLiCollector struct {
	Collector
}

func NewDSLiCollector(license *LicenseCollectInfo) (c *DSLiCollector) {
	re, _ := regexp.Compile(DSLi)
	c = &DSLiCollector{}
	c.License = license
	c.Regexp = re
	c.Components = make(map[string]Component)
	return
}

func (c *DSLiCollector) parse(raw string) {
	parsed := c.Regexp.FindAllStringSubmatch(raw, -1)
	for _, item := range parsed {
		if len(item) < 5 {
			continue
		}
		name := item[1]
		total, _ := strconv.Atoi(item[3])
		used, _ := strconv.Atoi(item[4])
		c.Components[name] = Component{
			Total: float64(total),
			Used:  float64(used),
		}
	}
}

// Collect 采集
func (c *DSLiCollector) Collect() error {
	path := strings.Split(c.License.LicensePath, " ")
	host, port := path[0], ""
	if len(path) > 1 {
		port = path[1]
	}
	command := fmt.Sprintf("%v -admin -run \"connect %v %v;getLicenseUsage\"", c.License.LmstatPath, host, port)
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
func (c *DSLiCollector) GetComponents() map[string]Component { return c.Components }
