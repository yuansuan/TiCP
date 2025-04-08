package api

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/user/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/service/impl"
)

type RouteService struct {
	AuthService    service.AuthService
	UserService    service.UserService
	LicenseService service.LicenseService
	OrgService     service.OrgService
}

func NewUserService() (*RouteService, error) {
	authService := impl.NewAuthService()
	userService := impl.NewUserService()
	licenseService := impl.NewLicenseService()
	orgService := impl.NewOrgService()

	return &RouteService{
		AuthService:    authService,
		UserService:    userService,
		LicenseService: licenseService,
		OrgService:     orgService,
	}, nil
}

// InitAPI 初始化API服务
func InitAPI(drv *http.Driver) {
	logger := logging.Default()

	s, err := NewUserService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	authGroup := drv.Group("/api/v1/auth")
	{
		authGroup.POST("/login", s.Login)
		authGroup.POST("/logout", s.Logout)
		authGroup.GET("/ping/ldap", s.PingLdap)

		authGroup.POST("/onlineList", s.OnlineList)
		authGroup.POST("/onlineListByUser", s.OnlineListByUser)
		authGroup.POST("/offlineByUserName", s.OfflineByUserName)
		authGroup.POST("/offlineByJti", s.OfflineByJti)

		authGroup.GET("/license/machineID", s.GetMachineID)
		authGroup.GET("/license", s.GetLicense)
		authGroup.POST("/license", s.UpdateLicense)
	}

	userGroup := drv.Group("/api/v1/user")
	{
		userGroup.POST("/add", s.AddUser)
		userGroup.POST("/query", s.Query)
		userGroup.GET("/current", s.Current)
		userGroup.GET("/get", s.Get)
		userGroup.PUT("/active", s.Active)
		userGroup.PUT("/inactive", s.Inactive)
		userGroup.DELETE("/delete", s.Delete)
		userGroup.PUT("/updatePassword", s.UpdatePassword)
		userGroup.GET("/getDataConfig", s.GetDataConfig)
		userGroup.PUT("/update", s.Update)
		userGroup.GET("/optionList", s.OptionList)
		userGroup.POST("/resetPassword", s.ResetPassword)
		userGroup.PUT("/genOpenapiCertificate", s.GenOpenapiCertificate)
	}

	orgGroup := drv.Group("/api/v1/org")
	{
		// 创建组织
		orgGroup.POST("/create", s.CreateOrg)
		// 删除组织
		orgGroup.DELETE("/delete", s.DeleteOrg)
		// 修改组织
		orgGroup.PUT("/update", s.UpdateOrg)
		// 新增组织成员
		orgGroup.POST("/member/add", s.AddOrgMember)
		// 删除组织成员
		orgGroup.DELETE("/member/delete", s.DeleteOrgMember)
		// 修改组织成员
		orgGroup.PUT("/member/update", s.UpdateOrgMember)
		// 查询某组织架构下的成员
		orgGroup.GET("/member/list", s.ListOrgMember)
	}
}
