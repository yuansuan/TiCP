package dataloader

import (
	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/util"
)

// InitDaemon ...
func InitDaemon() {
	// 初始化主机名映射
	util.InitHostNameMapping()

	nodeLoader, err := NewNodeLoader()
	if err != nil {
		panic(errors.Wrap(err, "failed to init daemon"))
	}

	nodeLoader.NodeLoaderStart()
}
