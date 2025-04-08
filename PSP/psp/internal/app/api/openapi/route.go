package openapi

import (
	"github.com/go-playground/validator/v10"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/service/impl"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

type RouteOpenapiService struct {
	appService service.AppService
	validate   *validator.Validate
}

func NewRouteOpenapiService() (*RouteOpenapiService, error) {
	logger := logging.Default()
	appService, err := impl.NewAppService()
	if err != nil {
		logger.Errorf("init app server service err: %v", err)
		return nil, err
	}

	return &RouteOpenapiService{
		appService: appService,
		validate:   validator.New(),
	}, nil
}

func InitOpenapiAPI(drv *http.Driver) {
	logger := logging.Default()

	s, err := NewRouteOpenapiService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	openapiAppGroup := drv.Group("/api/v1/openapi/app")
	{
		openapiAppGroup.GET("/list", s.ListApp)
	}

}
