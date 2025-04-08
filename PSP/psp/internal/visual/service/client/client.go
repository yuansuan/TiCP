package client

import (
	"sync"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
)

type Client struct {
	Project project.ProjectClient `grpc_client_inject:"project"`
	Notice  notice.NoticeClient   `grpc_client_inject:"notice"`

	RBAC struct {
		Permission rbac.PermissionManagerClient `grpc_client_inject:"rbac"`
		Role       rbac.RoleManagerClient       `grpc_client_inject:"rbac"`
	}
}

var (
	client *Client
	once   sync.Once
)

func GetInstance() *Client {
	once.Do(func() {
		client = &Client{}
		grpc_boot.InjectAllClient(client)
	})
	return client
}
