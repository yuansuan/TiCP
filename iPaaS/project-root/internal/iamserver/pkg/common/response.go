package common

import (
	"net/http"

	"github.com/google/uuid"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
)

type UserInfo struct {
	UserID snowflake.ID
	Tag    string
	IsTmp  bool
}

func (u *UserInfo) IsAdmin() bool {
	return u.Tag == "IamAdmin"
}

const (
	// UserInfoKey in context.Context
	UserInfoKey = "x-ys-user-id"
	// RequestIDKey in context.Context
	RequestIDKey    = "x-ys-request-id"
	UserAccessKeyId = "x-ys-user-access-key-id"
)

const (
	// Success ...
	Success = "Success"
)

// GetUserInfo from context
func GetUserInfo(c *gin.Context) *UserInfo {
	data, exists := c.Get(UserInfoKey)
	if !exists {
		// never reach here
		panic(1)
	}
	return data.(*UserInfo)
}

func SetUserInfo(c *gin.Context, u *UserInfo) {
	c.Set(UserInfoKey, u)
}

func SetRequestID(c *gin.Context) string {
	reqId := uuid.New().String()
	c.Set(RequestIDKey, reqId)
	return reqId
}

// SuccessResp ...
func SuccessResp(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &iam_api.BasicResponse{Data: data, RequestId: GetRequestID(c)})
}

func GetRequestID(c *gin.Context) string {
	value, _ := c.Get(RequestIDKey)
	if value == nil {
		value = "10086"
	}
	return value.(string)
}

// ErrorResp ...
func ErrorResp(c *gin.Context, code int, errorCode string, errorMsg string) {
	resp := iam_api.BasicResponse{
		ErrorCode:    errorCode,
		ErrorMessage: errorMsg,
		RequestId:    GetRequestID(c),
	}
	c.JSON(code, resp)
	return
}

func ErrorRespWithAbort(c *gin.Context, code int, errorCode string, errorMsg string) {
	resp := iam_api.BasicResponse{
		ErrorCode:    errorCode,
		ErrorMessage: errorMsg,
		RequestId:    GetRequestID(c),
	}
	c.AbortWithStatusJSON(code, resp)
	return
}

// InvalidParams ...
func InvalidParams(c *gin.Context, msg string) {
	// resp := Response{
	// 	ErrorCode: InvalidArgumentErrorCode,
	// 	ErrorMsg:  msg,
	// 	RequestID: GetRequestID(c),
	// }
	// logging.GetLogger(c).Warnf("invalid params %s", msg)
	// c.AbortWithStatusJSON(http.StatusBadRequest, resp)
	return
}

// InternalServerError ...
func InternalServerError(c *gin.Context, msg string) {
	resp := iam_api.BasicResponse{
		ErrorCode:    "InternalServerError",
		ErrorMessage: msg,
		RequestId:    GetRequestID(c),
	}
	//logging.GetLogger(c).Warnf("internal server error %s", msg)
	c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
	return
}
