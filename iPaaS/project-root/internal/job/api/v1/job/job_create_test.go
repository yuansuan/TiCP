package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type JobCreateSuite struct {
	JobSuite

	// data
	mockJob *models.Job
	req     jobcreate.Request
}

func TestJobCreate(t *testing.T) {
	suite.Run(t, new(JobCreateSuite))
}
func (s *JobCreateSuite) SetupSuite() {
	s.JobSuite.SetupSuite()
	s.mockConfig()
}

func (s *JobCreateSuite) SetupTest() {
	s.JobSuite.SetupTest()
	s.mockJob = &models.Job{
		ID:     snowflake.ID(12345),
		UserID: snowflake.ID(54321),
	}
}

func (s *JobCreateSuite) SetupSubTest() {
	s.JobSuite.SetupSubTest()
	cores := 1
	memory := 1
	s.req = jobcreate.Request{
		Name: "test name",
		Params: jobcreate.Params{
			Application: jobcreate.Application{
				Command: "sleep 10",
				AppID:   "4WLUvKxq7S5",
			},
			Resource: &jobcreate.Resource{
				Cores:  &cores,
				Memory: &memory,
			},
			EnvVars: map[string]string{},
			Input: &jobcreate.Input{
				Type:        "cloud_storage",
				Source:      "https://jn_storage_endpoint:8899/h9z/input",
				Destination: "",
			},
			Output: &jobcreate.Output{
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
	}
}

type jobCreateTestCase struct {
	JobTestCase
	setHeaderFunc func(ctx *gin.Context)
	setReq        func()
}

func (s *JobCreateSuite) reqFunc(setHeaderFunc func(ctx *gin.Context), setReq func()) func() {
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

func (s *JobCreateSuite) doFunc() {
	s.handler.Create(s.ctx)
}

func (s *JobCreateSuite) respFunc(tc *jobCreateTestCase) func() {
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

func (s *JobCreateSuite) MostCasesMockExpectFunc() {
	s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
		PublishStatus: string(update.Published),
		Image:         "image",
	}, nil)
	s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), gomock.Any()).Return(&schema.ApplicationAllow{
		ID:    "5ajX8WkAXdd",
		AppID: "4WLUvKxq7S5",
	}, fmt.Errorf("not allow error"))
	s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.ApplicationQuota{
		AppID:  "4WLUvKxq7S5",
		UserID: "h9z",
	}, nil)
}

func (s *JobCreateSuite) TestCreate() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
					s.mockJobSrv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("4ER", nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job nocommand success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus:   string(update.Published),
						Image:           "image",
						Command:         "sleep 10",
						ExtentionParams: "{\"YS_MAIN_FILE\":{\"Type\":\"String\",\"ReadableName\":\"主文件\",\"Must\":true},\"LIST\":{\"Type\":\"StringList\",\"ReadableName\":\"List\",\"Values\":[\"one\",\"two\"]}}",
					}, nil)
					s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), gomock.Any()).Return(&schema.ApplicationAllow{
						ID:    "5ajX8WkAXdd",
						AppID: "4WLUvKxq7S5",
					}, fmt.Errorf("not allow error"))
					s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.ApplicationQuota{
						AppID:  "4WLUvKxq7S5",
						UserID: "h9z",
					}, nil)
					s.mockJobSrv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("4ER", nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Application.Command = ""
				s.req.Params.EnvVars = map[string]string{"YS_MAIN_FILE": "123", "LIST": "one", "OTHER": "other"}
			},
		},
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}

func (s *JobCreateSuite) TestCreateValidate() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job bind error",
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
				Name:              "test create job invalid user id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidUserID,
			},
		},
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}

