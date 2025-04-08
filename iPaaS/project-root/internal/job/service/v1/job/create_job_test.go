package job

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

type CreateSuite struct {
	jobServiceSuite

	// data
	convert createConvertFunc
}

func TestCreate(t *testing.T) {
	suite.Run(t, new(CreateSuite))
}

type JobCreateTestCase struct {
	name           string
	expectedJobID  string
	expectedError  error
	mockExpectFunc func()
	setReq         func()
}

func (tc *JobCreateTestCase) Run(s *CreateSuite) {
	s.Run(tc.name, func() {
		// Mock EXPECT
		if tc.mockExpectFunc != nil {
			tc.mockExpectFunc()
		}

		// Set req
		if tc.setReq != nil {
			tc.setReq()
		}

		// do
		resp, err := s.jobSrv.Create(s.ctx, &jobcreate.Request{}, snowflake.Zero(), schema.ChargeParams{}, &models.Application{}, nil, s.convert)

		// assert
		if tc.expectedError != nil {
			if s.Error(err) {
				s.ErrorContains(err, tc.expectedError.Error())
			}
			return
		}

		if s.NoError(err) {
			s.Equal(tc.expectedJobID, resp)
		}
	})
}

func (s *CreateSuite) TestCreate() {
	testCases := []JobCreateTestCase{
		{
			name:          "normal",
			expectedJobID: snowflake.ID(10086).String(),
			expectedError: nil,
			mockExpectFunc: func() {
				mockJobID := snowflake.ID(10086)
				s.mockIDGen.EXPECT().GenID(gomock.Any()).Return(mockJobID, nil).Times(1)

				s.sqlmock.ExpectExec("INSERT INTO `job`").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			setReq: func() {
				s.convert = func(ctx context.Context, jobID snowflake.ID) (*models.Job, error) {
					return &models.Job{
						ID: snowflake.ID(10086),
					}, nil
				}
			},
		},
		{
			name:          "gen job id error",
			expectedJobID: "",
			expectedError: fmt.Errorf("gen job id error"),
			mockExpectFunc: func() {
				s.mockIDGen.EXPECT().GenID(gomock.Any()).Return(snowflake.ID(0), fmt.Errorf("gen job id error")).Times(1)
			},
		},
		{
			name:          "convert error",
			expectedJobID: "",
			expectedError: fmt.Errorf("convert error"),
			mockExpectFunc: func() {
				mockJobID := snowflake.ID(10086)
				s.mockIDGen.EXPECT().GenID(gomock.Any()).Return(mockJobID, nil).Times(1)
			},
			setReq: func() {
				s.convert = func(ctx context.Context, jobID snowflake.ID) (*models.Job, error) {
					return nil, fmt.Errorf("convert error")
				}
			},
		},
		{
			name: "insert error",
			mockExpectFunc: func() {
				mockJobID := snowflake.ID(10086)
				s.mockIDGen.EXPECT().GenID(gomock.Any()).Return(mockJobID, nil).Times(1)

				s.sqlmock.ExpectExec("INSERT INTO `job`").
					WillReturnError(fmt.Errorf("insert error"))
			},
			expectedError: fmt.Errorf("insert error"),
			setReq: func() {
				s.convert = func(ctx context.Context, jobID snowflake.ID) (*models.Job, error) {
					return &models.Job{
						ID: snowflake.ID(10086),
					}, nil
				}
			},
		},
	}

	for _, tc := range testCases {
		tc.Run(s)
	}
}
