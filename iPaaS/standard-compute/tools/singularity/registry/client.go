package registry

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry/image"

	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil/filemode"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xhttp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xoss"
)

type Client interface {
	Pull(ctx context.Context, locator image.Locator, options ...PullOption) (*image.Locally, error)
	Push(locator image.Locator, options ...PushOption) (*image.Remotely, error)
	Search(pattern string, options ...SearchOption) ([]*image.DefaultedLocator, error)
}

// Client 表示 Singularity 的镜像仓库客户端
type client struct {
	storage *fsutil.Filesystem
	remote  *xhttp.Client
	oss     *xoss.ObjectStorageService
	cfg     *config.Singularity
}

// LoadImageLocally 加载本地镜像描述
func (c *client) LoadImageLocally(locator image.Locator) (*image.Locally, error) {
	l, err := image.LoadLocally(locator, c.storage)
	return l, errors.Wrap(err, "registry.load")
}

// LoadImageRemotely 加载远程镜像描述
func (c *client) LoadImageRemotely(ctx context.Context, locator image.Locator) (*image.Remotely, error) {
	r, err := image.LoadRemotely(ctx, locator, c.oss)
	return r, errors.Wrap(err, "registry.load")
}

// NewClient 创建一个 Singularity 的镜像仓库客户端实例
func NewClient(cfg *config.Singularity) (Client, error) {
	if cfg.IsMock {
		return &mockClient{}, nil
	}

	if err := fsutil.MkdirAuto(cfg.Storage, filemode.Directory); err != nil {
		return nil, errors.Wrap(err, "registry.client")
	}

	c, err := xhttp.New(xhttp.WithBaseUrl(cfg.Registry.BaseURL()))
	if err != nil {
		return nil, errors.Wrap(err, "registry.client")
	}

	fs, err := fsutil.Chroot(cfg.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "registry.client")
	}

	oss, err := xoss.New(cfg.Registry)
	if err != nil {
		return nil, errors.Wrap(err, "registry.client")
	}

	return &client{storage: fs, remote: c, oss: oss, cfg: cfg}, nil
}
