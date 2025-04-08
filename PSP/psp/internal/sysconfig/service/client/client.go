package client

import (
	"sync"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
)

type Client struct {
	Perm rbac.PermissionManagerClient `grpc_client_inject:"rbac"`
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
