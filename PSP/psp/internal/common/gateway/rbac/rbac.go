package rbac

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/gateway/service"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ApiAuth 接口权限校验
func ApiAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		url := ctx.Request.URL.Path
		method := ctx.Request.Method
		userId := ginutil.GetUserID(ctx)

		pass, _ := service.RouteSrv.PermService.CheckApiPermission(ctx, userId, url, method)
		if !pass {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//fmt.Println("url-path: ", ctx.Request.URL.Path)
		//fmt.Println("token: ", token)
		//if ctx.Request.URL.Path == "/api/job/detail" {
		//	ctx.Header("Authenticate", token)
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
	}
}
