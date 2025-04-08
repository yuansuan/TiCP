package job

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobbatchget"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

type JobBatchGetSuite struct {
	JobSuite

	// data
	mockJobs []*models.Job
}

func TestJobBatchGet(t *testing.T) {
	suite.Run(t, new(JobBatchGetSuite))
}

func (s *JobBatchGetSuite) SetupTest() {
	s.JobSuite.SetupTest()
	p := models.AdminParams{}
	ps, _ := json.Marshal(p)
	s.mockJobs = append(s.mockJobs, &models.Job{
		ID:     snowflake.ID(12345),
		UserID: snowflake.ID(54321),
		Params: string(ps),
	})
}

type jobBatchGetTestCase struct {
	JobTestCase
	req           *jobbatchget.Request
	setHeaderFunc func(ctx *gin.Context)
}

func (s *JobBatchGetSuite) reqFunc(req *jobbatchget.Request, setHeaderFunc func(ctx *gin.Context)) func() {
	return func() {
		// Set custom headers
		if setHeaderFunc != nil {
			setHeaderFunc(s.ctx)
		}

		mock.HTTPRequest(s.ctx, http.MethodPost, req, nil, nil)
	}
}

func (s *JobBatchGetSuite) doFunc() {
	s.handler.BatchGet(s.ctx)
}

func (s *JobBatchGetSuite) respFunc(tc *jobBatchGetTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp jobbatchget.Response

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

func (s *JobBatchGetSuite) TestBatchGet() {
	testCase := []jobBatchGetTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test batch get job success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().BatchGet(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.mockJobs, nil)
				},
			},
			req: &jobbatchget.Request{
				JobIDs: []string{"4ER", "O"}, // invalid job id "O" will be ignored
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCase {
		testJobBatchGet := tc.MakeTestFunc(s.reqFunc(tc.req, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobBatchGet(s.T(), &s.JobSuite)
	}
}

func (s *JobBatchGetSuite) TestBatchGetValidate() {
	testCase := []jobBatchGetTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test batch get job bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentErrorCode,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test batch get job invalid user id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidUserID,
			},
			req: &jobbatchget.Request{
				JobIDs: []string{"4ER"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test batch get job invalid jobs length",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentJobIDs,
			},
			req: &jobbatchget.Request{
				JobIDs: func() []string {
					return make([]string, 101)
				}(),
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCase {
		testJobBatchGet := tc.MakeTestFunc(s.reqFunc(tc.req, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobBatchGet(s.T(), &s.JobSuite)
	}
}
