package image

import (
	_errors "errors"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil/filemode"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/locker"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/sbs"
)

var (
	// ErrInvalidLocallyLocator 无效的镜像定位器
	ErrInvalidLocallyLocator = _errors.New("image.locally: invalid locally image locator")
	// ErrLocallyImageNotFound 本地找不到指定的镜像
	ErrLocallyImageNotFound = _errors.New("image.locally: image not found")
)

// Locally 一个本地的镜像
type Locally struct {
	name    string
	tag     string
	hash    string
	tagFs   *fsutil.Filesystem
	imageFs *fsutil.Filesystem
}

// Exists 返回当前镜像是否已经存在
func (l *Locally) Exists() (bool, error) {
	exists, err := l.tagFs.FileExists(l.hash)
	return exists, errors.Wrap(err, "image.locally")
}

// Locator 生成当前镜像的定位器
func (l *Locally) Locator() Locator {
	return Locate(l.name, WithLocateTag(l.tag), WithLocateHash(l.hash))
}

// RealPath 返回当前镜像的实际位置
func (l *Locally) RealPath() (string, error) {
	return l.tagFs.RealPath(l.hash)
}

// Open 打开镜像文件
func (l *Locally) Open() (io.ReadSeekCloser, error) {
	return l.tagFs.Open(l.hash)
}

// RewriteDefaults 重写本地镜像的默认值
func (l *Locally) RewriteDefaults(raw Locator) error {
	if len(raw.Tag()) == 0 {
		if err := l.imageFs.WriteToFile(_LockFile, sbs.Bytes(l.tag)); err != nil {
			return errors.Wrap(err, "image.locally")
		}
	}

	if len(raw.Hash()) == 0 {
		if err := l.tagFs.WriteToFile(_LockFile, sbs.Bytes(l.hash)); err != nil {
			return errors.Wrap(err, "image.locally")
		}
	}

	return nil
}

// Writer 获取一个写入器
func (l *Locally) Writer() (io.WriteCloser, error) {
	f, err := l.tagFs.OpenFile("_"+l.hash, os.O_CREATE|os.O_RDWR|os.O_TRUNC, filemode.RegularFile)
	if err != nil {
		return nil, errors.Wrap(err, "image.locally")
	}

	return f, nil
}

// Occupy 将当前临时镜像标记为可用状态
func (l *Locally) Occupy() error {
	if exists, err := l.tagFs.FileExists("_" + l.hash); err != nil {
		return errors.Wrap(err, "image.locally")
	} else if exists {
		if err = l.tagFs.Remove(l.hash); err != nil && !os.IsNotExist(err) {
			return errors.Wrap(err, "image.locally")
		}
		return errors.Wrap(l.tagFs.Rename("_"+l.hash, l.hash), "image.locally")
	}
	return nil
}

// WriterLocker 创建一个写入锁
func (l *Locally) WriterLocker() locker.Locker {
	return locker.NewFileLocker(func() (*os.File, func(), error) {
		lockfile := "_" + l.hash + ".lock"
		f, err := l.tagFs.OpenFile(lockfile, os.O_CREATE|os.O_RDWR, filemode.RegularFile)
		if err != nil {
			return nil, nil, err
		}
		// closing file on cleanup closure

		sys, ok := fsutil.SysFile(f)
		if !ok {
			return nil, nil, errors.New("unsupported filesystem")
		}

		return sys, func() {
			_ = f.Close()
			_ = l.tagFs.Remove(lockfile)
		}, nil
	})
}

// NewLocally 创建一个新的本地镜像
func NewLocally(locator Locator, root *fsutil.Filesystem) (*Locally, error) {
	if len(locator.Tag()) == 0 || len(locator.Hash()) == 0 {
		return nil, ErrInvalidLocallyLocator
	}

	fully := fsutil.Join(locator.Name(), locator.Tag())
	if err := root.MkdirAuto(fully, filemode.Directory); err != nil {
		return nil, errors.Wrap(err, "image.locally")
	}

	imageFs, err := root.Chroot(locator.Name())
	if err != nil {
		return nil, errors.Wrap(err, "image.locally")
	}

	tagFs, err := imageFs.Chroot(locator.Tag())
	if err != nil {
		return nil, errors.Wrap(err, "image.locally")
	}

	return &Locally{
		name:    locator.Name(),
		tag:     locator.Tag(),
		hash:    locator.Hash(),
		tagFs:   tagFs,
		imageFs: imageFs,
	}, nil
}

// LoadLocally 从本地加载一个镜像文件
func LoadLocally(locator Locator, root *fsutil.Filesystem) (*Locally, error) {
	imageFs, err := root.Chroot(locator.Name())
	if err != nil {
		return nil, ErrLocallyImageNotFound
	}

	l := &Locally{name: locator.Name(), imageFs: imageFs, tag: locator.Tag(), hash: locator.Hash()}
	if len(locator.Tag()) == 0 {
		l.tag, err = l.imageFs.ReadFileAsString(_LockFile)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, ErrNoDefaultImageTagFound
			}
			return nil, errors.Wrap(err, "image.locally")
		}
	}

	l.tagFs, err = l.imageFs.Chroot(l.tag)
	if err != nil {
		return nil, ErrLocallyImageNotFound
	}

	if len(locator.Hash()) == 0 {
		l.hash, err = l.tagFs.ReadFileAsString(_LockFile)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, ErrNoDefaultImageHashFound
			}
			return nil, errors.Wrap(err, "image.locally")
		}
	}

	return l, nil
}
