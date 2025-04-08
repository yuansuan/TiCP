package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/allowadd"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/allowdelete"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	srvv1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

// suite
type ApplicationAllowSuite struct {
	// suite
	suite.Suite

	// base
	ctrl *gomock.Controller

	w      *httptest.ResponseRecorder
	ctx    *gin.Context
	logger *logging.Logger

	// service
	mockSrv         *srvv1.MockService
	mockAppSrv      *srvv1.MockAppSrv
	mockAppAllowSrv *srvv1.MockAppAllowSrv

	// controller
	controller *Controller

	// data
	mockAppAllow *schema.ApplicationAllow
}

func TestApplicationAllowSuite(t *testing.T) {
	suite.Run(t, new(ApplicationAllowSuite))
}

func (s *ApplicationAllowSuite) SetupSuite() {
	s.T().Log("setup suite")

	// mock base controller
	s.ctrl = gomock.NewController(s.T())

	// mock gin context
	s.w = httptest.NewRecorder()
	s.ctx = mock.GinContext(s.w)

	// logger
	logger, err := logging.NewLogger(logging.WithDefaultLogConfigOption())
	if !s.NoError(err) {
		return
	}
	s.logger = logger
	s.ctx.Set(logging.LoggerName, s.logger)

	// mock service
	mockSrv := srvv1.NewMockService(s.ctrl)
	mockAppSrv := srvv1.NewMockAppSrv(s.ctrl)
	mockAppAllowSrv := srvv1.NewMockAppAllowSrv(s.ctrl)

	mockSrv.EXPECT().Apps().Return(mockAppSrv).AnyTimes()
	mockSrv.EXPECT().AppsAllow().Return(mockAppAllowSrv).AnyTimes()

	s.mockSrv = mockSrv
	s.mockAppSrv = mockAppSrv
	s.mockAppAllowSrv = mockAppAllowSrv

	// controller
	s.controller = &Controller{srv: s.mockSrv}
}

func (s *ApplicationAllowSuite) TearDownTest() {
	s.T().Log("teardown suite")
}

func (s *ApplicationAllowSuite) SetupTest() {
	s.T().Log("setup test")
	s.w = httptest.NewRecorder()
	s.ctx = mock.GinContext(s.w)
	s.ctx.Set(logging.LoggerName, s.logger)
	s.ctx.Request.Header.Set("x-ys-request-id", gofakeit.UUID())

	mockAppAllow := &schema.ApplicationAllow{
		ID:    "4Eb",
		AppID: "4ER",
	}

	s.mockAppAllow = mockAppAllow
	s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), snowflake.ID(12345)).Return(mockAppAllow, nil).AnyTimes()
	s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), snowflake.ID(44444)).Return(nil, common.ErrAppAllowNotFound).AnyTimes()
	s.mockAppAllowSrv.EXPECT().GetAllow(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()

	s.mockAppAllowSrv.EXPECT().AddAllow(gomock.Any(), snowflake.ID(12345)).Return(mockAppAllow, nil).AnyTimes()
	s.mockAppAllowSrv.EXPECT().AddAllow(gomock.Any(), snowflake.ID(44444)).Return(nil, common.ErrAppIDNotFound).AnyTimes()
	s.mockAppAllowSrv.EXPECT().AddAllow(gomock.Any(), snowflake.ID(33333)).Return(nil, common.ErrAppAllowAlreadyExist).AnyTimes()
	s.mockAppAllowSrv.EXPECT().AddAllow(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()

	s.mockAppAllowSrv.EXPECT().DeleteAllow(gomock.Any(), snowflake.ID(12345)).Return(nil).AnyTimes()
	s.mockAppAllowSrv.EXPECT().DeleteAllow(gomock.Any(), snowflake.ID(44444)).Return(common.ErrAppAllowNotFound).AnyTimes()
	s.mockAppAllowSrv.EXPECT().DeleteAllow(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error")).AnyTimes()
}

func (s *ApplicationAllowSuite) TearDownSuite() {
	s.T().Log("teardown suite")
}

func (s *ApplicationAllowSuite) SetupSubTest() {
	s.T().Log("setup sub test")
	s.w = httptest.NewRecorder()
	s.ctx = mock.GinContext(s.w)
	s.ctx.Set(logging.LoggerName, s.logger)
	s.ctx.Request.Header.Set("x-ys-request-id", gofakeit.UUID())
}

func (s *ApplicationAllowSuite) TearDownSubTest() {
	s.T().Log("teardown sub test")
}

func (s *ApplicationAllowSuite) TestController_AllowGet() {
	s.Run("success", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		// configure query params
		u := url.Values{}

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)
		s.controller.AllowGet(s.ctx)
		s.Equal(http.StatusOK, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}

		s.T().Log(spew.Sdump(resp))
	})

	s.Run("empty AppID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: ""},
		}
		// configure query params
		u := url.Values{}

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)
		s.controller.AllowGet(s.ctx)
		s.Equal(http.StatusBadRequest, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidAppID, resp.ErrorCode)
	})

	s.Run("error AppID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "#$@@#$"},
		}
		// configure query params
		u := url.Values{}
		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)
		s.controller.AllowGet(s.ctx)
		s.Equal(http.StatusBadRequest, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidAppID, resp.ErrorCode)
	})

	s.Run("allow not found", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "edh"}, // 44444
		}
		// configure query params
		u := url.Values{}
		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)
		s.controller.AllowGet(s.ctx)
		s.Equal(http.StatusNotFound, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.AppAllowNotFound, resp.ErrorCode)
	})

	s.Run("internal error", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "h9z"}, // 54321
		}
		// configure query params
		u := url.Values{}
		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)
		s.controller.AllowGet(s.ctx)
		s.Equal(http.StatusInternalServerError, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InternalServerErrorCode, resp.ErrorCode)
	})
}

