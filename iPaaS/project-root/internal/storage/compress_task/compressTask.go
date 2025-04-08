package compress_task

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	compressInfoService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/compressInfo"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"xorm.io/xorm"
)

var (
	CompressCache = cache.New(10*time.Hour, 24*time.Hour)
)

var (
	Semaphores = new(sync.Map)
)

const (
	ZipCacheExpireTime = 12 * time.Hour
	ZipCachePrefix     = "compress_zip"
)

type ZipCacheContent struct {
	IsFinished bool               `json:"isFinished"`
	Status     int                `json:"status"`
	TargetPath string             `json:"targetPath"`
	Cancel     context.CancelFunc `json:"cancel"`
	Err        string             `json:"err"`
}

type ZipCache struct {
	ContentMap *sync.Map
}

type CompressTask struct {
	Engine          *xorm.Engine
	CompressInfoDao dao.CompressInfoDao
	rootPath        string
}

func NewCompressTask(compressInfoDao dao.CompressInfoDao, engine *xorm.Engine, rootPath string) *CompressTask {
	if compressInfoDao == nil {
		return nil
	}

	return &CompressTask{
		Engine:          engine,
		CompressInfoDao: compressInfoDao,
		rootPath:        rootPath,
	}
}

func (t *CompressTask) Recover() {
	ctx := context.Background()
	logger := logging.GetLogger(ctx).With("func", "CompressTask.Recover")
	compressInfoList, err := compressInfoService.ListCompressInfo(ctx, t.Engine, t.CompressInfoDao)
	if err != nil {
		logger.Warnf("get unfinished compress task from database error:%v", err)
		return
	}
	logger.Infof("unfinished compress task number:%v", len(compressInfoList))
	for _, info := range compressInfoList {
		//删除临时文件
		tmpPath := filepath.Join(t.rootPath, fsutil.TrimPrefix(info.TmpPath, "/"))
		err := os.Remove(tmpPath)
		if err != nil {
			logger.Warnf("remove tmp path error:%v,path:%v", err, tmpPath)
		}
		logger.Infof("remove tmp path success,path:%v", tmpPath)
		//获得需要压缩文件路径
		pathsMap := make(map[string]struct{})
		emptyDirs := make([]string, 0)
		if t.GetFilePaths(&gin.Context{}, strings.Split(info.Paths, ","), pathsMap, emptyDirs, logger, false) {
			logger.Warnf("get file paths error,compressID:%v", info.Id)
			continue
		}

		var paths []string
		for path := range pathsMap {
			paths = append(paths, path)
		}
		//开始压缩
		compressCacheKey := fmt.Sprintf("%s_%s", ZipCachePrefix, info.UserId)
		zipCache := t.GetOrCreateZipCache(info.UserId, compressCacheKey)
		logger.Infof("recover compress task,compressID:%v", info.Id)
		go t.StartCompressTask(context.Background(), &CompressParams{
			TargetPath:          filepath.Join(t.rootPath, fsutil.TrimPrefix(info.TargetPath, "/")),
			TmpCompressFilePath: tmpPath,
			BasePath:            filepath.Join(t.rootPath, fsutil.TrimPrefix(info.BasePath, "/")),
			Paths:               paths,
			EmptyDirs:           emptyDirs,
			CompressCacheKey:    compressCacheKey,
			CompressID:          info.Id,
			Logger:              logger,
			ZipCache:            zipCache,
			CompressInfo:        info,
		})
	}
}

func (t *CompressTask) InsertCompressInfo(logger *logging.Logger, ctx *gin.Context, data *model.CompressInfo) error {
	err := compressInfoService.InsertCompressInfo(ctx, t.Engine, t.CompressInfoDao, data)
	if err != nil {
		logger.Warnf("insert compress info to database error:%v,compressID:%v", err, data.Id)
		return err
	}
	return nil
}

func (t *CompressTask) GetCompressInfo(logger *logging.Logger, ctx *gin.Context, compressID string) (bool, *model.CompressInfo, error) {
	exist, compressInfo, err := compressInfoService.GetCompressInfo(ctx, t.Engine, t.CompressInfoDao, compressID)
	if err != nil {
		logger.Warnf("get compress info from database error:%v,compressID:%v", err, compressID)
		return false, nil, err
	}
	if !exist {
		logger.Warnf("compress info not exist,compressID:%v", compressID)
		return false, nil, nil
	}
	return true, compressInfo, nil

}

func (t *CompressTask) UpdateCompressInfo(logger *logging.Logger, ctx context.Context, data *model.CompressInfo) {
	err := compressInfoService.UpdateCompressInfo(ctx, t.Engine, t.CompressInfoDao, data)
	if err != nil {
		logger.Warnf("update compress info to database error:%v,compressID:%v", err, data.Id)
		return
	}
	return
}

