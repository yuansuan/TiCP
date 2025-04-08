package client

import (
	"sync"

	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
)

// Client 作业服务依赖的rpc服务客户端
type Client struct {
	User      user.UsersClient          `grpc_client_inject:"user"`
	SysConfig sysconfig.SysConfigClient `grpc_client_inject:"sysconfig"`
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
