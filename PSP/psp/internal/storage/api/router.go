package api

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service/impl"
)

type RouteService struct {
	LocalFileService service.FileService
}

func NewStorageService() (*RouteService, error) {
	localFileService, err := impl.NewLocalFileService()
	if err != nil {
		return nil, err
	}

	return &RouteService{
		LocalFileService: localFileService,
	}, nil
}

// InitAPI 初始化API服务
func InitAPI(drv *http.Driver) {
	logger := logging.Default()

	s, err := NewStorageService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	group := drv.Group("/api/v1/storage")
	{
		group.POST("/preUpload", s.PreUpload)
		group.POST("/upload", s.Upload)
		group.GET("/batchDownload", s.BatchDownload)
		group.POST("/batchDownloadPre", s.BatchDownloadPre)
		group.PUT("/copy", s.Copy)
		group.GET("/copyStatus", s.CopyStatus)
		group.POST("/list", s.List)
		group.POST("/listOfRecur", s.ListOfRecur)
		group.POST("/get", s.Get)
		group.PUT("/rename", s.Rename)
		group.PUT("/move", s.Move)
		group.POST("/createDir", s.CreateDir)
		group.POST("/remove", s.Remove)
		group.POST("/read", s.Read)
		group.POST("/link", s.Link)
		group.POST("/compress", s.Compress)
		group.GET("/compressTasks", s.CompressTasks)
		group.POST("/share/send", s.GenerateAndSendShareCode)
		group.POST("/share/recordList", s.GetRecordList)
		group.GET("/share/get", s.GetShareFile)
		group.GET("/share/count", s.ShareCount)
		group.PUT("/share/updateRecordState", s.UpdateRecordState)
		group.PUT("/share/readAll", s.ReadAll)

	}

}