func (t *CompressTask) GetFilePaths(ctx *gin.Context, paths []string, pathsMap map[string]struct{}, emptyDirs []string, logger *zap.SugaredLogger, isHttp bool) bool {
	for _, path := range paths {
		// generate absolute path
		absPath := filepath.Join(t.rootPath, fsutil.TrimPrefix(path, "/"))
		// recursively get all files if the path is a directory
		p, fileInfo, err := fsutil.ReadFinalPath(absPath)
		absPath = p
		if err == nil {
			if fileInfo.IsDir() {
				err = fsutil.GetAllFileAndEmptyDir(absPath, pathsMap, &emptyDirs, logger, ctx, isHttp)
				if err != nil {
					return true
				}
			} else {
				pathsMap[absPath] = struct{}{}
			}
		} else {
			if os.IsNotExist(err) {
				msg := fmt.Sprintf("path not found, path: %s", path)
				logger.Info(msg)
				if isHttp {
					common.ErrorResp(ctx, http.StatusNotFound, commoncode.PathNotFound, msg)
				}
				return true
			}
			msg := fmt.Sprintf("Lstat error, path: %s, err: %v", path, err)
			logger.Error(msg)
			if isHttp {
				common.InternalServerError(ctx, "lstat error")
			}
			return true
		}
	}
	return false
}

type CompressParams struct {
	TargetPath          string
	TmpCompressFilePath string
	BasePath            string
	Paths               []string
	EmptyDirs           []string
	CompressCacheKey    string
	Logger              *zap.SugaredLogger
	ZipCache            *ZipCache
	CompressID          string
	Sem                 *semaphore.Weighted
	CompressInfo        *model.CompressInfo
}

func (t *CompressTask) StartCompressTask(c context.Context, params *CompressParams) {
	ctx, cancel := context.WithCancel(c)
	params.ZipCache.ContentMap.Store(params.CompressID, &ZipCacheContent{IsFinished: false, Status: 0, Cancel: cancel, TargetPath: params.TargetPath})
	t.CompressWork(ctx, params)
}

func (t *CompressTask) CancelCompressTask(c *gin.Context, logger *logging.Logger, compressID string, cache *ZipCache) error {
	value, ok := cache.ContentMap.Load(compressID)
	if !ok {
		msg := fmt.Sprintf("can not find compress task, compressID:%s", compressID)
		logger.Info(msg)
		common.ErrorResp(c, http.StatusNotFound, commoncode.CompressTaskNotFound, msg)
		return errors.New(msg)
	} else {
		cacheContent, ok := value.(*ZipCacheContent)
		if !ok {
			msg := fmt.Sprintf("zipCache.ContentMap.Load(compressID) is not *ZipCacheContent")
			logger.Info(msg)
			common.InternalServerError(c, "internal server error")
			return errors.New("internal server error")
		}
		if cacheContent.Status != int(model.CompressTaskRunning) {
			msg := fmt.Sprintf("compress task is finished, compressID:%s", compressID)
			logger.Info(msg)
			common.ErrorResp(c, http.StatusBadRequest, commoncode.CompressTaskIsFinished, msg)
			return errors.New(msg)
		}
		cacheContent.Cancel()
		logger.Infof("cancel compress task, compressID:%s", compressID)
		return nil
	}

}
func (t *CompressTask) CompressWork(ctx context.Context, params *CompressParams) {
	var cacheContent *ZipCacheContent
	value, ok := params.ZipCache.ContentMap.Load(params.CompressID)
	if ok {
		cacheContent, ok = value.(*ZipCacheContent)
		if !ok {
			params.Logger.Errorf("compressWork:zipCache.ContentMap.Load(compressID) is not *ZipCacheContent")
			return
		}
	}

	// 无论是否正常打包结束，都需要将打包任务停止
	deferFunc := func() {
		if params.Sem != nil {
			params.Sem.Release(1)
		}
		cacheContent.IsFinished = true
		if cacheContent.Err != "" {
			params.CompressInfo.ErrorMsg = cacheContent.Err
			params.CompressInfo.Status = model.CompressTaskFailed
		}
		if params.CompressInfo.Status == model.CompressTaskRunning {
			params.CompressInfo.Status = model.CompressTaskFinished
		}
		cacheContent.Status = int(params.CompressInfo.Status)
		CompressCache.Set(params.CompressCacheKey, params.ZipCache, ZipCacheExpireTime)
		params.CompressInfo.UpdateTime = time.Now()

		t.UpdateCompressInfo(params.Logger, context.Background(), params.CompressInfo)
	}
	defer deferFunc()

	select {
	case <-ctx.Done():
		params.Logger.Infof("compress task canceled, compressID:%s", params.CompressID)
		params.CompressInfo.Status = model.CompressTaskCanceled
		return
	default:
	}
	if err := os.MkdirAll(filepath.Dir(params.TmpCompressFilePath), filemode.Directory); err != nil {
		cacheContent.Err = errors.Wrap(err, "compressWork:mkdir package tmp dir").Error()
		params.Logger.Errorf("compressWork:mkdir package tmp dir err:%+v", err)
		return
	}

	outputFile, err := os.Create(params.TmpCompressFilePath)
	if err != nil {
		cacheContent.Err = errors.Wrap(err, "compressWork:create package tmp file").Error()
		params.Logger.Errorf("compressWork:create package tmp file err:%+v", err)
		return
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			params.Logger.Errorf("compressWork:close package tmp file err:%+v", err)
		}
	}(outputFile)

	zipWriter := zip.NewWriter(outputFile)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			params.Logger.Errorf("compressWork:close zipWriter err:%+v", err)
		}
	}(zipWriter)

	for _, path := range params.Paths {
		select {
		case <-ctx.Done():
			params.Logger.Infof("compress task canceled, compressID:%s", params.CompressID)
			params.CompressInfo.Status = model.CompressTaskCanceled
			return
		default:
		}
		params.Logger.Infof("add file:%s to package", path)
		if err := t.addFileToZip(path, params.BasePath, cacheContent, zipWriter); err != nil {
			params.Logger.Errorf("addFileToPackage, filename:%s, cacheContetnErr:%s, err:%v", path, cacheContent.Err, err)
			break
		}
	}

	for _, emptyDir := range params.EmptyDirs {
		select {
		case <-ctx.Done():
			params.Logger.Infof("compress task canceled, compressID:%s", params.CompressID)
			params.CompressInfo.Status = model.CompressTaskCanceled
			return
		default:
		}
		params.Logger.Infof("add emptyDir:%s to package", emptyDir)
		if err := t.addFileToZip(emptyDir, params.BasePath, cacheContent, zipWriter); err != nil {
			params.Logger.Errorf("addFileToPackage, filename:%s, cacheContentErr:%s, err:%v", emptyDir, cacheContent.Err, err)
			break
		}
	}
	select {
	case <-ctx.Done():
		params.Logger.Infof("compress task canceled, compressID:%s", params.CompressID)
		params.CompressInfo.Status = model.CompressTaskCanceled
		return
	default:
	}
	if err := zipWriter.Flush(); err != nil {
		cacheContent.Err = errors.Wrap(err, "compressWork:Flush file").Error()
		return
	}

	// 打包过程遇到err
	if cacheContent.Err != "" {
		if err := os.Remove(params.TmpCompressFilePath); err != nil {
			cacheContent.Err = errors.Wrap(err, "compressWork:remove package tmp file").Error()
			params.Logger.Errorf("compressWork:remove package tmp file err:%+v", err)
		}
		return
	}

	// 打包结束，将临时文件移动到目标路径
	if err := os.Rename(params.TmpCompressFilePath, params.TargetPath); err != nil {
		cacheContent.Err = errors.Wrap(err, "compressWork:Rename package tmp file  to target path").Error()
		params.Logger.Errorf("compressWork:rename %s to %s,err:%+v", params.TmpCompressFilePath, params.TargetPath, err)
	}

	params.Logger.Infof("compressWork:compress success, targetPath:%s", params.TargetPath)
}

