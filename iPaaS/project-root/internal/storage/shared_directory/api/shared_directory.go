package api

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/pathchecker"
	"xorm.io/xorm"
)

var _hc *resty.Client
var once sync.Once

// GetHC 获取http client
func GetHC() *resty.Client {
	once.Do(func() {
		_hc = resty.New().SetTimeout(5 * time.Second)
	})
	return _hc
}

// SharedDirectory 共享目录api
type SharedDirectory struct {
	Engine                    *xorm.Engine
	StorageSharedDirectoryDao dao.StorageSharedDirectoryDao

	hc *resty.Client

	pathchecker.PathAccessCheckerImpl
}

// NewSharedDirectory 新建共享目录api
func NewSharedDirectory(storageSharedDirectoryDao dao.StorageSharedDirectoryDao, engine *xorm.Engine, hc *resty.Client, pathchecker pathchecker.PathAccessCheckerImpl) *SharedDirectory {
	if storageSharedDirectoryDao == nil {
		return nil
	}

	return &SharedDirectory{
		Engine:                    engine,
		StorageSharedDirectoryDao: storageSharedDirectoryDao,

		hc: hc,

		PathAccessCheckerImpl: pathchecker,
	}
}

// GetUserIDAndAKAndHandleError 获取用户id和ak
func (s *SharedDirectory) GetUserIDAndAKAndHandleError(ctx *gin.Context) (string, string, bool, error) {
	return s.PathAccessCheckerImpl.GetUserIDAndAKAndHandleError(ctx, pathchecker.SystemURLPrefix)
}
