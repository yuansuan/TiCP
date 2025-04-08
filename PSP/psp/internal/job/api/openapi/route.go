package openapi

import (
	"github.com/go-playground/validator/v10"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/service/impl"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

type openapiApiRoute struct {
	jobService service.JobService
	validate   *validator.Validate
}

func NewRouteOpenapiService() (*openapiApiRoute, error) {
	jobService, err := impl.NewJobService()
	if err != nil {
		return nil, err
	}

	return &openapiApiRoute{
		jobService: jobService,
		validate:   validator.New(),
	}, nil
}

func InitOpenapiAPI(drv *http.Driver) {
	logger := logging.Default()

	api, err := NewRouteOpenapiService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	openapiGroup := drv.Group("/api/v1/openapi/job")
	{
		openapiGroup.POST("/createTempDir", api.CreateTempDir)
		openapiGroup.POST("/submit", api.Submit)
		openapiGroup.GET("/detail", api.JobDetail)
		openapiGroup.POST("/terminate", api.JobTerminate)

	}

}
