package alarm

import (
	"bytes"
	"context"
	"text/template"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/wx"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/wx/text"
)

const LongRunningJobAlarmTemplate = `存在长时间运行作业: 
作业ID: {{.Job.ID}}
作业名称: {{.Job.Name}}
作业区域: {{.Job.Zone}}
应用ID: {{.Job.AppID}}
提交用户: {{.Job.UserID}}
运行时长: {{.Job.ExecutionDuration}} (秒)
分配核数: {{.Job.ResourceAssignCpus}} (核)
核秒: {{.CoreSeconds}} 超过阈值 {{.Threshold}}

请管理员及时关注和处理
`

type LongRunningJobAlarmData struct {
	Job         *models.Job
	CoreSeconds int64
	Threshold   int64
}

func SendLongRunningJobAlarm(ctx context.Context, sender wx.Sender, job *models.Job, threshold int64) {
	logger := logging.GetLogger(ctx)
	coreSeconds := job.ExecutionDuration * job.ResourceAssignCpus
	if coreSeconds < threshold {
		return
	}

	data := LongRunningJobAlarmData{
		Job:         job,
		CoreSeconds: coreSeconds,
		Threshold:   threshold,
	}

	// 写入template
	tmpl, err := template.New("longRunningJobAlarm").Parse(LongRunningJobAlarmTemplate)
	if err != nil {
		logger.Warnf("parse template failed, err: %v", err)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		logger.Warnf("execute template failed, err: %v", err)
		return
	}

	wxmsg := text.NewWXTextMessage(buf.String())

	// 发送告警
	resp, err := sender.Send(wxmsg)
	if err != nil {
		logger.Warnf("send wx message failed, err: %v", err)
		return
	}
	logger.Infof("send wx message success, resp: %v", resp)
}
