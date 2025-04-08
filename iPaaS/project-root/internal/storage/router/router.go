package router

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging/middleware"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	iamclient "github.com/yuansuan/ticp/common/project-root-iam/iam-client"
	trafficStat "github.com/yuansuan/ticp/iPaaS/project-root/internal/common/middleware"
	directoryUsageApi "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/directory_usage/api"
	operationLogApi "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/operationlog/api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/quota/api"
	SharedDirectoryApi "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/shared_directory/api"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/config"
)

const (
	VersionHeader     = "X-Ys-Version"
	ApiVersion        = "api-version"
	requestIDKey      = "x-ys-request-id"
	requestIDKeyInCtx = "RequestId"
)

// Init 路由
func Init(drv *http.Driver) {
	storage.New()

	apiGroup := drv.Group("/api", VersionHeaderMiddleware)
	setupIamMiddleware(apiGroup)
	apiGroup.Use(trafficStat.HttpStatMiddleware())
	apiGroup.Use(ingressLoggerMiddleware)
	storageGroup := apiGroup.Group("/storage")
	{
		storageGroup.GET("/healthz", health)
		storageGroup.GET("/lsWithPage", ls)
		storageGroup.POST("/mkdir", mkdir)
		storageGroup.GET("/download", download)
		storageGroup.GET("/batchDownload", batchDownload)
		storageGroup.POST("/batchDownload", batchDownload)
		storageGroup.POST("/rm", rm)
		storageGroup.POST("/mv", mv)
		storageGroup.GET("/stat", stat)
		storageGroup.POST("/upload/init", uploadInit)
		storageGroup.POST("/upload/slice", uploadSlice)
		storageGroup.POST("/upload/complete", uploadComplete)
		storageGroup.POST("/upload/file", uploadFile)
		storageGroup.POST("/link", link)
		storageGroup.POST("/copy", copy2)
		storageGroup.POST("/copyRange", copyRange)
		storageGroup.POST("/writeAt", writeAt)
		storageGroup.GET("/readAt", readAt)
		storageGroup.POST("/create", create)
		storageGroup.POST("/truncate", truncate)
		storageGroup.POST("/compress/start", compressStart)
		storageGroup.GET("/compress/status", compressStatus)
		storageGroup.POST("/compress/cancel", compressCancel)
		storageGroup.GET("/storageQuota", storageQuotaAPI)
		storageGroup.GET("/operationLog", listOperationLogApi)
		storageGroup.POST("/sharedDirectorys", CheckSharedDirectoryConfigMiddleware, createSharedDirectory)
		storageGroup.DELETE("/sharedDirectorys", CheckSharedDirectoryConfigMiddleware, deleteSharedDirectory)
		storageGroup.GET("/sharedDirectorys", listSharedDirectorys)
		storageGroup.POST("/directoryUsage/start", directoryUsageStart)
		storageGroup.GET("/directoryUsage/status", directoryUsageStatus)
		storageGroup.POST("/directoryUsage/cancel", directoryUsageCancel)
	}

	sysGroup := drv.Group("/system", VersionHeaderMiddleware)
	setupIamMiddleware(sysGroup)
	sysGroup.Use(ingressLoggerMiddleware)

	sysStorageGroup := sysGroup.Group("/storage")
	{
		sysStorageGroup.GET("/realpath", realpath)
		sysStorageGroup.GET("/lsWithPage", ls)
		sysStorageGroup.POST("/mkdir", mkdir)
		sysStorageGroup.GET("/download", download)
		sysStorageGroup.GET("/batchDownload", batchDownload)
		sysStorageGroup.POST("/batchDownload", batchDownload)
		sysStorageGroup.POST("/rm", rm)
		sysStorageGroup.POST("/mv", mv)
		sysStorageGroup.GET("/stat", stat)
		sysStorageGroup.POST("/upload/init", uploadInit)
		sysStorageGroup.POST("/upload/slice", uploadSlice)
		sysStorageGroup.POST("/upload/complete", uploadComplete)
		sysStorageGroup.POST("/upload/file", uploadFile)
		sysStorageGroup.POST("/link", link)
		sysStorageGroup.POST("/copy", copy2)
		sysStorageGroup.POST("/copyRange", copyRange)
		sysStorageGroup.POST("/writeAt", writeAt)
		sysStorageGroup.GET("/readAt", readAt)
		sysStorageGroup.POST("/create", create)
		sysStorageGroup.POST("/truncate", truncate)
		sysStorageGroup.POST("/compress/start", compressStart)
		sysStorageGroup.GET("/compress/status", compressStatus)
		sysStorageGroup.POST("/compress/cancel", compressCancel)
		sysStorageGroup.POST("/directoryUsage/start", directoryUsageStart)
		sysStorageGroup.GET("/directoryUsage/status", directoryUsageStatus)
		sysStorageGroup.POST("/directoryUsage/cancel", directoryUsageCancel)
	}

	adminGroup := drv.Group("/admin", VersionHeaderMiddleware)
	setupIamMiddleware(adminGroup)
	adminGroup.Use(ingressLoggerMiddleware)
	adminStorageGroup := adminGroup.Group("/storage")
	{
		adminStorageGroup.GET("/:UserID/storageQuota", storageQuotaAdmin)
		adminStorageGroup.GET("/storageQuota/total", storageQuotaTotal)
		adminStorageGroup.GET("/storageQuota", listStorageQuota)
		adminStorageGroup.PUT("/:UserID/storageQuota", putStorageQuota)
		adminStorageGroup.GET("/operationLog", listOperationLogAdmin)
	}

}

