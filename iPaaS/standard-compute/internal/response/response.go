package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

func BadRequestIfError(c *gin.Context, err error, errCode string) error {
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err, errCode)
	}

	return err
}

func NotfoundIfError(c *gin.Context, err error, errCode string) error {
	if err != nil {
		errorResponse(c, http.StatusNotFound, err, errCode)
	}

	return err
}

func InternalErrorIfError(c *gin.Context, err error, errCode string) error {
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err, errCode)
	}

	return err
}

func ServiceUnavailableIfError(c *gin.Context, err error, errCode string) error {
	if err != nil {
		errorResponse(c, http.StatusServiceUnavailable, err, errCode) // 响应503
	}

	return err
}

func ForbiddenIfError(c *gin.Context, err error, errCode string) error {
	if err != nil {
		errorResponse(c, http.StatusForbidden, err, errCode)
	}

	return err
}

type Response struct {
	v20230530.Response

	Data interface{} `json:"Data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, Response{
		Response: v20230530.Response{
			RequestID: trace.GetRequestId(c),
		},
		Data: data,
	})
}

func errorResponse(c *gin.Context, statusCode int, err error, errCode string) {
	c.AbortWithStatusJSON(statusCode, v20230530.Response{
		ErrorCode: errCode,
		ErrorMsg:  err.Error(),
		RequestID: trace.GetRequestId(c),
	})
}
