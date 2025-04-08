package fsutil

import (
	"bytes"
	_errors "errors"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil/filemode"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/sbs"
)

var (
	// ErrUnsupportedOperation 表示不支持的文件系统操作
	ErrUnsupportedOperation = _errors.New("filesystem: unsupported operation")

	// root 系统的默认根文件系统
	root *Filesystem
)

func init() {
	root = &Filesystem{
		realpath: "/",
		Fs:       afero.NewBasePathFs(afero.NewOsFs(), "/"),
	}
}

// MkdirAuto 检查目录是否存在, 如果不存在的话则尝试创建
func MkdirAuto(dirname string, mode os.FileMode) error {
	return root.MkdirAuto(dirname, mode)
}

// IsDir 判断路径是否是一个目录
func IsDir(dirname string) (bool, error) {
	return root.IsDir(dirname)
}

// Chroot 切换到指定目录之内
func Chroot(dirname string) (*Filesystem, error) {
	return root.Chroot(dirname)
}

// FileExists 检查文件是否存在
func FileExists(filename string) (bool, error) {
	return root.FileExists(filename)
}

// Filesystem 是一个基于 afero.Fs 封装的文件系统抽象层
type Filesystem struct {
	realpath string
	afero.Fs
}

// MkdirAuto 检查目录是否存在, 如果不存在的话则尝试创建
func (fs *Filesystem) MkdirAuto(dirname string, mode os.FileMode) error {
	exists, err := afero.DirExists(fs, dirname)
	if err != nil {
		return err
	} else if !exists {
		err = fs.MkdirAll(dirname, mode)
	}

	return err
}

// IsDir 判断路径是否是目录路径
func (fs *Filesystem) IsDir(dirname string) (bool, error) {
	isDir, err := afero.IsDir(fs, dirname)
	if err != nil {
		return false, err
	}
	return isDir, nil
}

// RealPath 获取该虚拟文件系统下真实的文件路径
func (fs *Filesystem) RealPath(name string) (string, error) {
	if len(fs.realpath) != 0 {
		return Join(fs.realpath, name), nil
	}
	return "", ErrUnsupportedOperation
}

// Chroot 将根目录切换到某个子目录上
func (fs *Filesystem) Chroot(dirname string) (*Filesystem, error) {
	if exists, err := afero.DirExists(fs, dirname); err != nil {
		return nil, err
	} else if !exists {
		return nil, os.ErrNotExist
	}

	return &Filesystem{
		realpath: Join(fs.realpath, dirname),
		Fs:       afero.NewBasePathFs(fs, dirname),
	}, nil
}

// FileExists 检查指定的文件是否存在
func (fs *Filesystem) FileExists(filename string) (bool, error) {
	info, err := fs.Stat(filename)
	if err == nil {
		return !info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// DirExists 检查指定的文件夹是否存在
func (fs *Filesystem) DirExists(dirname string) (bool, error) {
	return afero.DirExists(fs, dirname)
}

// ReadFileAsString 打开文件并读取内容转换成字符串后返回
func (fs *Filesystem) ReadFileAsString(filename string) (string, error) {
	exists, err := fs.FileExists(filename)
	if err != nil {
		return "", err
	} else if !exists {
		return "", os.ErrNotExist
	}

	f, err := fs.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return sbs.String(data), nil
}

// WriteToFile 将内容写入到文件中
func (fs *Filesystem) WriteToFile(filename string, data []byte) error {
	f, err := fs.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, filemode.RegularFile)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = io.Copy(f, bytes.NewReader(data))
	return err
}

// Join 将多个路径拼接起来
func Join(root string, vs ...string) string {
	for _, s := range vs {
		s = filepath.Clean(filepath.Join("/", s))
		root = filepath.Join(root, s)
	}
	return root
}

// SysFile 将 afero.File 接口类型的文件转换为系统文件抽象
func SysFile(f afero.File) (*os.File, bool) {
	for {
		if sys, ok := f.(*os.File); ok {
			return sys, true
		}

		switch f.(type) {
		case *afero.BasePathFile:
			f = f.(*afero.BasePathFile).File
		default:
			return nil, false
		}
	}
}
