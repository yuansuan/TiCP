package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

type ErrorResp struct {
	code    string
	message string
}

func WrapErrorResp(code, message string) ErrorResp {
	return ErrorResp{code, message}
}

func BadRequestIfError(c *gin.Context, err error, errResp ErrorResp) error {
	if err != nil {
		errResponse(c, errResp, http.StatusBadRequest)
	}

	return err
}

func InternalErrorIfError(c *gin.Context, err error, errResp ErrorResp) error {
	if err != nil {
		errResponse(c, errResp, http.StatusInternalServerError)
	}

	return err
}

func UnauthorizedIfError(c *gin.Context, err error, errResp ErrorResp) error {
	if err != nil {
		errResponse(c, errResp, http.StatusUnauthorized)
	}

	return err
}

func ForbiddenIfError(c *gin.Context, err error, errResp ErrorResp) error {
	if err != nil {
		errResponse(c, errResp, http.StatusForbidden)
	}

	return err
}

func NotFoundIfError(c *gin.Context, err error, errResp ErrorResp) error {
	if err != nil {
		errResponse(c, errResp, http.StatusNotFound)
	}

	return err
}

func ConflictIfError(c *gin.Context, err error, errResp ErrorResp) error {
	if err != nil {
		errResponse(c, errResp, http.StatusConflict)
	}

	return err
}

func errResponse(c *gin.Context, err ErrorResp, errorCode int) {
	c.JSON(errorCode, v20230530.Response{
		ErrorCode: err.code,
		ErrorMsg:  err.message,
		RequestID: trace.GetRequestId(c),
	})
}

func RenderJson(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, common.Response{
		Data:      data,
		RequestID: trace.GetRequestId(c),
	})
}
