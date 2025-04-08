package registry

import (
	"context"
	_errors "errors"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/locker"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xhttp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xio"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xoss"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry/image"
	"io"
	"time"
)

var (
	// PullFastFail 快速失败拉取
	PullFastFail = locker.FastFail
	// PullUntilComplete 直到完成
	PullUntilComplete = locker.UntilAcquire

	// ErrPullInProgress 表示当前还有其他的进程正在拉取镜像
	ErrPullInProgress = _errors.New("registry.pull: another still in progress")
	// ErrPullLockTimeout 表示获取锁超时了
	ErrPullLockTimeout = _errors.New("registry.pull: timeout")
)

// PullOption 用于配置拉取镜像的参数
type PullOption func(p *Puller) error

// WithPullTimeout 为拉取设置超时时间
func WithPullTimeout(timeout time.Duration) PullOption {
	return func(p *Puller) error {
		p.timeout = timeout
		return nil
	}
}

// WithPullBlocking 一直阻塞直到拉取完成或者发生错误
func WithPullBlocking() PullOption {
	return func(p *Puller) error {
		p.timeout = PullUntilComplete
		return nil
	}
}

// WithPullFastFail 快速失败的拉取镜像
func WithPullFastFail() PullOption {
	return func(p *Puller) error {
		p.timeout = PullFastFail
		return nil
	}
}

// WithPullContext 为拉取镜像指定上下文
func WithPullContext(ctx context.Context) PullOption {
	return func(p *Puller) error {
		p.ctx = ctx
		return nil
	}
}

// WithPullLogger 为拉取镜像指定日志记录器
func WithPullLogger(log Logger) PullOption {
	return func(p *Puller) error {
		p.log = log
		return nil
	}
}

// WithPullProgressBar 为拉取镜像创建进度条
func WithPullProgressBar(opts []progressbar.Option) PullOption {
	return func(p *Puller) error {
		p.pbOpts = opts
		return nil
	}
}

// WithPullProgressSubscribe 订阅拉取镜像的进度
func WithPullProgressSubscribe(aph func(*AtomicProgress)) PullOption {
	return func(p *Puller) error {
		p.aph = aph
		return nil
	}
}

// WithPullForce 配置是否强制拉取镜像
func WithPullForce(force bool) PullOption {
	return func(p *Puller) error {
		p.force = force
		return nil
	}
}

// Pull 从远程拉取指定的镜像文件
func (c *client) Pull(ctx context.Context, locator image.Locator, options ...PullOption) (*image.Locally, error) {
	puller := &Puller{
		remote:  c.remote,
		storage: c.storage,
		oss:     c.oss,
		ctx:     ctx,
		log:     NewDiscardLogger(),
	}

	for _, option := range options {
		if err := option(puller); err != nil {
			return nil, errors.Wrap(err, "registry.pull")
		}
	}

	return puller.Pull(locator)
}

// Puller 是一个特定的镜像拉取工具
type Puller struct {
	timeout time.Duration
	remote  *xhttp.Client
	oss     *xoss.ObjectStorageService
	storage *fsutil.Filesystem

	force bool

	ctx    context.Context
	log    Logger
	aph    func(*AtomicProgress)
	pbOpts []progressbar.Option
}

// Pull 根据配置的参数执行拉取镜像
func (p *Puller) Pull(locator image.Locator) (*image.Locally, error) {
	remotely, err := image.LoadRemotely(p.ctx, locator, p.oss)
	if err != nil {
		return nil, errors.Wrap(err, "registry.pull")
	}

	if len(locator.Tag()) == 0 {
		p.log(KindInfo, "Using default tag", remotely.Locator().Tag(), nil)
	}

	fullyLocator := image.Locate(remotely.Name(), image.WithLocateTag(remotely.Tag()),
		image.WithLocateHash(remotely.Hash()))

	locally, err := image.NewLocally(fullyLocator, p.storage)
	if err != nil {
		return nil, errors.Wrap(err, "registry.pull")
	}

	if exists, err := locally.Exists(); err != nil {
		return nil, errors.Wrap(err, "registry.pull")
	} else if exists && !p.force {
		return locally, locally.RewriteDefaults(locator)
	}

	lock := locally.WriterLocker()
	locked, err := lock.Lock(p.timeout)
	if err != nil {
		if err == locker.ErrLockTimeout {
			return nil, ErrPullLockTimeout
		}
		return nil, errors.Wrap(err, "registry.pull")
	}
	if !locked {
		return locally, ErrPullInProgress
	}
	defer func() { _ = lock.Unlock() }()

	// 在获取到锁之后二次确认是否需要拉取
	if exists, err := locally.Exists(); err != nil {
		return nil, errors.Wrap(err, "registry.pull")
	} else if exists && !p.force {
		return locally, locally.RewriteDefaults(locator)
	}

	p.log(KindDebug, "Pulling from", fullyLocator.String(), nil)
	w, err := locally.Writer()
	if err != nil {
		return nil, errors.Wrap(err, "registry.pull")
	}
	defer func() { _ = w.Close() }()

	rr, sz, err := remotely.Reader(p.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "registry.pull")
	}
	defer func() { _ = rr.Close() }()
	var r io.Reader = rr

	// progressbar
	if len(p.pbOpts) != 0 {
		bar := progressbar.NewOptions64(sz, p.pbOpts...)
		defer func() { _ = bar.Close() }()

		prd := progressbar.NewReader(r, bar)
		r = &prd
	}

	// progress subscriber
	if p.aph != nil {
		ap := NewAtomicProgress(sz, r)
		p.aph(ap)
		r = ap
	}

	if _, err = xio.Copy(p.ctx, w, r); err != nil {
		return nil, errors.Wrap(err, "registry.pull")
	}

	if err = locally.Occupy(); err != nil {
		return nil, errors.Wrap(err, "registry.pull")
	}

	if err = locally.RewriteDefaults(locator); err != nil {
		return nil, errors.Wrap(err, "registry.pull")
	}

	return locally, nil
}