func (s *JobCreateSuite) TestCreateValidateApplication() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job app id invalid",
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
				Name:              "test create job app not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.AppIDNotFoundErrorCode,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(nil, xorm.ErrNotExist)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job app info error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job app not published",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.AppNotPublished,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus: string(update.Unpublished),
						Image:         "image",
					}, nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job app quota not found",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.UserNoAppQuota,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus: string(update.Published),
						Image:         "image",
					}, nil)
					s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), gomock.Any()).Return(&schema.ApplicationAllow{
						ID:    "5ajX8WkAXdd",
						AppID: "4WLUvKxq7S5",
					}, fmt.Errorf("not allow error"))
					s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job app image and binpath empty",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus: string(update.Published),
					}, nil)
					s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), gomock.Any()).Return(&schema.ApplicationAllow{
						ID:    "5ajX8WkAXdd",
						AppID: "4WLUvKxq7S5",
					}, fmt.Errorf("not allow error"))
					s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.ApplicationQuota{
						AppID:  "4WLUvKxq7S5",
						UserID: "h9z",
					}, nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job command is too long",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentCommand,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Application.Command = strings.Repeat("a", consts.MaxCommandLength+1)
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job and app command is empty",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Application.Command = ""
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create parse ExtentionParams error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus: string(update.Published),
						Image:         "image",
						Command:       "sleep 10",
					}, nil)
					s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), gomock.Any()).Return(&schema.ApplicationAllow{
						ID:    "5ajX8WkAXdd",
						AppID: "4WLUvKxq7S5",
					}, fmt.Errorf("not allow error"))
					s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.ApplicationQuota{
						AppID:  "4WLUvKxq7S5",
						UserID: "h9z",
					}, nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Application.Command = ""
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create missing required environment variable",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentEnv,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus:   string(update.Published),
						Image:           "image",
						Command:         "sleep 10",
						ExtentionParams: "{\"YS_MAIN_FILE\":{\"Type\":\"String\",\"ReadableName\":\"主文件\",\"Must\":true}}",
					}, nil)
					s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), gomock.Any()).Return(&schema.ApplicationAllow{
						ID:    "5ajX8WkAXdd",
						AppID: "4WLUvKxq7S5",
					}, fmt.Errorf("not allow error"))
					s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.ApplicationQuota{
						AppID:  "4WLUvKxq7S5",
						UserID: "h9z",
					}, nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Application.Command = ""
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create environment variable not in the allowed range",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentEnv,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus:   string(update.Published),
						Image:           "image",
						Command:         "sleep 10",
						ExtentionParams: "{\"YS_MAIN_FILE\":{\"Type\":\"StringList\",\"ReadableName\":\"主文件\",\"Values\":[\"abc\"]}}",
					}, nil)
					s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), gomock.Any()).Return(&schema.ApplicationAllow{
						ID:    "5ajX8WkAXdd",
						AppID: "4WLUvKxq7S5",
					}, fmt.Errorf("not allow error"))
					s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.ApplicationQuota{
						AppID:  "4WLUvKxq7S5",
						UserID: "h9z",
					}, nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Application.Command = ""
				s.req.Params.EnvVars = map[string]string{"YS_MAIN_FILE": "123"}
			},
		},
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}

func (s *JobCreateSuite) TestCreateValidateResource() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job core less than Min",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.QuotaExhaustedResource,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				*s.req.Params.Resource.Cores = -1
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job memory less than Min",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.QuotaExhaustedResource,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				*s.req.Params.Resource.Memory = -1
			},
		},
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}

