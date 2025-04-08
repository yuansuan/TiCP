package fsutil

import (
	"errors"
	"fmt"
	"github.com/coreos/etcd/pkg/ioutil"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	NotADirectoryErrorString = "not a directory"
	MaxFileNameLength        = 255  // Linux file system limit for file name length
	MaxFilePathLength        = 4000 // Linux file system limit for file path length
)

var (
	ErrFileNameTooLong = errors.New("file name too long")
	ErrPathTooLong     = errors.New("file path too long")
)

// emptyReadCloser a mock ReadCloser
type emptyReadCloser struct{}

// check the length of filename and path is valid
func ValidateFileNameLength(path string) error {
	if len(filepath.Base(path)) > MaxFileNameLength {
		return ErrFileNameTooLong
	} else if len(path) > MaxFilePathLength {
		return ErrPathTooLong
	}
	return nil
}

func (emptyReadCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (emptyReadCloser) Close() error {
	return nil
}

// ReadAt  ...
func ReadAt(path string, offset int64, length int64) (io.ReadCloser, error) {
	if ok, err := IsSafePath(path); !ok {
		msg := fmt.Sprintf("path is not a safe path, path: %s, err: %s", path, err)
		logging.Default().Error(msg)
		return nil, errors.New(err)
	}

	// if file is symlink, return an empty reader
	fileInfo, err := os.Lstat(path)
	if err != nil {
		msg := fmt.Sprintf("lstat error, path: %s, err: %v", path, err)
		logging.Default().Errorf(msg)
		return nil, errors.New(msg)
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		return &emptyReadCloser{}, nil
	}

	// open it
	file, err := os.Open(path)
	if err != nil {
		msg := fmt.Sprintf("file open error, path: %s, err: %v", path, err)
		logging.Default().Errorf(msg)
		return nil, errors.New(msg)
	}

	if fileInfo.IsDir() {
		err := file.Close()
		if err != nil {
			logging.Default().Infof("file close error, err: %v", err)
			return nil, err
		}
		msg := fmt.Sprintf("%s is not a file (is a folder)", path)
		return nil, errors.New(msg)
	}

	// seek it
	ret, err := file.Seek(offset, io.SeekStart)
	if err != nil {
		logging.Default().Errorf("file seek error, path: %s, offset: %d, err: %v", path, offset, err)
		return nil, err
	}
	if ret != offset {
		msg := fmt.Sprintf("seeked but offset not same, err: %v", err)
		logging.Default().Errorf(msg)
		return nil, errors.New(msg)
	}

	reader := (io.Reader)(file)
	if length >= 0 {
		reader = io.LimitReader(file, length)
	}

	// read it
	return ioutil.ReaderAndCloser{
		Reader: reader,
		Closer: file,
	}, nil
}

// Move ...
func Move(srcPath, destPath string, logger *logging.Logger) error {
	if flag, msg := IsSafePath(srcPath); !flag {
		logger.Errorf("srcPath is not a safe path, src path: %s, err: %v", srcPath, msg)
		return errors.New(msg)
	}
	if flag, msg := IsSafePath(destPath); !flag {
		logger.Errorf("destPath is not a safe path, dest path: %s, err: %v", destPath, msg)
		return errors.New(msg)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), filemode.Directory); err != nil {
		logger.Errorf("make dir of file error, path: %s, err: %v", destPath, err)
		return err
	}

	if err := os.Rename(srcPath, destPath); err != nil {
		logger.Errorf("rename err, src path: %s, dest path: %s, err: %v", srcPath, destPath, err)
		return err
	}

	return nil
}

// Create ...
func Create(path string, size int64, logger *logging.Logger) error {
	file, err := CreateAndReturnFile(path, size, logger)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		if file == nil {
			return
		}
		err := file.Close()
		if err != nil {
			logger.Errorf("close file error, path: %s, err: %v", path, err)
		}
	}(file)

	return nil
}

