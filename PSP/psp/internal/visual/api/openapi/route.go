package openapi

import (
	"github.com/go-playground/validator/v10"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/service/impl"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

type RouteOpenapiService struct {
	visualService service.VisualService
	validate      *validator.Validate
}

func NewRouteOpenapiService() (*RouteOpenapiService, error) {
	logger := logging.Default()
	visualService, err := impl.NewVisualService()
	if err != nil {
		logger.Errorf("init visual server service err: %v", err)
		return nil, err
	}
	return &RouteOpenapiService{
		visualService: visualService,
		validate:      validator.New(),
	}, nil
}

func InitOpenapiAPI(drv *http.Driver) {
	logger := logging.Default()

	s, err := NewRouteOpenapiService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	openapiAppGroup := drv.Group("/api/v1/openapi/vis")
	{
		openapiAppGroup.GET("/software", s.ListSoftware)
		openapiAppGroup.GET("/hardware", s.ListHardware)
		openapiAppGroup.POST("/session/close", s.CloseSession)
		openapiAppGroup.POST("/session", s.StartSession)
		openapiAppGroup.GET("/session/list", s.ListSession)
		openapiAppGroup.GET("/session", s.SessionInfo)
		openapiAppGroup.POST("/session/reboot", s.RebootSession)

	}

}
