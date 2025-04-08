package job

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobneedsyncfile"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

type ListNeedSyncFileJobSuit struct {
	jobServiceSuite

	// data
}

func TestListNeedSyncFileJobs(t *testing.T) {
	suite.Run(t, new(ListNeedSyncFileJobSuit))
}

func (s *ListNeedSyncFileJobSuit) TestListNeedSyncFileJobs() {
	jobList := make([]*models.Job, 0)
	job := &models.Job{}

	jobList = append(jobList, job)

	s.mockJobDao.EXPECT().ListNeedFileSyncJobs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(jobList, int64(100), nil).Times(1)

	req := &jobneedsyncfile.Request{
		Zone:       "",
		PageOffset: new(int64),
		PageSize:   new(int64),
	}

	resp, err := s.jobSrv.ListNeedSyncFileJobs(s.ctx, req)
	if s.NoError(err) {
		return
	}

	s.Equal(int64(100), resp.Total)
	s.Equal(1, len(resp.Jobs))
}

func (s *ListNeedSyncFileJobSuit) TestListNeedSyncFileJobsError() {
	s.mockJobDao.EXPECT().ListNeedFileSyncJobs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, int64(0), fmt.Errorf("list job error"))

	req := &jobneedsyncfile.Request{
		Zone:       "",
		PageOffset: new(int64),
		PageSize:   new(int64),
	}

	_, err := s.jobSrv.ListNeedSyncFileJobs(s.ctx, req)
	s.Error(err)
	s.T().Log(err)
}
