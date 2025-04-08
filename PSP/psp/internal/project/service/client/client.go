package client

import (
	"sync"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
)

// GRPC 作业服务依赖的grpc服务客户端
type GRPC struct {
	User    user.UsersClient             `grpc_client_inject:"user"`
	Notice  notice.NoticeClient          `grpc_client_inject:"notice"`
	Storage storage.StorageClient        `grpc_client_inject:"storage"`
	Rbac    rbac.PermissionManagerClient `grpc_client_inject:"rbac"`
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
