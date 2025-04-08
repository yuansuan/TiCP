package alarm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/wx"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/wx/markdown"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/wx/text"
)

func TestSendLongRunningJobAlarm(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSender := wx.NewMockSender(ctrl)
	ctx := context.Background()
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}
	ctx = context.WithValue(ctx, logging.LoggerName, logger)
	type args struct {
		ctx       context.Context
		job       *models.Job
		threshold int64
		mockFunc  func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal long running job",
			args: args{
				ctx: ctx,
				job: &models.Job{
					ID:                 snowflake.ID(1816762938233458688),
					Name:               "mock数据",
					Zone:               "az-jinan",
					AppID:              snowflake.ID(1816762964422692864),
					UserID:             snowflake.ID(1816762966540816384),
					ExecutionDuration:  256,
					ResourceAssignCpus: 56,
				},
				threshold: 10000,
				mockFunc: func() {
					mockSender.EXPECT().Send(gomock.Any()).Do(func(msg wx.WXMessage) {
						switch msg := msg.(type) {
						case *markdown.WXMarkdownMessage:
							t.Logf("markdown message: %s", spew.Sdump(msg))
						case *text.WXTextMessage:
							t.Logf("text message: %s", spew.Sdump(msg))
						}
					}).Return(&resty.Response{
						RawResponse: &http.Response{
							Status: "200 OK",
						},
					}, nil)
				},
			},
		},
		{
			name: "less than a threshold",
			args: args{
				ctx: ctx,
				job: &models.Job{
					ID:                 snowflake.ID(1816762938233458688),
					Name:               "mock数据",
					Zone:               "az-jinan",
					AppID:              snowflake.ID(1816762964422692864),
					UserID:             snowflake.ID(1816762966540816384),
					ExecutionDuration:  256,
					ResourceAssignCpus: 56,
				},
				threshold: 20000,
				mockFunc:  nil,
			},
		},
		{
			name: "send error",
			args: args{
				ctx: ctx,
				job: &models.Job{
					ID:                 snowflake.ID(1816762938233458688),
					Name:               "mock数据",
					Zone:               "az-jinan",
					AppID:              snowflake.ID(1816762964422692864),
					UserID:             snowflake.ID(1816762966540816384),
					ExecutionDuration:  256,
					ResourceAssignCpus: 56,
				},
				threshold: 10000,
				mockFunc: func() {
					mockSender.EXPECT().Send(gomock.Any()).Return(nil, errors.New("send error"))
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mockFunc != nil {
				tt.args.mockFunc()
			}
			SendLongRunningJobAlarm(tt.args.ctx, mockSender, tt.args.job, tt.args.threshold)
		})

	}
}
