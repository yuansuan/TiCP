package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
)

const (
	// UserInfoKey in context.Context
	UserInfoKey = "x-ys-user-id"
	// RequestIDKey in context.Context
	RequestIDKey = "x-ys-request-id"
)

// Response ...
type Response struct {
	Data      interface{} `json:"Data"`
	ErrorCode string      `json:"ErrorCode"`
	ErrorMsg  string      `json:"ErrorMsg"`
	RequestID string      `json:"RequestID"`
}

type Error struct {
	ErrorCode string `json:"ErrorCode"`
	ErrorMsg  string `json:"ErrorMsg"`
}

func WrapError(errorCode string, errorMsg string) *Error {
	return &Error{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
	}
}

// SuccessResp ...
func SuccessResp(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{Data: data, RequestID: GetRequestID(c)})
}

func GetRequestID(c *gin.Context) string {
	return trace.GetRequestId(c)
}

// ErrorResp ...
func ErrorResp(c *gin.Context, code int, errorCode string, errorMsg string) {
	resp := Response{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
		RequestID: GetRequestID(c),
	}
	c.AbortWithStatusJSON(code, resp)
}

// InvalidParams ...
func InvalidParams(c *gin.Context, msg string) {
	resp := Response{
		ErrorCode: api.InvalidArgumentErrorCode,
		ErrorMsg:  msg,
		RequestID: GetRequestID(c),
	}
	logging.GetLogger(c).Warnf("invalid params %s", msg)
	c.AbortWithStatusJSON(http.StatusBadRequest, resp)
}

// InternalServerError ...
func InternalServerError(c *gin.Context, msg string) {
	resp := Response{
		ErrorCode: api.InternalServerErrorCode,
		ErrorMsg:  msg,
		RequestID: GetRequestID(c),
	}
	logging.GetLogger(c).Errorf("internal server error %s", msg)
	c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
}
