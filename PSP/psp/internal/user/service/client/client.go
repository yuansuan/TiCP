package client

import (
	"sync"

	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
)

// Client 用户服务依赖的rpc服务客户端
type Client struct {
	Role      rbac.RoleManagerClient          `grpc_client_inject:"rbac"`
	Perm      rbac.PermissionManagerClient    `grpc_client_inject:"rbac"`
	SysConfig sysconfig.SysConfigClient       `grpc_client_inject:"sysconfig"`
	Storage   storage.StorageClient           `grpc_client_inject:"storage"`
	Approve   approve.ApproveManagementClient `grpc_client_inject:"approve"`
}

var (
	once   sync.Once
	client *Client
)

// GetInstance 获取客户端
func GetInstance() *Client {
	once.Do(func() {
		client = &Client{}
		grpc_boot.InjectAllClient(client)
	})

	return client
}
