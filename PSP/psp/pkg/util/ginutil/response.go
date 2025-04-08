package ginutil

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"google.golang.org/grpc/codes"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
)

const (
	UserID   = "user_id"
	UserName = "user_name"
)

// Error 执行失败返回数据
func Error(ctx *gin.Context, code codes.Code, msg string) {
	ctx.JSON(nethttp.StatusOK, http.Resp{
		Success: false,
		Code:    code,
		Message: msg,
	})
}

// Success 执行成功返回数据
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(nethttp.StatusOK, http.Resp{
		Success: true,
		Data:    data,
	})
}

func GetUserID(ctx *gin.Context) int64 {
	return ctx.GetInt64(UserID)
}

func GetUserName(ctx *gin.Context) string {
	return ctx.GetString(UserName)
}

func GetTraceID(ctx *gin.Context) string {
	return ctx.GetString(common.TraceId)
}

func SetTraceID(ctx *gin.Context, traceId string) {
	ctx.Set(common.TraceId, traceId)
}

func SetUser(ctx *gin.Context, userID int64, userName string) {
	ctx.Set(UserID, userID)
	ctx.Set(UserName, userName)
}

func GetRequestIP(ctx *gin.Context) string {
	reqIP := ctx.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}
