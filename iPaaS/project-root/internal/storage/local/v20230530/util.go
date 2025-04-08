package v20230530

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/pathchecker"
	"golang.org/x/sync/errgroup"
)

const (
	maxGoroutine                  = 16
	pathNotSafeErrorMsg           = "path is not safe"
	writeAtLengthNotMatchErrorMsg = "writeAt length not match"
)

func (s *Storage) writeAt(path string, offset, length int64, reader io.Reader, logger *logging.Logger) error {
	if flag, msg := fsutil.IsSafePath(path); !flag {
		logger.Errorf("[method]:storage.writeAt.IsSafePath, path is not safe, path: %s, err: %v", path, msg)
		return errors.New(fmt.Sprintf(pathNotSafeErrorMsg+", path: %s", path))
	}
	if err := os.MkdirAll(filepath.Dir(path), filemode.Directory); err != nil {
		logger.Errorf("[method]:storage.writeAt.MkdirOfFile, make dir of file error, path: %s, err: %v", path, err)
		return err
	}
	// open it
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Errorf("[method]:storage.WriteAt.OpenFile, open file error, path: %s, err: %v", path, err)
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.Errorf("open file error, err: %v", err)
		}
	}(f)

	// seek it
	if err = fsutil.Seek(f, offset, logger); err != nil {
		logger.Errorf("[method]:storage.WriteAt.seek, seek file error, path: %s, err: %v", path, err)
		return err
	}

	buf := s.pool.Get()
	defer s.pool.Put(buf)

	if n, err := io.CopyBuffer(f, io.LimitReader(reader, length), buf.Bytes()); err != nil {
		return err
	} else if n != length {
		logger.Warnf("[method]:storage.WriteAt.copyBuffer, length is %d, body length is %d", length, n)
		return errors.New(fmt.Sprintf(writeAtLengthNotMatchErrorMsg+"length is %d, body length is %d", length, n))
	}

	// sync it
	if err = f.Sync(); err != nil {
		msg := fmt.Sprintf("sync file error, err: %v", err)
		logger.Error(msg)
		return err
	}

	return nil
}

func (s *Storage) uploadSlice(path string, offset int64, data []byte, f *os.File, logger *logging.Logger) error {

	if flag, msg := fsutil.IsSafePath(path); !flag {
		logger.Errorf("[method]:storage.uploadSlice.IsSafePath, path is not safe, path: %s, err: %v", path, msg)
		return errors.New(msg)
	}
	if err := os.MkdirAll(filepath.Dir(path), filemode.Directory); err != nil {
		logger.Errorf("[method]:storage.uploadSlice.MkdirOfFile, make dir of file error, path: %s, err: %v", path, err)
		return err
	}

	//GetQuota or create a lock for the given file path
	fileLock, _ := s.fileLocks.LoadOrStore(path, &sync.Mutex{})
	fileLock.(*sync.Mutex).Lock()
	defer fileLock.(*sync.Mutex).Unlock()

	// seek it
	if err := fsutil.Seek(f, offset, logger); err != nil {
		logger.Errorf("[method]:storage.uploadSlice.seek, seek file error, path: %s, err: %v", path, err)
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

// GetUserIDAndAKAndHandleError 获取用户id和ak
func (s *Storage) GetUserIDAndAKAndHandleError(ctx *gin.Context) (string, string, bool, error) {
	return s.PathAccessCheckerImpl.GetUserIDAndAKAndHandleError(ctx, pathchecker.SystemURLPrefix)
}

func (s *Storage) GetTotalSize(logger *logging.Logger, paths ...string) float64 {
	var totalSize float64
	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			logger.Errorf("Error accessing path: %s", path)
			continue
		}

		if fileInfo.IsDir() {
			err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					logger.Errorf("Error accessing path: %s", filePath)
					return err
				}

				if !info.IsDir() {
					totalSize += float64(info.Size())
				}
				return nil
			})
			if err != nil {
				return 0
			}
		} else {
			totalSize += float64(fileInfo.Size())

		}
	}
	logger.Infof("total size: %f", totalSize)
	return totalSize
}