func (t *CompressTask) addFileToZip(filePath, basePath string, zipCacheContent *ZipCacheContent, zipWriter *zip.Writer) (err error) {
	relPath, err := filepath.Rel(basePath, filePath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("createZipFile:Stat filePath:%s", relPath))
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("createZipFile:FileInfoHeader filePath:%s", relPath))
	}

	fileHeader, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		zipCacheContent.Err = errors.Wrap(err, fmt.Sprintf("createZipFile:FileInfoHeader filePath:%s", relPath)).Error()
		return
	}
	fileHeader.Name = relPath
	fileHeader.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(fileHeader)
	if err != nil {
		zipCacheContent.Err = errors.Wrap(err, fmt.Sprintf("addFileToZip:CreateHeader header:%+v", fileHeader)).Error()
		return
	}
	if !fileInfo.IsDir() {
		readCloser, err := fsutil.ReadAt(filePath, 0, -1)
		if err != nil {
			zipCacheContent.Err = errors.Wrap(err, fmt.Sprintf("addFileToZip:ReadAt file:%s", relPath)).Error()
			return err
		}
		defer func(readCloser io.ReadCloser) {
			err := readCloser.Close()
			if err != nil {
				zipCacheContent.Err = errors.Wrap(err, fmt.Sprintf("addFileToZip:Close file:%s", relPath)).Error()
			}
		}(readCloser)
		if _, err = io.Copy(writer, readCloser); err != nil {
			zipCacheContent.Err = errors.Wrap(err, "addFileToZip:Copy file").Error()
			return err
		}
	}

	return
}

func (t *CompressTask) GetOrCreateZipCache(userID, compressCacheKey string) *ZipCache {

	if cache, ok := CompressCache.Get(compressCacheKey); ok {
		return cache.(*ZipCache)
	}

	newCache := &ZipCache{ContentMap: &sync.Map{}}
	CompressCache.Set(compressCacheKey, newCache, ZipCacheExpireTime)
	CompressCache.OnEvicted(func(key string, value interface{}) {
		Semaphores.Delete(userID)
	})
	return newCache
}
