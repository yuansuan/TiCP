package helloworld

import grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

var _ grpc_boot.ServerType

func init() {
	grpc_boot.RegisterClient("helloworld", NewGreeterClient)
}