func (s *Storage) createZipFile(filePath, basePath string, zipWriter *zip.Writer, logger *logging.Logger) error {

	logger.Infof("createZipFile:filePath:%s, basePath:%s", filePath, basePath)

	relPath, err := filepath.Rel(basePath, filePath)
	if err != nil {
		return errors.Wrapf(err, "createZipFile:Rel basePath:%s, filePath:%s", basePath, filePath)
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return errors.Wrapf(err, "createZipFile:Stat filePath:%s", relPath)
	}

	fileHeader, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return errors.Wrapf(err, "createZipFile:FileInfoHeader filePath:%s", relPath)
	}
	fileHeader.Name = relPath
	fileHeader.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(fileHeader)
	if err != nil {
		return errors.Wrapf(err, "createZipFile:CreateHeader header:%+v", fileHeader)
	}

	if !fileInfo.IsDir() {

		file, err := os.Open(filePath)
		if err != nil {
			return errors.Wrapf(err, "createZipFile:Open filePath:%s", relPath)
		}
		defer func(file *os.File) {
			if closeErr := file.Close(); closeErr != nil {
				logger.Errorf("close file error, err: %v", closeErr)
			}
		}(file)

		if _, err := io.Copy(writer, file); err != nil {
			return errors.Wrap(err, "createZipFile:Copy file")
		}
	}

	return err
}

func (s *Storage) createZipFiles(basePath string, zipWriter *zip.Writer, logger *logging.Logger, paths ...string) error {
	var g errgroup.Group

	sem := make(chan struct{}, maxGoroutine)
	var mu sync.Mutex
	for _, path := range paths {
		path := path
		g.Go(func() error {
			defer func() {
				<-sem
			}()
			sem <- struct{}{}

			mu.Lock()
			defer mu.Unlock()

			return s.createZipFile(path, basePath, zipWriter, logger)
		})

	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) HandlePathContainsFileError(ctx *gin.Context, logger *logging.Logger, filePath, msg string) (bool, os.FileInfo, error) {
	fileInfo, err := os.Lstat(filePath)
	if err != nil && strings.Contains(err.Error(), fsutil.NotADirectoryErrorString) {
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.PathContainsFile, msg)
		return false, nil, err
	}
	return true, fileInfo, err
}

func (s *Storage) HandleSymbolicLink(absPath string, logger *logging.Logger, fileInfo os.FileInfo, ctx *gin.Context) (bool, string, os.FileInfo) {

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		logger.Infof("path %s is a symlink", absPath)
		sourceFilePath, err := os.Readlink(absPath)
		logger.Infof("source file path: %s", sourceFilePath)
		if err != nil {
			msg := fmt.Sprintf("read symlink error, err: %v", err)
			logger.Errorf(msg)
			common.InternalServerError(ctx, "internal server error")
			return false, "", nil
		}
		if filepath.IsAbs(sourceFilePath) {
			absPath = sourceFilePath
		} else {
			absPath = filepath.Join(filepath.Dir(absPath), sourceFilePath)
		}
		fileInfo, err = os.Lstat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				msg := fmt.Sprintf("file or directory not found, path: %v", strings.TrimPrefix(absPath, s.rootPath+"/"))
				logger.Info(msg)
				common.ErrorResp(ctx, http.StatusNotFound, commoncode.PathNotFound, msg)
				return false, "", nil
			}
			msg := fmt.Sprintf("stat file error, err: %v , path: %v", err, strings.TrimPrefix(absPath, s.rootPath+"/"))
			logger.Errorf(msg)
			common.InternalServerError(ctx, "internal server error")
			return false, "", nil
		}

	}
	return true, absPath, fileInfo
}
