package client

import (
	"sync"

	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
)

type Client struct {
	Rbac      rbac.RoleManagerClient       `grpc_client_inject:"rbac"`
	Perm      rbac.PermissionManagerClient `grpc_client_inject:"rbac"`
	User      user.UsersClient             `grpc_client_inject:"user"`
	Notice    notice.NoticeClient          `grpc_client_inject:"notice"`
	SysConfig sysconfig.SysConfigClient    `grpc_client_inject:"sysconfig"`
}

var (
	once   sync.Once
	client *Client
)

func GetInstance() *Client {
	once.Do(func() {
		client = &Client{}
		grpc_boot.InjectAllClient(client)
	})

	return client
}
