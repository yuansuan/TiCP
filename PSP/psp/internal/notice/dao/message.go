package dao

import (
	"context"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/service/mq"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/dbutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

const (
	StateUnRead = 1
	StateRead   = 2
)

type messageDaoImpl struct {
	sid      *snowflake.Node
	producer *mq.KafkaProducer
}

// NewMessageDao 创建MessageDao
func NewMessageDao() (MessageDao, error) {
	node, err := snowflake.GetInstance()
	if err != nil {
		return nil, err
	}

	producer, err := mq.NewKafkaProducer(common.NoticeWebsocketTopic)
	if err != nil {
		return nil, err
	}

	return &messageDaoImpl{
		sid:      node,
		producer: producer,
	}, nil
}

// SaveAndSendMessage 保存并发送消息
func (d *messageDaoImpl) SaveAndSendMessage(ctx context.Context, message *model.Message) (snowflake.ID, error) {
	msgID, err := boot.MW.DefaultTransaction(ctx, func(session *xorm.Session) (interface{}, error) {
		message.Id = d.sid.Generate()
		message.State = StateUnRead
		_, err := session.InsertOne(message)
		if err != nil {
			return nil, err
		}

		msg := util.ConvertModelMessage(message)
		if err = d.producer.SendMessage(ctx, common.NoticeWebsocketKey, msg); err != nil {
			return nil, err
		}

		return message.Id, nil
	})

	if err != nil {
		return 0, err
	}

	return msgID.(snowflake.ID), nil
}

// InsertMessage 保存消息信息
func (d *messageDaoImpl) InsertMessage(ctx context.Context, message *model.Message) error {
	session := boot.MW.DefaultSession(ctx)

	message.Id = d.sid.Generate()
	_, err := session.InsertOne(message)
	if err != nil {
		return err
	}

	return nil
}

// ReadMessage 消息设置已读
func (d *messageDaoImpl) ReadMessage(ctx context.Context, userID snowflake.ID, ids []snowflake.ID) error {
	session := boot.MW.DefaultSession(ctx)

	message := &model.Message{
		State: StateRead,
	}
	_, err := session.Where("user_id=?", userID).In("id", ids).Cols("state").Update(message)
	if err != nil {
		return err
	}

	return nil
}

// ReadAllMessage 消息设置已读
func (d *messageDaoImpl) ReadAllMessage(ctx context.Context, userID snowflake.ID) error {
	session := boot.MW.DefaultSession(ctx)

	message := &model.Message{
		State: StateRead,
	}
	_, err := session.Where("user_id=?", userID).Cols("state").Update(message)
	if err != nil {
		return err
	}

	return nil
}

// GetMessageCount 获取消息数量统计
func (d *messageDaoImpl) GetMessageCount(ctx context.Context, userID snowflake.ID, state int) (int64, error) {
	session := boot.MW.DefaultSession(ctx)

	session.Where("user_id=?", userID)
	if state > 0 {
		session.Where("state=?", state)
	}

	messages := &model.Message{}
	total, err := session.Count(messages)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// GetMessageList 获取消息列表
func (d *messageDaoImpl) GetMessageList(ctx context.Context, msg *dto.MessagePage) ([]*model.Message, int64, error) {
	session := boot.MW.DefaultSession(ctx)

	var messages []*model.Message
	index, size := msg.Page.Index, msg.Page.Size
	offset, err := xtype.GetPageOffset(index, size)
	if err != nil {
		return nil, 0, err
	}

	session.Where("user_id=?", msg.UserID)
	total, err := wrapListSession(session, msg.Filter, msg.OrderSort).Limit(int(size), int(offset)).FindAndCount(&messages)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func wrapListSession(session *xorm.Session, filter *dto.MessageFilter, orderSort *xtype.OrderSort) *xorm.Session {
	filterSession := wrapFilterSession(session, filter)
	return dbutil.WrapSortSession(filterSession, orderSort)
}

func wrapFilterSession(session *xorm.Session, filter *dto.MessageFilter) *xorm.Session {
	if filter != nil {
		if filter.State > 0 {
			session.Where("state=?", filter.State)
		}

		if !strutil.IsEmpty(filter.Content) {
			session.Where("content like ?", "%"+filter.Content+"%")
		}
	}

	return session
}
