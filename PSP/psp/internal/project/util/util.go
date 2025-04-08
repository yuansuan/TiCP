package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

func ValidOwnerID(ID string) (bool, error) {
	resultID, err := snowflake.ParseString(ID)
	if err != nil {
		return false, err
	}

	if resultID == 0 {
		return false, fmt.Errorf("snowflake parse id:[%s] err", ID)
	}

	return true, nil
}

func CheckErrorIfNoPermission(ctx *gin.Context, err error, errCode codes.Code) {
	if status.Code(err) == errcode.ErrProjectAccessPermission {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectAccessPermission)
		return
	} else if status.Code(err) == errcode.ErrProjectNotFound {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectNotFound)
		return
	}

	errcode.ResolveErrCodeMessage(ctx, err, errCode)
}
