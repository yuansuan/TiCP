package scheduler

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

type PostPaidBillSliceMessage struct {
	AccountId      snowflake.ID
	PayByAccountId snowflake.ID
	SessionId      snowflake.ID
	ResourceId     snowflake.ID
	ResourceType   resourceType
	IsInitial      bool
	IsFirst        bool
	IsFinish       bool
	StartTime      time.Time
	EndTime        time.Time
	// 保证一个Session只会传递一个消息标记Session付费结束（一个Session包含software/hardware两个资源，无法使用消息体中的IsFinished标识Session付费结束）
	MarkSessionIsFinishedPaid bool
}
