package api

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/project/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/service/impl"
)

type apiRoute struct {
	projectService       service.ProjectService
	projectMemberService service.ProjectMemberService
}

func NewAPIRoute() (*apiRoute, error) {
	projectService, err := impl.NewProjectService()
	if err != nil {
		return nil, err
	}

	memberService, err := impl.NewProjectMemberService()
	return &apiRoute{
		projectService:       projectService,
		projectMemberService: memberService,
	}, nil
}

// InitAPI 初始化API服务
func InitAPI(drv *http.Driver) {
	logger := logging.Default()

	api, err := NewAPIRoute()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	group := drv.Group("/api/v1")

	{
		projectGroup := group.Group("/project")
		projectGroup.POST("/save", api.ProjectSave)
		projectGroup.POST("/list", api.ProjectList)
		projectGroup.GET("/list/current", api.CurrentProjectList)
		projectGroup.GET("/listForParam", api.CurrentProjectListForParam)
		projectGroup.GET("/detail", api.ProjectDetail)
		projectGroup.POST("/delete", api.ProjectDelete)
		projectGroup.POST("/terminate", api.ProjectTerminate)
		projectGroup.POST("/edit", api.ProjectEdit)
		projectGroup.POST("/modifyOwner", api.ProjectModifyOwner)
	}
	{
		projectMemberGroup := group.Group("/projectMember")
		projectMemberGroup.POST("/save", api.BatchUpdateProjectMember)
	}

}
