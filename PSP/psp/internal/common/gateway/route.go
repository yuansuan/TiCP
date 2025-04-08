package gateway

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/gateway/auth"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/gateway/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/gateway/service"
)

// InitGateway ...
func InitGateway(drv *http.Driver) {

	routeSrv, err := service.NewRouteService()
	if err != nil {
		panic(err)
	}
	service.RouteSrv = routeSrv

	drv.Use(auth.BasicAuth())
	drv.Use(rbac.ApiAuth())

}