func CheckSharedDirectoryConfigMiddleware(c *gin.Context) {
	if config.GetConfig().ShareRegisterAddress == "" || config.GetConfig().SharedHost == "" {
		common.ErrorResp(c, nethttp.StatusForbidden, commoncode.AccessDeniedErrorCode, "not support this api")
		return
	}

	c.Next()
}

func setupIamMiddleware(group *gin.RouterGroup) {
	// 启动vip-box，关闭校验
	if *config.GetConfig().AuthEnabled {
		iamConfig := iamclient.IamConfig{
			Endpoint:  config.GetConfig().IamServerUrl,
			AppKey:    config.GetConfig().AccessKeyId,
			AppSecret: config.GetConfig().AccessKeySecret,
		}
		group.Use(iamclient.SignatureValidateMiddleware(iamConfig))
	} else {
		group.Use(requestIDMiddleware)
	}
}

func health(ctx *gin.Context) {
	common.SuccessResp(ctx, "")
}

func ls(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.Ls(ctx)
	})
}

func mkdir(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.Mkdir(ctx)
	})
}

func download(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.Download(ctx)
	})
}

func batchDownload(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.BatchDownload(ctx)
	})
}

func rm(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.Rm(ctx)
	})
}

func mv(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.Mv(ctx)
	})
}

func stat(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.Stat(ctx)
	})
}

func uploadInit(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.UploadInit(ctx)
	})
}

func uploadSlice(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.UploadSlice(ctx)
	})
}

func uploadComplete(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.UploadComplete(ctx)
	})
}

func realpath(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.Realpath(ctx)
	})
}

func copy2(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.Copy(ctx)
	})
}

func copyRange(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.CopyRange(ctx)
	})
}

func create(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.Create(ctx)
	})
}

func writeAt(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.WriteAt(ctx)
	})
}

func truncate(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.Truncate(ctx)
	})
}

func readAt(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.ReadAt(ctx)
	})
}

func uploadFile(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.UploadFile(ctx)
	})
}

func link(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.Link(ctx)
	})
}

func compressStart(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.CompressStart(ctx)
	})
}

func compressStatus(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.CompressStatus(ctx)
	})
}

func compressCancel(ctx *gin.Context) {
	doAction(ctx, func(storage storage.Storage) {
		storage.CompressCancel(ctx)
	})
}

func directoryUsageStart(ctx *gin.Context) {
	directoryUsage, err := getDirectoryUsage(ctx)
	if err != nil {
		logging.Default().Infof("get directory usage error, error:%v", err)
		return
	}
	directoryUsage.DirectoryUsageStart(ctx)
}

func directoryUsageStatus(ctx *gin.Context) {
	directoryUsage, err := getDirectoryUsage(ctx)
	if err != nil {
		logging.Default().Infof("get directory usage error, error:%v", err)
		return
	}
	directoryUsage.DirectoryUsageStatus(ctx)
}

func directoryUsageCancel(ctx *gin.Context) {
	directoryUsage, err := getDirectoryUsage(ctx)
	if err != nil {
		logging.Default().Infof("get directory usage error, error:%v", err)
		return
	}
	directoryUsage.DirectoryUsageCancel(ctx)
}

func storageQuotaAPI(ctx *gin.Context) {
	quota, _, err := getStorageQuotaAndOperationLog(ctx)
	if err != nil {
		logging.Default().Infof("get storage quota error, error:%v", err)
		return
	}
	quota.GetStorageQuotaAPI(ctx)
}

func storageQuotaAdmin(ctx *gin.Context) {
	quota, _, err := getStorageQuotaAndOperationLog(ctx)
	if err != nil {
		logging.Default().Infof("get storage quota error, error:%v", err)
		return
	}
	quota.GetStorageQuotaAdmin(ctx)
}

func storageQuotaTotal(ctx *gin.Context) {
	quota, _, err := getStorageQuotaAndOperationLog(ctx)
	if err != nil {
		logging.Default().Infof("get storage quota error, error:%v", err)
		return
	}
	quota.GetStorageQuotaTotal(ctx)
}

func listStorageQuota(ctx *gin.Context) {
	quota, _, err := getStorageQuotaAndOperationLog(ctx)
	if err != nil {
		logging.Default().Infof("get storage quota error, error:%v", err)
		return
	}
	quota.ListStorageQuota(ctx)
}

func putStorageQuota(ctx *gin.Context) {
	quota, _, err := getStorageQuotaAndOperationLog(ctx)
	if err != nil {
		logging.Default().Infof("get storage quota error, error:%v", err)
		return
	}
	quota.PutStorageQuota(ctx)
}

