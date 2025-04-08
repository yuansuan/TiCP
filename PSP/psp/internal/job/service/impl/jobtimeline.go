package impl

import (
	"context"
	"strings"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
)

// GetJobTimeline 获取作业时间线
func (s *jobServiceImpl) GetJobTimeline(ctx context.Context, jobID, uploadFileTaskID, jobState, dataState string) ([]*dto.JobTimeLine, error) {
	logger := logging.GetLogger(ctx)

	sid, err := snowflake.ParseString(jobID)
	if err != nil {
		return nil, err
	}

	timelines, err := s.jobTimelineDao.GetJobTimeline(ctx, sid)
	if err != nil {
		logger.Errorf("get job [%v] timeline err: %v", jobID, err)
		return nil, err
	}

	jobTimeLines := make([]*dto.JobTimeLine, 0, len(timelines))
	for _, timeline := range timelines {
		jobTimeLine := &dto.JobTimeLine{
			EventName: timeline.EventName,
			EventTime: timeutil.DefaultFormatTime(timeline.EventTime),
		}

		if strings.HasSuffix(timeline.EventName, consts.JobStateBursting) {
			var uploadProgress int
			switch dataState {
			case consts.JobDataStateUploadFailed:
				uploadProgress = -1
			default:
				uploadProgress = 100
			}

			if uploadFileTaskID != "" && dataState == consts.JobDataStateUploading && jobState == consts.JobStateBursting {
				uploadFileTaskStatusReq := &pb.GetUploadHpcFileTaskStatusReq{
					TaskKey: uploadFileTaskID,
				}
				taskResp, err := s.rpc.Storage.GetUploadHpcFileTaskStatus(ctx, uploadFileTaskStatusReq)
				if err != nil {
					logger.Errorf("get job [%v] upload file task status err: %v", jobID, err)
					return nil, err
				}

				if taskResp.TotalSize > 0 {
					uploadProgress = int(float64(taskResp.CurrentSize) / float64(taskResp.TotalSize) * 100)
				}
			}

			jobTimeLine.Progress = uploadProgress
		} else {
			jobTimeLine.Progress = -1
		}

		jobTimeLines = append(jobTimeLines, jobTimeLine)
	}

	return jobTimeLines, nil
}
