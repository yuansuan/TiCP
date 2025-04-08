package service

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dto"
)

type MessageService interface {
	// ReadAllMessage 消息设置已读
	ReadAllMessage(ctx context.Context, userID int64) error
	// ReadMessage 消息设置已读
	ReadMessage(ctx context.Context, userID int64, ids []string) error
	// GetMessageCount 消息数量统计
	GetMessageCount(ctx context.Context, userID int64, state int) (int64, error)
	// SendWebsocketMessage 发送websocket消息
	SendWebsocketMessage(ctx context.Context, msg *dto.WebsocketMessage) error
	// GetMessageList 获取消息列表
	GetMessageList(ctx context.Context, msg *dto.MessagePage) ([]*model.Message, int64, error)
}

type EmailService interface {
	GetEmail(ctx context.Context) (*dto.EmailConfig, error)
	SetEmail(ctx context.Context, email *dto.EmailConfig) error
	SendEmail(ctx context.Context, receiver, emailTemplate, jsonData string) error
}