func listOperationLogApi(ctx *gin.Context) {
	_, operationLog, err := getStorageQuotaAndOperationLog(ctx)
	if err != nil {
		logging.Default().Infof("get storage quota error, error:%v", err)
		return
	}
	operationLog.ListOperationLogApi(ctx)
}

func listOperationLogAdmin(ctx *gin.Context) {
	_, operationLog, err := getStorageQuotaAndOperationLog(ctx)
	if err != nil {
		logging.Default().Infof("get storage quota error, error:%v", err)
		return
	}
	operationLog.ListOperationLogAdmin(ctx)
}

func createSharedDirectory(ctx *gin.Context) {
	sharedDirectory, err := getSharedDirectoryAPI(ctx)
	if err != nil {
		logging.Default().Infof("get sharedDirectory api error, error:%v", err)
		return
	}
	sharedDirectory.Create(ctx)
}

func deleteSharedDirectory(ctx *gin.Context) {
	sharedDirectory, err := getSharedDirectoryAPI(ctx)
	if err != nil {
		logging.Default().Infof("get sharedDirectory api error, error:%v", err)
		return
	}
	sharedDirectory.Delete(ctx)
}

func listSharedDirectorys(ctx *gin.Context) {
	sharedDirectory, err := getSharedDirectoryAPI(ctx)
	if err != nil {
		logging.Default().Infof("get sharedDirectory api error, error:%v", err)
		return
	}
	sharedDirectory.List(ctx)
}

func getStorageQuotaAndOperationLog(ctx *gin.Context) (*api.Quota, *operationLogApi.OperationLog, error) {
	version := ctx.GetString(ApiVersion)
	_, quota, operationLog, _, _, err := storage.GetStorage(version)
	if err != nil {
		common.ErrorResp(ctx, nethttp.StatusBadRequest, commoncode.InvalidStorageVersion, err.Error())
		return nil, nil, err
	}
	return quota, operationLog, nil
}

func getDirectoryUsage(ctx *gin.Context) (*directoryUsageApi.DirectoryUsage, error) {
	version := ctx.GetString(ApiVersion)
	_, _, _, _, directoryUsage, err := storage.GetStorage(version)
	if err != nil {
		common.ErrorResp(ctx, nethttp.StatusBadRequest, commoncode.InvalidStorageVersion, err.Error())
		return nil, err
	}
	return directoryUsage, nil
}

func getSharedDirectoryAPI(ctx *gin.Context) (*SharedDirectoryApi.SharedDirectory, error) {
	version := ctx.GetString(ApiVersion)
	_, _, _, sharedDirectory, _, err := storage.GetStorage(version)
	if err != nil {
		common.ErrorResp(ctx, nethttp.StatusBadRequest, commoncode.InvalidStorageVersion, err.Error())
		return nil, err
	}
	return sharedDirectory, nil
}

func doAction(ctx *gin.Context, do func(storage storage.Storage)) {
	version := ctx.GetString(ApiVersion)
	storageImpl, _, _, _, _, err := storage.GetStorage(version)
	if err != nil {
		common.ErrorResp(ctx, nethttp.StatusBadRequest, commoncode.InvalidStorageVersion, err.Error())
		return
	}

	do(storageImpl)
}

func VersionHeaderMiddleware(ctx *gin.Context) {
	version := ctx.GetHeader(VersionHeader)
	// 如果header里面没有，就从url里面取
	if version == "" {
		version = ctx.Query(VersionHeader)
	}

	ctx.Set(ApiVersion, version)
	ctx.Next()
}

var logBodyWhiteList = []string{
	"/api/storage/lsWithPage",
	"/api/storage/mkdir",
	"/api/storage/rm",
	"/api/storage/mv",
	"/api/storage/stat",
	"/api/storage/upload/init",
	"/api/storage/upload/complete",
	"/api/storage/realpath",
	"/api/storage/link",
	"/api/storage/checksum",
	"/api/storage/copy",
	"/api/storage/copyRange",
	"/api/storage/create",
	"/api/storage/truncate",
	"/api/storage/compress/start",
	"/api/storage/compress/status",
	"/system/storage/realpath",
}

func inLogBodyWhiteList(path string) bool {
	for _, v := range logBodyWhiteList {
		if path == v {
			return true
		}
	}

	return false
}

func ingressLoggerMiddleware(c *gin.Context) {
	if inLogBodyWhiteList(c.Request.URL.Path) {
		middleware.IngressLogger(middleware.IngressLoggerConfig{
			IsLogRequestHeader:  true,
			IsLogRequestBody:    true,
			IsLogResponseHeader: true,
			IsLogResponseBody:   true,
		})(c)
	} else {
		middleware.IngressLogger(middleware.IngressLoggerConfig{
			IsLogRequestHeader:  true,
			IsLogRequestBody:    false,
			IsLogResponseHeader: true,
			IsLogResponseBody:   false,
		})(c)
	}

}

func requestIDMiddleware(c *gin.Context) {
	requestId := uuid.New().String()
	c.Set(requestIDKey, requestId)
	c.Set(requestIDKeyInCtx, requestId)
	c.Request.Header.Set(requestIDKey, requestId)
	c.Writer.Header().Set(requestIDKey, requestId)
}
