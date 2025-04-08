package rpc

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/util"
)

// SendWebsocketMessage 发送websocket消息
func (s *GRPCService) SendWebsocketMessage(ctx context.Context, req *pb.WebsocketMessage) (*pb.NoticeEmpty, error) {
	logger := logging.GetLogger(ctx)

	message := util.ConvertGRPCMessage(req)
	if message == nil {
		return nil, status.Errorf(errcode.ErrNoticeFailSend, "the send message is empty")
	}

	if err := s.messageService.SendWebsocketMessage(ctx, message); err != nil {
		logger.Errorf("send websocket message err: %v", err)
		return nil, status.Errorf(errcode.ErrNoticeFailSend, "failed to send websocket message")
	}

	return &pb.NoticeEmpty{}, nil
}
