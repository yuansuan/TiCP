package openapi

import (
	"github.com/go-playground/validator/v10"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service/impl"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

type RouteOpenapiService struct {
	LocalFileService service.FileService
	validate         *validator.Validate
}

func NewRouteOpenapiService() (*RouteOpenapiService, error) {
	localFileService, err := impl.NewLocalFileService()
	if err != nil {
		return nil, err
	}

	return &RouteOpenapiService{
		LocalFileService: localFileService,
		validate:         validator.New(),
	}, nil
}

func InitOpenapiAPI(drv *http.Driver) {
	logger := logging.Default()

	s, err := NewRouteOpenapiService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	openapiGroup := drv.Group("/api/v1/openapi/storage")
	{
		openapiGroup.POST("/preUpload", s.PreUpload)
		openapiGroup.POST("/upload", s.Upload)
		openapiGroup.POST("/list", s.List)
		openapiGroup.POST("/remove", s.Remove)
		openapiGroup.GET("/batchDownload", s.BatchDownload)
		openapiGroup.POST("/batchDownloadPre", s.BatchDownloadPre)

	}

}
