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
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/quotaadd"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/quotadelete"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	srvv1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

// suite
type ApplicationQuotaSuite struct {
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
	mockAppQuotaSrv *srvv1.MockAppQuotaSrv

	// controller
	controller *Controller

	// data
	mockAppQuota *schema.ApplicationQuota
}

func TestApplicationQuotaSuite(t *testing.T) {
	suite.Run(t, new(ApplicationQuotaSuite))
}

func (s *ApplicationQuotaSuite) SetupSuite() {
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
	mockAppQuotaSrv := srvv1.NewMockAppQuotaSrv(s.ctrl)

	mockSrv.EXPECT().Apps().Return(mockAppSrv).AnyTimes()
	mockSrv.EXPECT().AppsQuota().Return(mockAppQuotaSrv).AnyTimes()

	s.mockSrv = mockSrv
	s.mockAppSrv = mockAppSrv
	s.mockAppQuotaSrv = mockAppQuotaSrv

	// controller
	s.controller = &Controller{srv: s.mockSrv}
}

func (s *ApplicationQuotaSuite) TearDownTest() {
	s.T().Log("teardown suite")
}

func (s *ApplicationQuotaSuite) SetupTest() {
	s.T().Log("setup test")
	s.w = httptest.NewRecorder()
	s.ctx = mock.GinContext(s.w)
	s.ctx.Set(logging.LoggerName, s.logger)
	s.ctx.Request.Header.Set("x-ys-request-id", gofakeit.UUID())

	mockAppQuota := &schema.ApplicationQuota{
		ID:     "4Eb",
		AppID:  "4ER",
		UserID: "h9z",
	}

	s.mockAppQuota = mockAppQuota
	s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), snowflake.ID(12345), snowflake.ID(54321)).Return(mockAppQuota, nil).AnyTimes()
	s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), snowflake.ID(12345), snowflake.ID(12345)).Return(nil, common.ErrAppQuotaNotFound).AnyTimes()
	s.mockAppQuotaSrv.EXPECT().GetQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()

	s.mockAppQuotaSrv.EXPECT().AddQuota(gomock.Any(), snowflake.ID(12345), snowflake.ID(54321)).Return(mockAppQuota, nil).AnyTimes()
	s.mockAppQuotaSrv.EXPECT().AddQuota(gomock.Any(), snowflake.ID(998877), gomock.Any()).Return(nil, common.ErrAppIDNotFound).AnyTimes()
	s.mockAppQuotaSrv.EXPECT().AddQuota(gomock.Any(), snowflake.ID(12345), snowflake.ID(998877)).Return(nil, common.ErrUserNotExists).AnyTimes()
	s.mockAppQuotaSrv.EXPECT().AddQuota(gomock.Any(), snowflake.ID(12345), snowflake.ID(12345)).Return(nil, common.ErrAppQuotaAlreadyExist).AnyTimes()
	s.mockAppQuotaSrv.EXPECT().AddQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()

	s.mockAppQuotaSrv.EXPECT().DeleteQuota(gomock.Any(), snowflake.ID(12345), snowflake.ID(54321)).Return(nil).AnyTimes()
	s.mockAppQuotaSrv.EXPECT().DeleteQuota(gomock.Any(), snowflake.ID(12345), snowflake.ID(12345)).Return(common.ErrAppQuotaNotFound).AnyTimes()
	s.mockAppQuotaSrv.EXPECT().DeleteQuota(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error")).AnyTimes()
}

func (s *ApplicationQuotaSuite) TearDownSuite() {
	s.T().Log("teardown suite")
}

func (s *ApplicationQuotaSuite) SetupSubTest() {
	s.T().Log("setup sub test")
	s.w = httptest.NewRecorder()
	s.ctx = mock.GinContext(s.w)
	s.ctx.Set(logging.LoggerName, s.logger)
	s.ctx.Request.Header.Set("x-ys-request-id", gofakeit.UUID())
}

func (s *ApplicationQuotaSuite) TearDownSubTest() {
	s.T().Log("teardown sub test")
}

func (s *ApplicationQuotaSuite) TestController_QuotaGet() {
	s.Run("success", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		// configure query params
		u := url.Values{}
		u.Add("UserID", "h9z") // 54321

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)

		s.controller.QuotaGet(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

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
		u.Add("UserID", "h9z") // 54321

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)

		s.controller.QuotaGet(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

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
		u.Add("UserID", "h9z") // 54321

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)

		s.controller.QuotaGet(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidAppID, resp.ErrorCode)
	})

	s.Run("empty UserID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		// configure query params
		u := url.Values{}
		u.Add("UserID", "")

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)

		s.controller.QuotaGet(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidUserID, resp.ErrorCode)
	})

	s.Run("error UserID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		// configure query params
		u := url.Values{}
		u.Add("UserID", "%#$%#$")

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)

		s.controller.QuotaGet(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidUserID, resp.ErrorCode)
	})

	s.Run("quota not found", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		// configure query params
		u := url.Values{}
		u.Add("UserID", "4ER") // 12345

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)

		s.controller.QuotaGet(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.AppQuotaNotFound, resp.ErrorCode)
	})

	s.Run("internal error", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "h9z"}, // 54321
		}
		// configure query params
		u := url.Values{}
		u.Add("UserID", "h9z") // 54321

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, params, u)

		s.controller.QuotaGet(s.ctx)

		s.Equal(http.StatusInternalServerError, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InternalServerErrorCode, resp.ErrorCode)
	})
}

