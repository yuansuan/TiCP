package client

import (
	"sync"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/job"
)

// GRPC 服务依赖的grpc服务客户端
type GRPC struct {
	Job job.JobClient `grpc_client_inject:"job"`
}

var (
	grpc *GRPC
	once sync.Once
)

// GetInstance 获取客户端实例
func GetInstance() *GRPC {
	once.Do(func() {
		grpc = &GRPC{}
		grpc_boot.InjectAllClient(grpc)
	})

	return grpc
}
