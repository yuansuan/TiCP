package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type CreateSuite struct {
	suite.Suite
	ctx            *gin.Context
	httpRecord     *httptest.ResponseRecorder
	engine         *xorm.Engine
	sqlMock        sqlmock.Sqlmock
	licenseHandler *LicenseHandler
}

func TestLicense(t *testing.T) {
	suite.Run(t, new(CreateSuite))
}

func (s *CreateSuite) SetupSuite() {
	s.T().Log("setup suite")

	recorder := httptest.NewRecorder()
	ctx := mock.GinContext(recorder)
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}
	ctx.Set(logging.LoggerName, logger)

	mockEngine, sqlMock := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(logger, true))

	s.httpRecord = recorder
	s.engine = mockEngine
	s.sqlMock = sqlMock
	s.ctx = ctx

	licenseImpl := dao.NewLicenseImpl(s.engine)
	s.licenseHandler = NewLicenseHandler(licenseImpl)
}

func (s *CreateSuite) TearDownSuite() {
	s.T().Log("teardown suite")
}

func (s *CreateSuite) SetupTest() {
	s.T().Log("setup test")
}

func (s *CreateSuite) TearDownTest() {
	s.T().Log("teardown test")
}

func (s *CreateSuite) TestListLicenseManager() {
	s.T().Log("list license manage start")

	queryRows := s.sqlMock.NewRows([]string{
		"id",
		"app_type",
		"os",
		"status",
		"description",
		"compute_rule",
		"publish_time",
		"create_time",
		"update_time",
		"id",
		"manager_id",
		"provider",
		"license_server",
		"mac_addr",
		"license_url",
		"license_port",
		"license_proxies",
		"license_num",
		"weight",
		"begin_time",
		"end_time",
		"auth",
		"license_type",
		"tool_path",
		"collector_type",
		"hpc_endpoint",
		"allowable_hpc_endpoints",
		"create_time",
		"update_time",
		"id",
		"license_id",
		"module_name",
		"total",
		"used",
		"actual_total",
		"actual_used",
		"create_time",
		"update_time",
	}).AddRow("1680614605799297024",
		"starccm++",
		"2",
		"1",
		"aaa",
		"echo'{\"ccmppower\":1}'",
		"2023-07-31 20:00:56",
		"2023-07-17 00:25:20",
		"2023-07-31 20:00:56",
		"1680811786686697472",
		"1680614605799297024",
		"爱国者1111",
		"CDLMD_LICENSE_FILE",
		"00:00:00:00:00:00",
		"http://www.baidu.com",
		"39000",
		"{\"http://10.0.10.3:8080\":{\"Url\":\"1.1.1.1\",\"Port\":29000}}",
		"4E30AFEFCB5DD27D353F5B3D7D",
		"100",
		"2021-11-06 16:05:35",
		"2025-11-07 16:05:35",
		"1",
		"1",
		"/es01/home/zssoftware/Ansys/AnsysEM2021r1/shared_files/licensing/linx64/lmutil",
		"flex",
		"https://root-jn-hpc-test.yuansuan.com:31322",
		"[\"http://10.0.10.3:8080\"]",
		"2023-07-17 13:28:52",
		"2023-08-01 13:18:56",
		"1680849806248910848",
		"1680811786686697472",
		"ccmppower8",
		"22",
		"0",
		"0",
		"0",
		"2023-07-17 15:59:57",
		"2023-08-01 14:10:52")

	expectedQuery := regexp.QuoteMeta("SELECT * FROM `license_manager` LEFT JOIN `license_info` ON license_manager.id = license_info.manager_id LEFT JOIN `module_config` ON license_info.id = module_config.license_id")
	s.sqlMock.ExpectQuery(expectedQuery).WillReturnRows(queryRows)

	countQueryRows := s.sqlMock.NewRows([]string{"count"}).AddRow(1)
	exceptedCountQuery := regexp.QuoteMeta("SELECT count(*) FROM `license_manager` LEFT JOIN `license_info` ON license_manager.id = license_info.manager_id LEFT JOIN `module_config` ON license_info.id = module_config.license_id")
	s.sqlMock.ExpectQuery(exceptedCountQuery).WillReturnRows(countQueryRows)
	s.licenseHandler.ListLicenseManage(s.ctx)

	res := &licmanager.ListLicManagerResponse{}
	if err := json.Unmarshal(s.httpRecord.Body.Bytes(), res); err != nil {
		s.T().Error(err)
	}

	if res.ErrorCode != "" {
		s.T().Errorf("error code is %s, error message is %s", res.ErrorCode, res.ErrorMsg)
	}

	// 验证mock的期望行为是否被满足
	if err := s.sqlMock.ExpectationsWereMet(); err != nil {
		s.T().Error(err)
	}
}

func (s *CreateSuite) TestGetLicenseManage() {
	s.T().Log("list license manage start")

}
