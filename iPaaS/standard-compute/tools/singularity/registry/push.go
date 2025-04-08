package registry

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	_errors "errors"
	"github.com/schollz/progressbar/v3"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry/image"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/locker"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xhttp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xio"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xoss"
)

const (
	// _DefaultTagName 默认的镜像标签
	_DefaultTagName = "latest"
)

var (
	// ErrPushInProgress 有另外一个进程正在上传中
	ErrPushInProgress = _errors.New("registry.push: another still in progress")
)

// PushOption 推送镜像的参数
type PushOption func(*Pusher) error

// WithPushLocalFile 使用本地文件上传而不是从仓库中拉取
func WithPushLocalFile(filename string) PushOption {
	return func(p *Pusher) error {
		p.localFile = filename
		return nil
	}
}

// WithPushLocalAlias 使用本地的镜像进行上传
func WithPushLocalAlias(alias image.Locator) PushOption {
	return func(p *Pusher) error {
		p.alias = alias
		return nil
	}
}

// WithPushForce 是否强制上传镜像
func WithPushForce(force bool) PushOption {
	return func(p *Pusher) error {
		p.force = force
		return nil
	}
}

// WithPushContext 设置一个上下文对象
func WithPushContext(ctx context.Context) PushOption {
	return func(p *Pusher) error {
		p.ctx = ctx
		return nil
	}
}

// WithPushLogger 设置日志记录器
func WithPushLogger(log Logger) PushOption {
	return func(p *Pusher) error {
		p.log = log
		return nil
	}
}

// WithPushProgressBar 设置输出推送进度条
func WithPushProgressBar(opts []progressbar.Option) PushOption {
	return func(p *Pusher) error {
		p.pbOpts = opts
		return nil
	}
}

// WithPushProgressSubscribe 设置推送进度监控
func WithPushProgressSubscribe(aph func(*AtomicProgress)) PushOption {
	return func(p *Pusher) error {
		p.aph = aph
		return nil
	}
}

// Push 推送镜像到远端仓库中
func (c *client) Push(locator image.Locator, options ...PushOption) (*image.Remotely, error) {
	pusher := &Pusher{
		oss:     c.oss,
		remote:  c.remote,
		storage: c.storage,
		cfg:     c.cfg.Registry,
		ctx:     context.Background(),
		log:     NewDiscardLogger(),
	}

	for _, option := range options {
		if err := option(pusher); err != nil {
			return nil, errors.Wrap(err, "registry.push")
		}
	}

	return pusher.Push(locator)
}

// Pusher 镜像推送工具
type Pusher struct {
	ctx     context.Context
	oss     *xoss.ObjectStorageService
	storage *fsutil.Filesystem
	remote  *xhttp.Client
	cfg     *config.ObjectStorageService

	force     bool
	localFile string
	alias     image.Locator

	log    Logger
	aph    func(*AtomicProgress)
	pbOpts []progressbar.Option
}

// Push 执行推送镜像任务
func (p *Pusher) Push(locator image.Locator) (remotely *image.Remotely, err error) {
	var r io.ReadSeekCloser
	if len(p.localFile) == 0 {
		var locally *image.Locally
		if p.alias != nil {
			locally, err = image.LoadLocally(p.alias, p.storage)
		} else {
			locally, err = image.LoadLocally(locator, p.storage)
		}
		if err != nil {
			return nil, errors.Wrap(err, "registry.push")
		}
		p.log(KindInfo, "Using locally image", locally.Locator().String(), nil)

		r, err = locally.Open()
		if err != nil {
			return nil, errors.Wrap(err, "registry.push")
		}
		defer func() { _ = r.Close() }()

		ll := locally.Locator()
		remotely, err = p.doPush(locator.Name(), locator.Tag(), ll.Hash(), r)
	} else {
		r, err = os.Open(p.localFile)
		if err != nil {
			return nil, errors.Wrap(err, "registry.push")
		}
		defer func() { _ = r.Close() }()
		p.log(KindInfo, "Using local file", p.localFile, nil)

		tag := locator.Tag()
		if len(tag) == 0 {
			tag = _DefaultTagName
			p.log(KindInfo, "Using default tag", tag, nil)
		}

		var hash string
		hash, err = p.GenHash(r)
		if err != nil {
			return nil, errors.Wrap(err, "registry.push")
		}
		p.log(KindInfo, "Calculating image hash", hash, nil)

		remotely, err = p.doPush(locator.Name(), tag, hash, r)
	}

	if err != nil {
		return nil, err
	}

	if err = remotely.OverwriteDefaults(p.ctx, locator); err != nil {
		return nil, err
	}

	return remotely, nil
}

// doPush 执行上传逻辑
func (p *Pusher) doPush(name, tag, hash string, rs io.ReadSeeker) (*image.Remotely, error) {
	fullyLocator := image.Locate(name, image.WithLocateTag(tag), image.WithLocateHash(hash))
	remotely, err := image.LoadRemotely(p.ctx, fullyLocator, p.oss)
	if err == nil && !p.force { // returns remotely directly when image exists
		return remotely, nil
	}
	if err != nil && err != image.ErrRemotelyImageNotFound {
		return nil, errors.Wrap(err, "registry.push")
	}

	remotely, err = image.NewRemotely(p.ctx, fullyLocator, p.oss)
	if err != nil {
		return nil, errors.Wrap(err, "registry.push")
	}

	lock := remotely.WriterLocker()
	locked, err := lock.Lock(locker.FastFail)
	if err != nil {
		return nil, errors.Wrap(err, "registry.push")
	}
	if !locked {
		return nil, ErrPushInProgress
	}
	defer func() { _ = lock.Unlock() }()

	if err = remotely.Verify(p.ctx); err != nil && err != image.ErrRemotelyImageNotFound {
		return nil, errors.Wrap(err, "registry.push")
	} else if err == nil && !p.force {
		return remotely, nil
	}

	var size int64
	if size, err = xio.SeekerLen(rs); err != nil {
		return nil, errors.Wrap(err, "registry.push")
	}

	var rd io.Reader = rs

	// progressbar
	if len(p.pbOpts) != 0 {
		bar := progressbar.NewOptions64(size, p.pbOpts...)
		defer func() { _ = bar.Close() }()

		prd := progressbar.NewReader(rd, bar)
		rd = &prd
	}

	// progress subscribe
	if p.aph != nil {
		ap := NewAtomicProgress(size, rd)
		p.aph(ap)
		rd = ap
	}

	if err = remotely.WriteFrom(p.ctx, rd); err != nil {
		return nil, errors.Wrap(err, "registry.push")
	}

	if err = remotely.Occupy(p.ctx); err != nil {
		return nil, errors.Wrap(err, "registry.push")
	}

	if err = remotely.Verify(p.ctx); err != nil {
		return nil, errors.Wrap(err, "registry.push")
	}

	return remotely, nil
}

// GenHash 为文件生成哈希值
func (p *Pusher) GenHash(r io.ReadSeeker) (string, error) {
	_, err := r.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}
	defer func() { _, _ = r.Seek(0, io.SeekStart) }()

	h := md5.New()
	if _, err = io.Copy(h, r); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
