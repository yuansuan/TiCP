package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

// GetUserID 从HTTP Header中获取用户ID
func GetUserID(c *gin.Context) (snowflake.ID, error) {
	userIDStr := c.GetHeader(common.UserInfoKey)
	if userIDStr == "" {
		return 0, fmt.Errorf("%s is empty in HTTP Header", common.UserInfoKey)
	}

	userID, err := snowflake.ParseString(userIDStr)
	if err != nil {
		return 0, fmt.Errorf("parse userId failed, %w", err)
	}

	return userID, nil
}
