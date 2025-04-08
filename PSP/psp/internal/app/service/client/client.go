package client

import (
	"sync"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/license"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/monitor"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
)

type Client struct {
	RBAC struct {
		Permission rbac.PermissionManagerClient `grpc_client_inject:"rbac"`
		Role       rbac.RoleManagerClient       `grpc_client_inject:"rbac"`
	}

	Monitor monitor.MonitorClient `grpc_client_inject:"monitor"`

	License license.LicenseClient `grpc_client_inject:"license"`
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
