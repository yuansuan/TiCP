package handler_rpc

import (
	"context"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
	"github.com/yuansuan/ticp/common/project-root-api/proto/license"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/handler_rpc/impl"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/scheduler"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/leader"
)

const (
	LicenseMonitorRunDuration = 60 * time.Second
)

// InitGRPCServer ...
func InitGRPCServer(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer() //获取gRPC服务器
	util.ChkErr(err)
	//dao.SyncTable()
	licenseServer := impl.NewLicenseServer(boot.MW.DefaultORMEngine())
	//将自己的服务实现注册到gRPC服务器中（GRPC服务器就知道在处理对应的GRPC
	// -请求时应该调用 licenseServer 对象中实现的方法了）
	license.RegisterLicenseManagerServiceServer(s.Driver(), licenseServer)
	//监控license使用情况
	ctx := context.TODO()
	// 创建定时任务获取license的信息
	leader.Runner(ctx, "license_info_scheduler", scheduler.NewScheduler().Run,
		leader.SetInterval(LicenseMonitorRunDuration))
}
