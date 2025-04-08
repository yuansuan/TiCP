package dao

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// MessageDao 消息表数据访问
type MessageDao interface {
	// InsertMessage 保存消息信息
	InsertMessage(ctx context.Context, message *model.Message) error
	// ReadMessage 消息设置已读取
	ReadMessage(ctx context.Context, userID snowflake.ID, ids []snowflake.ID) error
	// ReadAllMessage 消息设置已读取
	ReadAllMessage(ctx context.Context, userID snowflake.ID) error
	// SaveAndSendMessage 保存并发送消息
	SaveAndSendMessage(ctx context.Context, message *model.Message) (snowflake.ID, error)
	// GetMessageCount 获取消息数量统计
	GetMessageCount(ctx context.Context, userID snowflake.ID, state int) (int64, error)
	// GetMessageList 分页获取消息列表
	GetMessageList(ctx context.Context, msg *dto.MessagePage) ([]*model.Message, int64, error)
}