func (s *ApplicationQuotaSuite) TestController_QuotaAdd() {
	s.Run("success", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		req := &quotaadd.Request{
			UserID: "h9z", // 54321
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, req, params, nil)

		s.controller.QuotaAdd(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)

		got := s.w.Body.String()
		s.T().Log(spew.Sdump(got))
	})

	s.Run("empty AppID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: ""},
		}
		req := &quotaadd.Request{
			UserID: "h9z", // 54321
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, req, params, nil)

		s.controller.QuotaAdd(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

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
		req := &quotaadd.Request{
			UserID: "h9z", // 54321
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, req, params, nil)

		s.controller.QuotaAdd(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

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
			{Key: "AppID", Value: "67W2"}, // 不存在的app
		}
		req := &quotaadd.Request{
			UserID: "h9z", // 54321
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, req, params, nil)

		s.controller.QuotaAdd(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.AppIDNotFoundErrorCode, resp.ErrorCode)
	})

	s.Run("empty UserID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		req := &quotaadd.Request{
			UserID: "",
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, req, params, nil)

		s.controller.QuotaAdd(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidUserID, resp.ErrorCode)
	})

	s.Run("error UserID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		req := &quotaadd.Request{
			UserID: "%#$%#$",
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, req, params, nil)

		s.controller.QuotaAdd(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidUserID, resp.ErrorCode)
	})

	s.Run("user not found", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		req := &quotaadd.Request{
			UserID: "67W2", // 不存在的user
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, req, params, nil)

		s.controller.QuotaAdd(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.UserNotExistsErrorCode, resp.ErrorCode)
	})

	s.Run("quota already exists", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		req := &quotaadd.Request{
			UserID: "4ER", // 12345
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, req, params, nil)

		s.controller.QuotaAdd(s.ctx)

		s.Equal(http.StatusConflict, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.AppQuotaAlreadyExist, resp.ErrorCode)
	})

	s.Run("internal error", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "h9z"}, // 54321
		}
		req := &quotaadd.Request{
			UserID: "h9z", // 54321
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, req, params, nil)

		s.controller.QuotaAdd(s.ctx)

		s.Equal(http.StatusInternalServerError, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotaadd.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InternalServerErrorCode, resp.ErrorCode)
	})
}

func (s *ApplicationQuotaSuite) TestController_QuotaDelete() {
	s.Run("success", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		req := &quotadelete.Request{
			UserID: "h9z", // 54321
		}

		mock.HTTPRequest(s.ctx, http.MethodDelete, req, params, nil)

		s.controller.QuotaDelete(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)

		got := s.w.Body.String()
		s.T().Log(spew.Sdump(got))
	})

	s.Run("empty AppID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: ""},
		}
		req := &quotadelete.Request{
			UserID: "h9z", // 54321
		}

		mock.HTTPRequest(s.ctx, http.MethodDelete, req, params, nil)

		s.controller.QuotaDelete(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotadelete.Response

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
		req := &quotadelete.Request{
			UserID: "h9z", // 54321
		}

		mock.HTTPRequest(s.ctx, http.MethodDelete, req, params, nil)

		s.controller.QuotaDelete(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotadelete.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidAppID, resp.ErrorCode)
	})

	s.Run("empty UserID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		req := &quotadelete.Request{
			UserID: "",
		}

		mock.HTTPRequest(s.ctx, http.MethodDelete, req, params, nil)

		s.controller.QuotaDelete(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotadelete.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidUserID, resp.ErrorCode)
	})

	s.Run("error UserID", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		req := &quotadelete.Request{
			UserID: "%#$%#$",
		}

		mock.HTTPRequest(s.ctx, http.MethodDelete, req, params, nil)

		s.controller.QuotaDelete(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotadelete.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InvalidUserID, resp.ErrorCode)
	})

	s.Run("quota not found", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "4ER"}, // 12345
		}
		req := &quotadelete.Request{
			UserID: "4ER", // 12345
		}

		mock.HTTPRequest(s.ctx, http.MethodDelete, req, params, nil)

		s.controller.QuotaDelete(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotadelete.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.AppQuotaNotFound, resp.ErrorCode)
	})

	s.Run("internal error", func() {
		// configure path params
		params := []gin.Param{
			{Key: "AppID", Value: "h9z"}, // 54321
		}
		req := &quotadelete.Request{
			UserID: "h9z", // 54321
		}

		mock.HTTPRequest(s.ctx, http.MethodDelete, req, params, nil)

		s.controller.QuotaDelete(s.ctx)

		s.Equal(http.StatusInternalServerError, s.w.Code)

		got := s.w.Body.Bytes()
		var resp quotadelete.Response

		// 解析到结构体
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.Equal(api.InternalServerErrorCode, resp.ErrorCode)
	})
}