// CreateAndReturnFile ...
func CreateAndReturnFile(path string, size int64, logger *logging.Logger) (*os.File, error) {
	if flag, msg := IsSafePath(path); !flag {
		logger.Errorf("path is not safe, path: %s, err: %v", path, msg)
		return nil, errors.New(msg)
	}

	if err := os.MkdirAll(filepath.Dir(path), filemode.Directory); err != nil {
		logger.Errorf("make dir of file error, path: %s, err: %v", path, err)
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.Errorf("open file error, path: %s, size: %d, err: %v", path, size, err)
		return nil, err
	}

	if err = file.Truncate(size); err != nil {
		logger.Errorf("truncate error, path: %s, size: %d, err: %v", path, size, err)
		return nil, err
	}

	return file, nil
}

// Seek ...
func Seek(f *os.File, offset int64, logger *logging.Logger) error {
	r, err := f.Seek(offset, io.SeekStart)
	if err != nil {
		logger.Errorf("file seek err: %v", err)
		return err
	}
	if r != offset {
		logger.Errorf("seeked but offset not same, err: %v", err)
		return err
	}
	return nil
}

// Ls ...
func Ls(path string, pageSize int64, logger *logging.Logger) ([]os.FileInfo, error) {
	size := pageSize
	if pageSize <= 0 {
		size = 500
	}

	var result []os.FileInfo

	for i := int64(0); ; i += size {
		infos, _, _, err := LsPage(path, nil, i, size, logger)
		if err != nil {
			logger.Warnf("error getting page info: %v, path:%v, pageOffset:%v ,pageSize:%v", err, path, i, size)
			return nil, err
		}
		for _, fo := range infos {
			result = append(result, fo)
		}
		if len(infos) < int(size) {
			break
		}
	}

	return result, nil
}

// LsPage ...
func LsPage(path string, filterRegexpList []*regexp.Regexp, pageOffset, pageSize int64, logger *logging.Logger) (result []os.FileInfo, nextMarker int64, total int64, err error) {

	// 检查路径安全性
	if flag, msg := IsSafePath(path); !flag {
		logger.Warnf("path is not safe, path: %s, err: %v", path, msg)
		return nil, 0, 0, errors.New(msg)
	}

	// 读取目录内容
	infos, err := os.ReadDir(path)
	if err != nil {
		logger.Warnf("read dir err, path: %s, err: %v", path, err)
		return nil, 0, 0, err
	}

	// 计算起始和结束索引
	startIndex := max(pageOffset, 0)
	endIndex := min(startIndex+pageSize, int64(len(infos)))
	// 初始化结果数组
	if startIndex < endIndex {
		result = make([]os.FileInfo, 0, endIndex-startIndex)
	}

	// 遍历目录内容
	count := int64(0)
	for i, info := range infos {
		// 匹配过滤条件
		matched := false
		for _, filterRegexp := range filterRegexpList {
			if filterRegexp != nil && filterRegexp.MatchString(info.Name()) {
				matched = true
				break
			}
		}
		if matched {
			continue
		}

		count++

		// 判断是否在分页范围内
		if count-1 >= startIndex && count-1 < endIndex {
			fi, err := infos[i].Info()
			if os.IsNotExist(err) {
				count-- // 文件不存在时减少计数
				continue
			}
			if err != nil {
				logger.Errorf("get file info err, path: %s, err: %v", path, err)
				return nil, 0, 0, err
			}

			// 处理符号链接文件
			if fi.Mode()&os.ModeSymlink != 0 {
				f, err := ReadSymlinksFile(path, fi.Name())
				if os.IsNotExist(err) {
					count-- // 符号链接文件不存在时减少计数
					continue
				}
				if err != nil {
					logger.Errorf("get file info err, path: %s, err: %v", path, err)
					return nil, 0, 0, err
				}
				fi = &linkFileInfo{
					name:    fi.Name(),
					size:    f.Size(),
					mode:    fi.Mode(),
					modTime: fi.ModTime(),
					isDir:   f.IsDir(),
					sys:     fi.Sys(),
				}
				logger.Debugf("symlink file info: %v", fi)
			}

			result = append(result, fi)
			nextMarker = count
		}
	}

	// 如果当前页面已经到达最后一页，设置 nextMarker 为 -1
	if nextMarker == count {
		nextMarker = -1
	}

	return result, nextMarker, count, nil
}

