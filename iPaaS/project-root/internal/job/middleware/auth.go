package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
)

// AdminAccessCheck 检查用户是否有管理员权限
func AdminAccessCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader(util.UserIdKeyInHeader)
		if userID == "" {
			common.ErrorResp(c, http.StatusNotFound, api.InvalidUserID, fmt.Sprintf("%s is empty in HTTP Header", util.UserIdKeyInHeader))
			return
		}

		adminYsIDs := config.GetConfig().SelfYsID

		if adminYsIDs == userID {
			c.Next()
		} else {
			common.ErrorResp(c, http.StatusForbidden, api.AccessDeniedErrorCode, "You don't have permission to access this resource")
		}
	}
}
