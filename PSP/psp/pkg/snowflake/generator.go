package snowflake

import (
	"sync"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

const (
	DefaultNodeNum = 1
)

var (
	node *Node
	once sync.Once
)

// GetInstance 获取Snowflake实例
func GetInstance() (*Node, error) {
	var err error
	once.Do(func() {
		node, err = NewNode(DefaultNodeNum)
		if err != nil {
			logging.Default().Errorf("create snowflake instance err: %v", err)
		}
	})

	return node, err
}
