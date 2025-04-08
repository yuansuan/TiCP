package api

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/service/impl"
)

type RouteService struct {
	RoleService service.RoleService
	PermService service.PermService
}

func NewRbacService() (*RouteService, error) {
	roleService, err := impl.NewRoleService()
	permService, err := impl.NewPermService()
	if err != nil {
		return nil, err
	}

	return &RouteService{
		RoleService: roleService,
		PermService: permService,
	}, nil
}

// InitAPI 初始化API服务
func InitAPI(drv *http.Driver) {
	logger := logging.Default()

	s, err := NewRbacService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	roleGroup := drv.Group("/api/v1/role")
	{
		roleGroup.POST("/add", s.AddRole)
		roleGroup.POST("/query", s.QueryRole)
		roleGroup.GET("/detail", s.GetRoleDetail)
		roleGroup.PUT("/update", s.UpdateRole)
		roleGroup.DELETE("/delete", s.DeleteRole)
		roleGroup.PUT("/setLdapUserDefRole", s.SetLdapUserDefRole)

	}

}