func (s *JobCreateSuite) TestCreateValidateInputFile() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file input type error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentInput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Input.Type = "xxx"
			},
		},
		// source
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file input path empty",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentErrorCode, // TODO: InvalidArgumentInput
				MockExpectFunc: func() {
					// s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
					// 	PublishStatus: string(update.Published),
					// 	Image:         "image",
					// }, nil)
					// s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.ApplicationQuota{
					// 	AppID:  "4WLUvKxq7S5",
					// 	UserID: "h9z",
					// }, nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Input.Source = ""
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file input not absolute with domain",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentInput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Input.Source = "xxx"
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file input not in zone list",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentInput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Input.Source = "http://xxx"
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file input ys_id empty",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentInput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Input.Source = "https://jn_storage_endpoint:8899/"
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file input ys_id invalid",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentInput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Input.Source = "https://jn_storage_endpoint:8899/OOO"
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file input ys_id Unauthorized",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobPathUnauthorized,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Input.Source = "https://jn_storage_endpoint:8899/4mvp"
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file input with \"..\"",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentInput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Input.Source = "https://jn_storage_endpoint:8899/h9z/../"
			},
		},
		// destination
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file input dest with domain",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentInput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Input.Source = "https://jn_storage_endpoint:8899/h9z/input"
				s.req.Params.Input.Destination = "https://jn_storage_endpoint:8899/h9z/dest"
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file input dest begin with '/'",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentInput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Input.Source = "https://jn_storage_endpoint:8899/h9z/input"
				s.req.Params.Input.Destination = "/h9z/dest"
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file input dest with \"..\"",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentInput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Input.Source = "https://jn_storage_endpoint:8899/h9z/input"
				s.req.Params.Input.Destination = "h9z/../"
			},
		},
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}

func (s *JobCreateSuite) TestCreateValidateOutputFile() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file output type error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentOutput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Output = &jobcreate.Output{
					Type:          "xxx",
					Address:       "",
					NoNeededPaths: "",
				}
			},
		},
		// address
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file output address empty",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentOutput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Output = &jobcreate.Output{
					Type:          "cloud_storage",
					Address:       "",
					NoNeededPaths: "",
				}
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file output address not absolute with domain",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentOutput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Output = &jobcreate.Output{
					Type:          "cloud_storage",
					Address:       "xxx",
					NoNeededPaths: "",
				}
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file output address not in zone list",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentOutput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Output = &jobcreate.Output{
					Type:          "cloud_storage",
					Address:       "http://xxx",
					NoNeededPaths: "",
				}
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file output address ys_id empty",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentOutput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Output = &jobcreate.Output{
					Type:          "cloud_storage",
					Address:       "https://jn_storage_endpoint:8899/",
					NoNeededPaths: "",
				}
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file output address ys_id invalid",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentOutput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Output = &jobcreate.Output{
					Type:          "cloud_storage",
					Address:       "https://jn_storage_endpoint:8899/OOO",
					NoNeededPaths: "",
				}
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file output address Unauthorized",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobPathUnauthorized,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Output = &jobcreate.Output{
					Type:          "cloud_storage",
					Address:       "https://jn_storage_endpoint:8899/4mvp",
					NoNeededPaths: "",
				}
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job file output address with \"..\"",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentOutput,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.Output = &jobcreate.Output{
					Type:          "cloud_storage",
					Address:       "https://jn_storage_endpoint:8899/h9z/../",
					NoNeededPaths: "",
				}
			},
		},
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}

func (s *JobCreateSuite) TestCreateValidateCustomRule() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job custom rule key too long",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentCustomStateRuleKeyStatement,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.CustomStateRule = &jobcreate.CustomStateRule{
					KeyStatement: strings.Repeat("a", consts.MaxCustomStateRuleLength+1),
					ResultState:  jobcreate.ResultStateFailed,
				}
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job custom rule ResultState is invalid",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentCustomStateRuleResultState,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Params.CustomStateRule = &jobcreate.CustomStateRule{
					KeyStatement: "error",
					ResultState:  "xxx",
				}
			},
		},
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}

func (s *JobCreateSuite) TestCreateValidateChargeParams() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job charge type is invalid",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentChargeParams,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				chargeType := schema.ChargeType("xxx")
				s.req.ChargeParam = schema.ChargeParams{
					ChargeType: &chargeType,
				}
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job charge type is unsupported",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentChargeParams,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				chargeType := schema.PrePaid
				s.req.ChargeParam = schema.ChargeParams{
					ChargeType: &chargeType,
				}
			},
		},
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}

