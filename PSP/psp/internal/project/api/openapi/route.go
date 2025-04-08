package openapi

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/project/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/service/impl"
)

type RouteOpenapiService struct {
	projectService       service.ProjectService
	projectMemberService service.ProjectMemberService
}

func NewRouteOpenapiService() (*RouteOpenapiService, error) {
	projectService, err := impl.NewProjectService()
	if err != nil {
		return nil, err
	}

	memberService, err := impl.NewProjectMemberService()
	return &RouteOpenapiService{
		projectService:       projectService,
		projectMemberService: memberService,
	}, nil
}

func InitOpenapiAPI(drv *http.Driver) {
	logger := logging.Default()

	api, err := NewRouteOpenapiService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	openapiGroup := drv.Group("/api/v1/openapi/project")
	{
		openapiGroup.POST("/list", api.ProjectList)

	}

}
