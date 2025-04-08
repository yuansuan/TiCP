package errcode

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// Common error codes from 10001 to 11000
const (
	ErrInternalServer codes.Code = 10001
	ErrInvalidParam   codes.Code = 10002
	ErrInvalidAction  codes.Code = 10003
	ErrOpenAPIAction  codes.Code = 10004
	ErrCloudDisable   codes.Code = 10005
)

// Common error message
const (
	MsgInternalServer = "服务器内部错误"
	MsgInvalidParam   = "请求参数错误"
	MsgCloudDisable   = "云端模式未开启"
)

var codeMsgOnce sync.Once
var codeMsg = map[codes.Code]string{
	ErrInternalServer: "服务器内部错误",
	ErrInvalidParam:   "请求参数错误",
	ErrInvalidAction:  "操作不支持",
	ErrOpenAPIAction:  "OpenAPI 操作失败",
	ErrCloudDisable:   "混合云模式未开启",
}

func ResolveErrCodeMessage(ctx *gin.Context, err error, defCode codes.Code) {
	logging.Default().Errorf("errMsg：%+v", err.Error())
	code := status.Code(err)
	msg := GetCodeMsg()[code]
	if msg == "" {
		ginutil.Error(ctx, defCode, GetCodeMsg()[defCode])
	} else {
		ginutil.Error(ctx, code, msg)
	}
}

func GetErrCodeMessage(err error, defCode codes.Code) string {
	code := status.Code(err)
	msg := GetCodeMsg()[code]
	if msg == "" {
		return GetCodeMsg()[defCode]
	} else {
		return msg
	}
}

// GetCodeMsg ...
func GetCodeMsg() map[codes.Code]string {
	codeMsgOnce.Do(func() {
		for k, v := range AppCodeMsg {
			codeMsg[k] = v
		}
		for k, v := range ProjectCodeMsg {
			codeMsg[k] = v
		}
		for k, v := range VisualCodeMsg {
			codeMsg[k] = v
		}
		for k, v := range RBACCodeMsg {
			codeMsg[k] = v
		}
		for k, v := range UserCodeMsg {
			codeMsg[k] = v
		}
		for k, v := range StorageCodeMsg {
			codeMsg[k] = v
		}
		for k, v := range NoticeCodeMsg {
			codeMsg[k] = v
		}
		for k, v := range SysConfigCodeMsg {
			codeMsg[k] = v
		}
		for k, v := range JobCodeMsg {
			codeMsg[k] = v
		}
		for k, v := range ApproveCodeMsg {
			codeMsg[k] = v
		}
	})
	return codeMsg
}
