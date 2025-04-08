package impl

import (
	"context"
	"fmt"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	openapijob "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

func (s *jobServiceImpl) GetJobSnapshotList(ctx context.Context, jobID string) (*dto.JobSnapshotListResponse, error) {
	logger := logging.GetLogger(ctx)

	exist, job, err := s.jobDao.GetJobDetail(ctx, snowflake.MustParseString(jobID))
	if err != nil {
		logger.Errorf("get the job: [%v] detail err: %v", jobID, err)
		return nil, err
	}
	if !exist {
		logger.Infof("the job: [%v] not found", jobID)
		return nil, JobNotFoundError
	}

	snapshotsData := make(map[string][]string)
	switch {
	case job.Type == common.Local:
		snapshots, err := openapijob.AdminListSnapshot(s.localAPI, job.OutJobId)
		if err != nil {
			logger.Errorf("get the admin job: [%v] snapshots err: %v", jobID, err)
			return nil, err
		}

		snapshotsData = *snapshots.Data
	default:
		return nil, fmt.Errorf("job compute type: [%v] not match", job.Type)
	}

	return &dto.JobSnapshotListResponse{Snapshots: snapshotsData}, nil
}

func (s *jobServiceImpl) GetJobSnapshot(ctx context.Context, jobID, path string) (*dto.JobSnapshotResponse, error) {
	logger := logging.GetLogger(ctx)

	exist, job, err := s.jobDao.GetJobDetail(ctx, snowflake.MustParseString(jobID))
	if err != nil {
		logger.Errorf("get the job: [%v] detail err: %v", jobID, err)
		return nil, err
	}
	if !exist {
		logger.Infof("the job: [%v] not found", jobID)
		return nil, JobNotFoundError
	}

	snapshotData := ""
	switch {
	case job.Type == common.Local:
		snapshot, err := openapijob.AdminGetSnapshot(s.localAPI, job.OutJobId, path)
		if err != nil {
			logger.Errorf("get the admin job: [%v] snapshot err: %v, path: [%v]", jobID, err, path)
			return nil, err
		}
		snapshotData = string(*snapshot.Data)
	default:
		return nil, fmt.Errorf("job compute type: [%v] not match", job.Type)
	}

	return &dto.JobSnapshotResponse{Snapshot: snapshotData}, nil
}
