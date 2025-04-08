package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	boothttp "github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/admin"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/middleware"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/scheduler"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/state"
)

func Init(drv *boothttp.Driver) {
	s, err := state.New()
	if err != nil {
		fmt.Println(err)
		logging.Default().Error(fmt.Sprintf("init state failed, administrator check! %v", err))
		os.Exit(-1)
	}

	drv.Use(func(c *gin.Context) {
		c.Set("state", s)
	})

	for _, mid := range middleware.Middlewares() {
		drv.Use(mid)
	}

	apiGrp := drv.Group("/api")
	initApiEndpoint(apiGrp)

	adminGrp := drv.Group("/admin")
	initAdminEndpoint(adminGrp)

	sched, err := scheduler.NewScheduler(s)
	if err != nil {
		logging.Default().Fatalf("new scheduler failed, %v", err)
	}

	go sched.Start()
}

func initApiEndpoint(apiGrp *gin.RouterGroup) {
	apiGrp.POST("/sessions", PostSessions)
	apiGrp.GET("/sessions/:SessionId", GetSession)
	apiGrp.GET("/sessions", ListSession)
	apiGrp.POST("/sessions/:SessionId/close", CloseSession)
	apiGrp.GET("/sessions/:SessionId/ready", SessionReady)
	apiGrp.DELETE("/sessions/:SessionId", DeleteSession)
	apiGrp.POST("/sessions/:SessionId/start", StartSession)
	apiGrp.POST("/sessions/:SessionId/stop", StopSession)
	apiGrp.POST("/sessions/:SessionId/restart", RestartSession)
	apiGrp.POST("/sessions/:SessionId/restore", RestoreSession)
	apiGrp.POST("/sessions/:SessionId/execScript", SessionExecScript)
	apiGrp.POST("/sessions/:SessionId/mount", SessionMount)
	apiGrp.POST("/sessions/:SessionId/umount", SessionUmount)

	apiGrp.GET("/sessions/:SessionId/remoteapps/:RemoteAppName", GetRemoteApp)

	apiGrp.GET("/hardwares/:HardwareId", GetHardWare)
	apiGrp.GET("/hardwares", ListHardWare)

	apiGrp.GET("/softwares/:SoftwareId", GetSoftWare)
	apiGrp.GET("/softwares", ListSoftWare)
}

func initAdminEndpoint(adminGrp *gin.RouterGroup) {
	adminGrp.GET("/sessions", admin.ListSession)
	adminGrp.POST("/sessions/:SessionId/close", admin.CloseSession)
	adminGrp.POST("/sessions/:SessionId/start", admin.StartSession)
	adminGrp.POST("/sessions/:SessionId/stop", admin.StopSession)
	adminGrp.POST("/sessions/:SessionId/restart", admin.RestartSession)
	adminGrp.POST("/sessions/:SessionId/restore", admin.RestoreSession)
	adminGrp.POST("/sessions/:SessionId/execScript", admin.ExecScript)
	adminGrp.POST("/sessions/:SessionId/mount", admin.SessionMount)
	adminGrp.POST("/sessions/:SessionId/umount", admin.SessionUmount)

	adminGrp.POST("/hardwares", admin.PostHardwares)
	adminGrp.PUT("/hardwares/:HardwareId", admin.PutHardware)
	adminGrp.PATCH("/hardwares/:HardwareId", admin.PatchHardware)
	adminGrp.DELETE("/hardwares/:HardwareId", admin.DeleteHardware)
	adminGrp.GET("/hardwares/:HardwareId", admin.GetHardware)
	adminGrp.GET("/hardwares", admin.ListHardware)
	adminGrp.POST("/hardwares/users", admin.PostHardwaresUsers)
	adminGrp.DELETE("/hardwares/users", admin.DeleteHardwaresUsers)

	adminGrp.POST("/softwares", admin.PostSoftwares)
	adminGrp.PUT("/softwares/:SoftwareId", admin.PutSoftware)
	adminGrp.PATCH("/softwares/:SoftwareId", admin.PatchSoftware)
	adminGrp.DELETE("/softwares/:SoftwareId", admin.DeleteSoftware)
	adminGrp.GET("/softwares/:SoftwareId", admin.GetSoftware)
	adminGrp.GET("/softwares", admin.ListSoftware)
	adminGrp.POST("/softwares/users", admin.PostSoftwaresUsers)
	adminGrp.DELETE("/softwares/users", admin.DeleteSoftwaresUsers)

	adminGrp.POST("/remoteapps", admin.PostRemoteApps)
	adminGrp.PUT("/remoteapps/:RemoteAppId", admin.PutRemoteApp)
	adminGrp.PATCH("/remoteapps/:RemoteAppId", admin.PatchRemoteApp)
	adminGrp.DELETE("/remoteapps/:RemoteAppId", admin.DeleteRemoteApp)
}
