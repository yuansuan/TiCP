package yamlutil

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const (
	PreKey = "yaml"
)

// ReplaceConfig 替换制定配置文件对应key的value,不会导致注释消失或重排序
// example:
// configPath: config/prod_custom.yml
// key: user.ldap.enable
// value: true
func ReplaceConfig(configPath, key, value string) error {
	// 读取 YAML 文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	node := yaml.Node{}

	err = yaml.Unmarshal(data, &node)
	if err != nil {
		return err
	}

	childNode, err := FindChildNode(&node, key)
	if err != nil {
		return err
	}

	childNode.Value = value

	b, err := yaml.Marshal(&node)
	if err != nil {
		return err
	}

	// 将修改后的数据写入文件
	err = os.WriteFile(configPath, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func FindChildNode(node *yaml.Node, identifier string) (*yaml.Node, error) {
	if identifier == "" {
		return nil, errors.New("identifier不能为空")
	}

	identifier = fmt.Sprintf("%s.%s", PreKey, identifier)
	identifiers := strings.Split(identifier, ".")
	result := findNode(node, identifiers, true)
	if result == nil {
		return nil, errors.New("未找到匹配的节点")
	}
	return result, nil
}

// yaml3的node结构很抽象，a节点的子节点b，位置是a的下标+1，也就是说为了找到子节点，需要做遍历而不是递归
func findNode(node *yaml.Node, identifiers []string, findNextNode bool) *yaml.Node {
	if len(identifiers) == 0 {
		return nil
	}
	// 需要匹配的节点名
	identifier := identifiers[0]

	// 匹配节点名，如果匹配上了就递归下标+1的节点内容
	for _, n := range node.Content {
		if findNextNode {
			nextIdentifiers := identifiers[1:]
			// 说明后续没有需要匹配的了，匹配完成，直接返回
			if len(nextIdentifiers) == 0 {
				return n
			}
			// 说明匹配到上一个节点，需要进入子节点的遍历
			if len(n.Content) > 0 {
				return findNode(n, nextIdentifiers, false)
			}
			break
		}

		if n.Value == identifier {
			findNextNode = true
			continue
		}
	}
	return nil
}