// 辅助函数：取最大值
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// 辅助函数：取最小值
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// ReadSymlinksFile ...
func ReadSymlinksFile(path, name string) (os.FileInfo, error) {
	targetPath, err := filepath.EvalSymlinks(filepath.Join(path, name))
	if err != nil {
		if os.IsNotExist(err) {
			logging.Default().Infof("eval symlinks not exist, path: %s, err: %v", path, err)
			return nil, err
		}
		logging.Default().Errorf("eval symlinks error, path: %s, err: %v", path, err)
		return nil, err
	}
	targetInfo, err := os.Stat(targetPath)
	if err != nil {
		logging.Default().Errorf("stat error, path: %s, err: %v", path, err)
		return nil, err
	}
	return targetInfo, nil
}

// the return path will be the real path of symbolic link
func ReadFinalPath(path string) (string, os.FileInfo, error) {
	// if path is a symbolic link, read the real path
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return "", nil, err
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		// to its final destination
		realPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return "", nil, err
		}
		fileInfo, err = os.Lstat(realPath)
		if err != nil {
			return "", nil, err
		}
		if !fileInfo.IsDir() {
			return path, fileInfo, nil
		}
		path = realPath
		return path, fileInfo, nil
	}
	return path, fileInfo, nil
}

func GetAllFileAndEmptyDir(path string, pathsMap map[string]struct{}, emptyDirs *[]string, logger *logging.Logger, ctx *gin.Context, isHttp bool) error {

	err := filepath.WalkDir(path, func(filePath string, file os.DirEntry, err error) error {
		if err != nil {
			msg := fmt.Sprintf("Error walking directory, path: %s, err: %v", path, err)
			logger.Error(msg)
			if isHttp {
				common.InternalServerError(ctx, "walk dir err")
			}
			return err
		}

		if file.IsDir() && strings.HasPrefix(file.Name(), ".") {
			return filepath.SkipDir
		}

		if file.IsDir() {
			files, err := os.ReadDir(filePath)
			if err != nil {
				msg := fmt.Sprintf("Error reading directory, path: %s, err: %v\n", filePath, err)
				logger.Error(msg)
				if isHttp {
					common.InternalServerError(ctx, "read dir err")
				}
				return err
			}
			if len(files) == 0 {
				*emptyDirs = append(*emptyDirs, filePath)
			}
		}

		if file.Type()&os.ModeSymlink != 0 {
			if file.IsDir() {
				logger.Infof("dir path %s is a symlink", filePath)
				return nil
			}
			sourceFilePath := filePath
			var absPath string
			if filepath.IsAbs(sourceFilePath) {
				absPath = sourceFilePath
			} else {
				absPath = filepath.Join(filepath.Dir(filePath), sourceFilePath)
			}

			pathsMap[absPath] = struct{}{}
			return nil
		}

		if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			pathsMap[filePath] = struct{}{}
		}

		return nil
	})

	return err
}

// Truncate ...
func Truncate(path string, size int64, logger *logging.Logger) error {
	if flag, msg := IsSafePath(path); !flag {
		logger.Errorf("path is not safe, path: %s, err: %v", path, msg)
		return errors.New(msg)
	}

	if err := os.MkdirAll(filepath.Dir(path), filemode.Directory); err != nil {
		logger.Errorf("make dir of file error, path: %s, err: %v", path, err)
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.Errorf("open file error, path: %s, size: %d, err: %v", path, size, err)
		return err
	}

	defer func(file *os.File) {
		if file == nil {
			return
		}
		err := file.Close()
		if err != nil {
			logger.Errorf("close file error, path: %s, err: %v", path, err)
		}
	}(file)

	if err = file.Truncate(size); err != nil {
		logger.Errorf("truncate error, path: %s, size: %d, err: %v", path, size, err)
		return err
	}

	return nil
}