func (s *ApplicationAllowSuite) TestController_AllowAdd() {
	s.Run("success", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		mock.HTTPRequest(s.ctx, http.MethodPatch, nil, params, nil)
		s.controller.AllowAdd(s.ctx)
		s.Equal(http.StatusOK, s.w.Code)
		got := s.w.Body.String()
		s.T().Log(spew.Sdump(got))
	})

	s.Run("empty AppID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: ""},
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, nil, params, nil)
		s.controller.AllowAdd(s.ctx)
		s.Equal(http.StatusBadRequest, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidAppID, resp.ErrorCode)
	})

	s.Run("error AppID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "#$@@#$"},
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, nil, params, nil)
		s.controller.AllowAdd(s.ctx)
		s.Equal(http.StatusBadRequest, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidAppID, resp.ErrorCode)
	})

	s.Run("app not found", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "edh"}, // 不存在的app(44444)
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, nil, params, nil)
		s.controller.AllowAdd(s.ctx)
		s.Equal(http.StatusNotFound, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.AppIDNotFoundErrorCode, resp.ErrorCode)
	})

	s.Run("allow already exists", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "aUH"}, // 33333
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, nil, params, nil)
		s.controller.AllowAdd(s.ctx)
		s.Equal(http.StatusConflict, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.AppAllowAlreadyExist, resp.ErrorCode)
	})

	s.Run("internal error", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "h9z"}, // 54321
		}
		mock.HTTPRequest(s.ctx, http.MethodPatch, nil, params, nil)
		s.controller.AllowAdd(s.ctx)
		s.Equal(http.StatusInternalServerError, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InternalServerErrorCode, resp.ErrorCode)
	})
}

func (s *ApplicationAllowSuite) TestController_AllowDelete() {
	s.Run("success", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}

		mock.HTTPRequest(s.ctx, http.MethodDelete, nil, params, nil)
		s.controller.AllowDelete(s.ctx)
		s.Equal(http.StatusOK, s.w.Code)
		got := s.w.Body.String()
		s.T().Log(spew.Sdump(got))
	})

	s.Run("empty AppID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: ""},
		}
		mock.HTTPRequest(s.ctx, http.MethodDelete, nil, params, nil)
		s.controller.AllowDelete(s.ctx)
		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp allowdelete.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidAppID, resp.ErrorCode)
	})

	s.Run("error AppID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "#$@@#$"},
		}

		mock.HTTPRequest(s.ctx, http.MethodDelete, nil, params, nil)
		s.controller.AllowDelete(s.ctx)
		s.Equal(http.StatusBadRequest, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowdelete.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidAppID, resp.ErrorCode)
	})

	s.Run("allow not found", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "edh"}, // 44444
		}

		mock.HTTPRequest(s.ctx, http.MethodDelete, nil, params, nil)
		s.controller.AllowDelete(s.ctx)
		s.Equal(http.StatusNotFound, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowdelete.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.AppAllowNotFound, resp.ErrorCode)
	})

	s.Run("internal error", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "h9z"}, // 54321
		}

		mock.HTTPRequest(s.ctx, http.MethodDelete, nil, params, nil)
		s.controller.AllowDelete(s.ctx)
		s.Equal(http.StatusInternalServerError, s.w.Code)
		got := s.w.Body.Bytes()
		var resp allowdelete.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InternalServerErrorCode, resp.ErrorCode)
	})
}
