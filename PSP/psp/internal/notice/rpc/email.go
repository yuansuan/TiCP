package rpc

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
)

func (s *GRPCService) SendEmail(ctx context.Context, in *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	logger := logging.GetLogger(ctx)

	err := s.emailService.SendEmail(ctx, in.Receiver, in.EmailTemplate, in.JsonData)
	if err != nil {
		logger.Errorf("send email err: %v, in: [%+v]", err, in)
		return nil, status.Errorf(errcode.ErrNoticeSendEmailFailed, "send email failed")
	}

	return &pb.SendEmailResponse{}, nil
}
