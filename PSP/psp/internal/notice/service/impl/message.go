package impl

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/service"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type messageServiceImpl struct {
	messageDao dao.MessageDao
}

func NewMessageService() (service.MessageService, error) {
	messageDao, err := dao.NewMessageDao()
	if err != nil {
		return nil, err
	}

	messageService := &messageServiceImpl{
		messageDao: messageDao,
	}

	return messageService, nil
}

// ReadAllMessage 消息设置已读
func (s *messageServiceImpl) ReadAllMessage(ctx context.Context, userID int64) error {
	logger := logging.GetLogger(ctx)

	if err := s.messageDao.ReadAllMessage(ctx, snowflake.ID(userID)); err != nil {
		logger.Errorf("batch read message err: %v", err)
		return err
	}

	return nil
}

// ReadMessage 消息设置已读
func (s *messageServiceImpl) ReadMessage(ctx context.Context, userID int64, ids []string) error {
	logger := logging.GetLogger(ctx)

	sids, fids := snowflake.BatchParseString(ids)
	if len(fids) > 0 {
		logger.Warnf("message ids [%v] is invalid, cannot parsed to snowflake id", fids)
	}

	if err := s.messageDao.ReadMessage(ctx, snowflake.ID(userID), sids); err != nil {
		logger.Errorf("batch read message err: %v", err)
		return err
	}

	return nil
}

// SendWebsocketMessage 发送websocket消息
func (s *messageServiceImpl) SendWebsocketMessage(ctx context.Context, msg *dto.WebsocketMessage) error {
	logger := logging.GetLogger(ctx)

	uid := snowflake.MustParseString(msg.UserId)
	message := &model.Message{
		UserId:  uid,
		Type:    msg.Type,
		Content: msg.Content,
	}

	msgID, err := s.messageDao.SaveAndSendMessage(ctx, message)
	if err != nil {
		logger.Errorf("save and send websocket message err: %v", err)
		return err
	}

	logger.Debugf("save and send websocket message success, message id: [%v]", msgID)

	return nil
}

// GetMessageList 获取消息列表
func (s *messageServiceImpl) GetMessageList(ctx context.Context, msg *dto.MessagePage) ([]*model.Message, int64, error) {
	logger := logging.GetLogger(ctx)

	messages, total, err := s.messageDao.GetMessageList(ctx, msg)
	if err != nil {
		logger.Errorf("get message list err: %v", err)
		return nil, 0, err
	}

	return messages, total, nil
}

// GetMessageCount 消息数量统计
func (s *messageServiceImpl) GetMessageCount(ctx context.Context, userID int64, state int) (int64, error) {
	logger := logging.GetLogger(ctx)

	count, err := s.messageDao.GetMessageCount(ctx, snowflake.ID(userID), state)
	if err != nil {
		logger.Errorf("get message count err: %v", err)
		return 0, err
	}

	return count, nil
}
