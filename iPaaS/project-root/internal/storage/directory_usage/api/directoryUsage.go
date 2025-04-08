package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/pathchecker"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/directoryUsage"
	directoryUsageService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/directoryUsage"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
	"xorm.io/xorm"
)

type DirectoryUsage struct {
	Engine            *xorm.Engine
	DirectoryUsageDao dao.DirectoryUsageDao
	CancelMap         sync.Map
	rootPath          string

	pathchecker.PathAccessCheckerImpl
}

func NewDirectoryUsage(directoryUsageDao dao.DirectoryUsageDao, engine *xorm.Engine, pathchecker pathchecker.PathAccessCheckerImpl, rootPath string) *DirectoryUsage {
	if directoryUsageDao == nil {
		return nil
	}

	return &DirectoryUsage{
		Engine:                engine,
		DirectoryUsageDao:     directoryUsageDao,
		CancelMap:             sync.Map{},
		PathAccessCheckerImpl: pathchecker,
		rootPath:              rootPath,
	}
}

// GetUserIDAndAKAndHandleError 获取用户id和ak
func (t *DirectoryUsage) GetUserIDAndAKAndHandleError(ctx *gin.Context) (string, string, bool, error) {
	return t.PathAccessCheckerImpl.GetUserIDAndAKAndHandleError(ctx, pathchecker.SystemURLPrefix)
}

func (t *DirectoryUsage) Recover() {
	ctx := context.Background()
	logger := logging.GetLogger(ctx).With("func", "DirectoryUsage.Recover")
	directoryUsageList, err := directoryUsageService.ListDirectoryUsage(ctx, t.Engine, t.DirectoryUsageDao)
	if err != nil {
		logger.Warnf("get unfinished directory usage task from database error:%v", err)
		return
	}
	logger.Infof("unfinished directory usage task number:%v", len(directoryUsageList))
	for _, directoryUsage := range directoryUsageList {
		logger.Infof("recover directory usage task,taskID:%v", directoryUsage.Id)
		go t.StartCalculateDirectoryUsage(context.Background(), logger, directoryUsage.Id, filepath.Join(t.rootPath, fsutil.TrimPrefix(directoryUsage.Path, "/")), directoryUsage.UserID, true)
	}
}
func (t *DirectoryUsage) StartCalculateDirectoryUsage(c context.Context, logger *logging.Logger, taskID, path, userID string, skipInsert bool) {
	ctx, cancel := context.WithCancel(c)
	t.CancelMap.Store(taskID, cancel)
	logger.Infof("start calculate directory usage, path: %v,taskID: %v", path, taskID)
	t.DirectoryUsage(path, logger, ctx, taskID, userID, skipInsert)
}

func (t *DirectoryUsage) DirectoryUsage(path string, logger *logging.Logger, ctx context.Context, taskID, userID string, skipInsert bool) {
	defer func() {
		t.CancelMap.Delete(taskID)
	}()
	var err error
	if !skipInsert {
		err = directoryUsage.InsertDirectoryUsage(context.Background(), t.Engine, t.DirectoryUsageDao, &model.DirectoryUsage{
			Id:         taskID,
			UserID:     userID,
			Path:       fsutil.TrimPrefix(path, t.rootPath),
			Status:     model.DirectoryUsageTaskCalculating,
			CreateTime: time.Now(),
		})
	}
	if err != nil {
		logger.Warnf("insert directory usage to database error:%v,taskID:%v", err, taskID)
		return
	}
	totalSize := int64(0)
	totalLogicSize := int64(0)
	if err = t.walkDir(ctx, path, &totalSize, &totalLogicSize); err != nil {
		if err == context.Canceled {
			err := directoryUsage.UpdateDirectoryUsage(context.Background(), t.Engine, t.DirectoryUsageDao, &model.DirectoryUsage{
				Id:         taskID,
				Status:     model.DirectoryUsageTaskCanceled,
				UpdateTime: time.Now(),
			})
			if err != nil {
				logger.Warnf("update directory usage to database error:%v,taskID:%v", err, taskID)
			}
			logger.Infof("calculate directory usage canceled, path:%s", path)
		} else {
			logger.Warnf("Error walking path: %v", err)
			err := directoryUsage.UpdateDirectoryUsage(context.Background(), t.Engine, t.DirectoryUsageDao, &model.DirectoryUsage{
				Id:         taskID,
				Status:     model.DirectoryUsageTaskFailed,
				ErrMsg:     err.Error(),
				UpdateTime: time.Now(),
			})
			if err != nil {
				logger.Warnf("update directory usage to database error:%v,taskID:%v", err, taskID)
			}
		}
		return
	}

	err = directoryUsage.UpdateDirectoryUsage(context.Background(), t.Engine, t.DirectoryUsageDao, &model.DirectoryUsage{
		Id:         taskID,
		Status:     model.DirectoryUsageTaskFinished,
		Size:       totalSize,
		LogicSize:  totalLogicSize,
		UpdateTime: time.Now(),
	})
	if err != nil {
		logger.Warnf("update directory usage to database error:%v,taskID:%v", err, taskID)
	}
	logger.Infof("calculate directory usage success, path:%s,size:%v", path, totalSize)
}

func (t *DirectoryUsage) walkDir(ctx context.Context, dir string, totalSize, totalLogicSize *int64) error {
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}
	*totalSize += dirInfo.Size()

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			subDir := filepath.Join(dir, file.Name())
			if err := t.walkDir(ctx, subDir, totalSize, totalLogicSize); err != nil {
				return err
			}
		} else {
			select {
			case <-ctx.Done():
				return context.Canceled
			default:
				// csp团队要求记录软链接指向文件的真实大小
				_, fileInfo, err := fsutil.ReadFinalPath(filepath.Join(dir, file.Name()))
				if err != nil {
					return err
				}
				*totalSize += fileInfo.Size()

				fileLogicInfo, err := file.Info()
				if err != nil {
					return err
				}
				*totalLogicSize += fileLogicInfo.Size()
			}
		}
	}

	return nil
}

func (t *DirectoryUsage) CancelCalculateDirectoryUsage(c *gin.Context, logger *logging.Logger, taskID string) error {
	value, ok := t.CancelMap.Load(taskID)
	if !ok {
		msg := fmt.Sprintf("Calculate directory usage task not found, taskID:%s", taskID)
		logger.Info(msg)
		common.ErrorResp(c, http.StatusNotFound, commoncode.DirectoryUsageTaskNotFound, msg)
		return errors.New(msg)
	} else {
		cancelFunc, ok := value.(context.CancelFunc)
		if !ok {
			msg := fmt.Sprintf("type assertion failed, taskID:%s", taskID)
			logger.Info(msg)
			common.InternalServerError(c, "internal server error")
			return errors.New("internal server error")
		}
		cancelFunc()
		logger.Infof("cancel calculate directory usage, taskID:%s", taskID)
		return nil
	}

}
