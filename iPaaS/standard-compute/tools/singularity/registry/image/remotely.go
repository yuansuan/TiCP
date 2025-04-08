package image

import (
	"context"
	_errors "errors"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/locker"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xhttp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xoss"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xurl"
)

var (
	// ErrRemotelyImageNotFound 无法在远程找到镜像
	ErrRemotelyImageNotFound = _errors.New("registry.remote: unable to find image remotely")
	// ErrNoDefaultImageTagFound 找不到镜像的默认标签
	ErrNoDefaultImageTagFound = _errors.New("registry.remote: no default image tag found")
	// ErrNoDefaultImageHashFound 找不到镜像的默认标签
	ErrNoDefaultImageHashFound = _errors.New("registry.remote: no default image hash found")
	// ErrFullQualityLocatorRequired 表示需要一个全路径的定位器
	ErrFullQualityLocatorRequired = _errors.New("registry.remote: full quality locator required")
)

// Remotely 表示一个远端的镜像
type Remotely struct {
	name string
	tag  string
	hash string
	oss  *xoss.ObjectStorageService

	h *xhttp.Client
}

// Locator 生成当前远程镜像的资源定位器
func (r *Remotely) Locator() Locator {
	return Locate(r.name, WithLocateTag(r.tag), WithLocateHash(r.hash))
}

// Name 返回镜像的名字
func (r *Remotely) Name() string {
	return r.name
}

// Tag 返回镜像的标签
func (r *Remotely) Tag() string {
	return r.tag
}

// Hash 返回镜像的哈希值
func (r *Remotely) Hash() string {
	return r.hash
}

// Verify 检查这个镜像是否有效
func (r *Remotely) Verify(ctx context.Context) error {
	_, err := r.oss.HeadObject(ctx, r.objectKey())
	if err != nil {
		if xoss.IsObjectNotExists(err) {
			return ErrRemotelyImageNotFound
		}
		return errors.Wrap(err, "image.remotely")
	}
	return nil
}

// Reader 打开远程镜像并返回, 注意一定要关闭对象
func (r *Remotely) Reader(ctx context.Context) (io.ReadCloser, int64, error) {
	resp, err := r.oss.GetObject(ctx, r.objectKey())
	if err != nil {
		return nil, 0, errors.Wrap(err, "image.remotely")
	}

	return resp.Body, *resp.ContentLength, nil
}

// WriteFrom 从 rs 中读取数据并上传到远端仓库
func (r *Remotely) WriteFrom(ctx context.Context, rd io.Reader) error {
	_, err := r.oss.UploadObject(ctx, r.tempObjectKey(), rd)
	if err != nil {
		return errors.Wrap(err, "image.remotely")
	}
	return nil
}

// Occupy 重新发布镜像版本
func (r *Remotely) Occupy(ctx context.Context) (err error) {
	if _, err = r.oss.HeadObject(ctx, r.tempObjectKey()); err == nil {
		if _, err = r.oss.CopyObject(ctx, r.tempObjectKey(), r.objectKey()); err == nil {
			_, err = r.oss.DeleteObject(ctx, r.tempObjectKey())
		}
	}
	return
}

// WriterLocker 创建写入锁
func (r *Remotely) WriterLocker() locker.Locker {
	return locker.NewS3Locker(r.oss, r.tempObjectKey()+".lock")
}

// OverwriteDefaults 覆盖默认配置
func (r *Remotely) OverwriteDefaults(ctx context.Context, raw Locator) error {
	if len(raw.Hash()) == 0 {
		hashLock := xurl.Join(r.name, r.tag, _LockFile)
		_, err := r.oss.PutObject(ctx, hashLock, strings.NewReader(r.hash))
		if err != nil {
			return errors.Wrap(err, "image.remotely")
		}
	}

	if len(raw.Tag()) == 0 {
		tagLock := xurl.Join(r.name, _LockFile)
		_, err := r.oss.PutObject(ctx, tagLock, strings.NewReader(r.tag))
		if err != nil {
			return errors.Wrap(err, "image.remotely")
		}
	}

	return nil
}

// objectKey 返回该镜像的资源位置
func (r *Remotely) objectKey() string {
	// convention: /<name>/<tag>/<hash>
	return xurl.Join(r.name, r.tag, r.hash)
}

// tempObjectKey 镜像的临时文件位置
func (r *Remotely) tempObjectKey() string {
	return xurl.Join(r.name, r.tag, "_"+r.hash)
}

// LoadRemotely 从远程加载镜像
func LoadRemotely(ctx context.Context, locator Locator, oss *xoss.ObjectStorageService) (r *Remotely, err error) {
	r = &Remotely{oss: oss, name: locator.Name(), tag: locator.Tag(), hash: locator.Hash()}
	if len(locator.Tag()) == 0 {
		if r.tag, err = r.oss.GetAsString(ctx, xurl.Join(r.name, _LockFile)); err != nil {
			if xoss.IsObjectNotExists(err) {
				return nil, ErrNoDefaultImageTagFound
			}
			return nil, errors.Wrap(err, "image.remotely")
		}
	}

	if len(locator.Hash()) == 0 {
		if r.hash, err = r.oss.GetAsString(ctx, xurl.Join(r.name, r.tag, _LockFile)); err != nil {
			if xoss.IsObjectNotExists(err) {
				return nil, ErrNoDefaultImageHashFound
			}
			return nil, errors.Wrap(err, "image.remotely")
		}
	}

	if err = r.Verify(ctx); err != nil {
		return nil, err
	}

	return
}

// NewRemotely 创建一个新的远程镜像
func NewRemotely(_ context.Context, locator Locator, oss *xoss.ObjectStorageService) (*Remotely, error) {
	if len(locator.Tag()) == 0 || len(locator.Hash()) == 0 {
		return nil, ErrFullQualityLocatorRequired
	}

	return &Remotely{name: locator.Name(), tag: locator.Tag(), hash: locator.Hash(), oss: oss}, nil
}