func (s *JobCreateSuite) TestCreateValidateShared() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job shared IsYsProductUser error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
					s.mockUserChecker.EXPECT().IsYsProductUser(gomock.Any()).Return(false, errors.New("error"))
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.NoRound = true
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job shared not ys product user",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobAccessDenied,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
					s.mockUserChecker.EXPECT().IsYsProductUser(gomock.Any()).Return(false, nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.NoRound = true
			},
		},
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}

func (s *JobCreateSuite) TestCreateValidateAllocType() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job alloc type IsYsProductUser error", // 检查用户是否是 YsProductUser 时发生错误
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus: string(update.Published),
						Image:         "image",
					}, nil)
					s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), gomock.Any()).Return(&schema.ApplicationAllow{
						ID:    "5ajX8WkAXdd",
						AppID: "4WLUvKxq7S5",
					}, fmt.Errorf("not allow error"))
					s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.ApplicationQuota{
						AppID:  "4WLUvKxq7S5",
						UserID: "h9z",
					}, nil)
					s.mockUserChecker.EXPECT().IsYsProductUser(gomock.Any()).Return(false, errors.New("error"))
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
				Name:              "test create job alloc type not ys product user",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobAccessDenied,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus: string(update.Published),
						Image:         "image",
					}, nil)
					s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), gomock.Any()).Return(&schema.ApplicationAllow{
						ID:    "5ajX8WkAXdd",
						AppID: "4WLUvKxq7S5",
					}, fmt.Errorf("not allow error"))
					s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.ApplicationQuota{
						AppID:  "4WLUvKxq7S5",
						UserID: "h9z",
					}, nil)
					s.mockUserChecker.EXPECT().IsYsProductUser(gomock.Any()).Return(false, nil)
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
				Name:              "test create job alloc type success with ys product user",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
					s.mockUserChecker.EXPECT().IsYsProductUser(gomock.Any()).Return(true, nil)
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
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}

func (s *JobCreateSuite) TestCreateValidateOther() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job name is too long",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentName,
				MockExpectFunc: func() {
					s.mockAppASrv.EXPECT().GetApp(gomock.Any(), gomock.Any()).Return(&models.Application{
						PublishStatus: string(update.Published),
						Image:         "image",
					}, nil)
					// s.mockAp
					s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), gomock.Any()).Return(&schema.ApplicationAllow{
						ID:    "5ajX8WkAXdd",
						AppID: "4WLUvKxq7S5",
					}, fmt.Errorf("not allow error"))
					s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.ApplicationQuota{
						AppID:  "4WLUvKxq7S5",
						UserID: "h9z",
					}, nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Name = strings.Repeat("a", consts.MaxNameLength+1)
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job comment is too long",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentComment,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Comment = strings.Repeat("a", consts.MaxCommentLength+1)
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job zone is unknown",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentZone,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Zone = "az-zhigu"
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job zone hpc endpoint is empty",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentZone,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
			setReq: func() {
				s.req.Zone = "az-wuxi"
			},
		},
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}

func (s *JobCreateSuite) TestCreateServiceError() {
	testCase := []jobCreateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test create job service internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
					s.mockJobSrv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("error"))
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job service invalid account id",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalErrorInvalidAccountId,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
					s.mockJobSrv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", common.ErrInvalidAccountId)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job service account not enough balance",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.InvalidAccountStatusNotEnoughBalance,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
					s.mockJobSrv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", common.ErrInvalidAccountStatusNotEnoughBalance)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test create job service account frozen",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.InvalidAccountStatusFrozen,
				MockExpectFunc: func() {
					s.MostCasesMockExpectFunc()
					s.mockJobSrv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", common.ErrInvalidAccountStatusFrozen)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCase {
		testJobCreate := tc.MakeTestFunc(s.reqFunc(tc.setHeaderFunc, tc.setReq), s.doFunc, s.respFunc(&tc))
		testJobCreate(s.T(), &s.JobSuite)
	}
}
