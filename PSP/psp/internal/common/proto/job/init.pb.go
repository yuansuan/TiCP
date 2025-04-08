package job

import boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

var _ boot.ServerType

func init() {
	boot.RegisterClient("job", NewJobClient)
}
