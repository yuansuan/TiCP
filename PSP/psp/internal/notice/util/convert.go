package util

import (
	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dto"
)

func ConvertMessage(msg *model.Message) *dto.Message {
	if msg != nil {
		return &dto.Message{
			ID:         msg.Id.String(),
			UserID:     msg.UserId.String(),
			Type:       msg.Type,
			Content:    msg.Content,
			State:      msg.State,
			CreateTime: msg.CreateTime,
			UpdateTime: msg.UpdateTime,
		}
	}

	return nil
}

func ConvertGRPCMessage(msg *pb.WebsocketMessage) *dto.WebsocketMessage {
	if msg != nil {
		return &dto.WebsocketMessage{
			UserId:  msg.UserId,
			Type:    msg.Type,
			Content: msg.Content,
		}
	}

	return nil
}

func ConvertModelMessage(msg *model.Message) *dto.WebsocketMessage {
	if msg != nil {
		return &dto.WebsocketMessage{
			UserId:  msg.UserId.String(),
			Type:    msg.Type,
			Content: msg.Content,
		}
	}

	return nil
}
