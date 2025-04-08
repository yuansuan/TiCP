package util

import "github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"

// UserChecker ...
//
//go:generate mockgen -destination mock_userchecker.go -package util github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util UserChecker
type UserChecker interface {
	IsYsProductUser(userID snowflake.ID) (bool, error)
}
