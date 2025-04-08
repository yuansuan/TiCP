package test

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/shared_directory/api"
	"xorm.io/xorm"
)

// GetTestGinContext returns a gin context for testing
func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

func BindRequest(ctx *gin.Context, req interface{}) {

	// 获取原始查询参数
	originalQuery := ctx.Request.URL.Query()

	// 使用反射获取结构体字段的名称和值，并将其设置为查询参数
	valueOfReq := reflect.ValueOf(req)
	typeOfReq := reflect.TypeOf(req)

	for i := 0; i < typeOfReq.NumField(); i++ {
		field := typeOfReq.Field(i)
		fieldValue := valueOfReq.Field(i)

		// 跳过未导出字段或空值字段
		if field.PkgPath != "" || fieldValue.IsZero() {
			continue
		}

		// 将字段名称和值转换为字符串，并设置为查询参数
		paramName := field.Name
		paramValue := fmt.Sprintf("%v", fieldValue.Interface())
		originalQuery.Set(paramName, paramValue)
	}

	// 生成新的 URL 对象并替换上下文中的 URL
	ctx.Request.URL.RawQuery = originalQuery.Encode()
	ctx.Request.Body = http.NoBody
}

func BindJsonRequest(ctx *gin.Context, req interface{}) error {
	bodyData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyData))
	ctx.Request.ContentLength = int64(len(bodyData))
	ctx.Request.Header.Set("Content-Type", "application/json")
	return nil

}

func CreateLocalStorage(userID, linkType string, storageQuotaDao dao.StorageQuotaDao, storageOperationLogDao dao.StorageOperationLogDao, engine *xorm.Engine) (string, *v20230530.Storage) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("get current dir failed: %v", err))
	}

	authEnabled := false
	cfg := &config.CustomT{
		Local: config.Local{
			RootPath: currentDir,
			LinkType: linkType,
		},
		YsId:        userID,
		AuthEnabled: &authEnabled,
	}
	config.Custom = *cfg
	storage, err := v20230530.NewStorage(cfg.Local.RootPath, engine, storageQuotaDao, storageOperationLogDao, nil, nil, nil)

	if err != nil {
		panic(fmt.Sprintf("init storage failed: %v", err))
	}
	return currentDir, storage
}

func CreateLocalStorageWithUploadDao(userID, linkType string, uploadInfo dao.UploadInfoDao, storageQuotaDao dao.StorageQuotaDao, storageOperationLogDao dao.StorageOperationLogDao, engine *xorm.Engine) (string, *v20230530.Storage) {
	currentDir, storage := CreateLocalStorage(userID, linkType, storageQuotaDao, storageOperationLogDao, engine)
	storage.UploadInfoDao = uploadInfo
	storage.Engine = engine
	return currentDir, storage
}

func CreateLocalStorageWithSharedDirectory(userID, linkType string, sharedDirectoryDao dao.StorageSharedDirectoryDao, storageQuotaDao dao.StorageQuotaDao, storageOperationLogDao dao.StorageOperationLogDao, engine *xorm.Engine) (string, *v20230530.Storage) {
	currentDir, storage := CreateLocalStorage(userID, linkType, storageQuotaDao, storageOperationLogDao, engine)
	storage.SharedDirectory = api.NewSharedDirectory(sharedDirectoryDao, engine, api.GetHC(), storage.PathAccessCheckerImpl)
	storage.Engine = engine
	return currentDir, storage
}

func GenerateRandomData(fileSize int64) (error, io.Reader) {
	data := make([]byte, fileSize*1024)
	_, err := rand.Read(data)
	if err != nil {
		return err, nil
	}
	randomReader := bytes.NewReader(data)

	return nil, randomReader
}
