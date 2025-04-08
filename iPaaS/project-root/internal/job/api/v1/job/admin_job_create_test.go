package job

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobcreate"
	job "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

// JobAdminCreateSuite
type JobAdminCreateSuite struct {
	JobSuite

	// data
	mockJob *models.Job
	req     jobcreate.Request
}

func TestJobAdminCreate(t *testing.T) {
	suite.Run(t, new(JobAdminCreateSuite))
}

func (s *JobAdminCreateSuite) SetupSuite() {
	s.JobSuite.SetupSuite()
	s.mockConfig()
}

func (s *JobAdminCreateSuite) SetupTest() {
	s.JobSuite.SetupTest()
	s.mockJob = &models.Job{
		ID:     snowflake.ID(12345),
		UserID: snowflake.ID(54321),
	}
}

// SetupSubTest
func (s *JobAdminCreateSuite) SetupSubTest() {
	s.JobSuite.SetupSubTest()
	cores := 1
	memory := 1
	s.req = jobcreate.Request{
		Request: job.Request{
			Name: "test name",
			Params: job.Params{
				Application: job.Application{
					Command: "sleep 10",
					AppID:   "4WLUvKxq7S5",
				},
				Resource: &job.Resource{
					Cores:  &cores,
					Memory: &memory,
				},
				EnvVars: map[string]string{},
				Input: &job.Input{
					Type:        "cloud_storage",
					Source:      "https://jn_storage_endpoint:8899/h9z/input",
					Destination: "",
				},
				Output: &job.Output{
					Type:          "cloud_storage",
					Address:       "https://jn_storage_endpoint:8899/h9z/output",
					NoNeededPaths: "",
				},
				TmpWorkdir:        false,
				SubmitWithSuspend: false,
				CustomStateRule:   nil,
			},
			Timeout:     0,
			Zone:        "az-jinan",
			Comment:     "comment",
			ChargeParam: schema.ChargeParams{},
			NoRound:     false,
			AllocType:   "",
		},
		Queue: "queue",
	}
}

type jobAdminCreateTestCase struct {
	JobTestCase
	setHeaderFunc func(ctx *gin.Context)
	setReq        func()
}

func (s *JobAdminCreateSuite) reqFunc(setHeaderFunc func(ctx *gin.Context), setReq func()) func() {
	return func() {
		// Set custom headers
		if setHeaderFunc != nil {
			setHeaderFunc(s.ctx)
		}

		// Set request body
		if setReq != nil {
			setReq()
		}

		mock.HTTPRequest(s.ctx, http.MethodPost, s.req, nil, nil)
	}
}

func (s *JobAdminCreateSuite) doFunc() {
	s.handler.AdminCreate(s.ctx)
}

func (s *JobAdminCreateSuite) respFunc(tc *jobAdminCreateTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp jobcreate.Response

		// unmarshal
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.T().Log(spew.Sdump(resp))

		tc.ErrorCode = resp.ErrorCode
		tc.ErrorMsg = resp.ErrorMsg
	}
}

func (s *JobAdminCreateSuite) TestAdminCreate() {
	testCase := []jobAdminCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test admin create job success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus: string(update.Published),
						Image:         "image",
					}, nil)
					s.mockJobSrv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("4ER", nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin create job bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentResource,
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Resource = nil
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin create job invalid user id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidUserID,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin create job app id invalid",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidAppID,
				MockExpectFunc:    func() {},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Application.AppID = "O" // O is invalid
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin create job service internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus: string(update.Published),
						Image:         "image",
					}, nil)
					s.mockJobSrv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("error"))
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin create job with average allocation",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus: string(update.Published),
						Image:         "image",
					}, nil)
					s.mockJobSrv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("4ER", nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.AllocType = "average"
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin create job with invalid allocation type",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.InvalidArgumentAllocType,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus: string(update.Published),
						Image:         "image",
					}, nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.AllocType = "auto" //非法的 AllocType 值: 随便填一个不是 average 的，返回403
			},
		},
	}

	for _, tc := range testCase {
		testJobAdminCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobAdminCreate(s.T(), &s.JobSuite)
	}
}
