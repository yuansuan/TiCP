package lockutil

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"time"
)

func TryLock(taskKey string) (bool, error) {
	redisClient := boot.Middleware.DefaultRedis()
	boolCmd := redisClient.SetNX(taskKey, true, 1*time.Second)
	successFlag, err := boolCmd.Result()
	if err != nil {
		return false, err
	}

	return successFlag, nil
}

func UnLock(taskKey string) {
	redisClient := boot.Middleware.DefaultRedis()
	redisClient.Del(taskKey)
}
