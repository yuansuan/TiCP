package util

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

func GetResourceId(ctx *gin.Context) (string, bool) {
	pathId := ctx.Param("id")
	id, err := snowflake.ParseString(pathId)
	if err != nil || id == 0 {
		logging.GetLogger(ctx).Warnf("invalid params invalid id")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return "", false
	}
	return pathId, true
}
