package client

import (
	"sync"

	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
)

// Client 作业服务依赖的rpc服务客户端
type Client struct {
	// UserManagement Invoke the user management gRPC APIs
	Users    user.UsersClient                 `grpc_client_inject:"user"`
	Notice   notice.NoticeClient              `grpc_client_inject:"notice"`
	AuditLog approve.AuditLogManagementClient `grpc_client_inject:"approve"`
	Project  project.ProjectClient            `grpc_client_inject:"project"`
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
