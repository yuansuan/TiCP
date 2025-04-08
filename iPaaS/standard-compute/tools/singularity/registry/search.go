package registry

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry/image"
	"os"

	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xoss"
)

// SearchOption 搜索的参数
type SearchOption func(s *Searcher) error

// WithSearchContext 为搜索配置上下文对象
func WithSearchContext(ctx context.Context) SearchOption {
	return func(s *Searcher) error {
		s.ctx = ctx
		return nil
	}
}

// WithSearchDefaults 是否列出所有的标签
func WithSearchDefaults(defaults bool) SearchOption {
	return func(s *Searcher) error {
		s.defaults = defaults
		return nil
	}
}

// WithSearchLocally 是否搜索本地仓库
func WithSearchLocally(locally bool) SearchOption {
	return func(s *Searcher) error {
		s.locally = locally
		return nil
	}
}

// Search 根据指定的模式搜索远程镜像
func (c *client) Search(pattern string, options ...SearchOption) ([]*image.DefaultedLocator, error) {
	s := &Searcher{oss: c.oss, ctx: context.Background(), storage: c.storage}
	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, "registry.search")
		}
	}

	return s.Search(pattern)
}

// Searcher 远程镜像搜索器
type Searcher struct {
	ctx     context.Context
	oss     *xoss.ObjectStorageService
	storage *fsutil.Filesystem

	defaults bool
	locally  bool
}

// Search 执行搜索操作
func (s *Searcher) Search(pattern string) ([]*image.DefaultedLocator, error) {
	var cellar image.Cellar
	if s.locally {
		cellar = &locallyCellar{storage: s.storage, oss: s.oss}
	} else {
		cellar = &s3Cellar{oss: s.oss}
	}

	return image.FindLocators(s.ctx, cellar, pattern, image.FindOptions{Defaults: s.defaults})
}

type s3Cellar struct {
	oss *xoss.ObjectStorageService
}

func (c *s3Cellar) FindDirectories(ctx context.Context, prefix string) ([]string, error) {
	return c.oss.ListDirectories(ctx, prefix)
}

func (c *s3Cellar) FindObjects(ctx context.Context, prefix string) ([]string, error) {
	return c.oss.ListFiles(ctx, prefix)
}

func (c *s3Cellar) ReadAsString(ctx context.Context, name string) (string, error) {
	return c.oss.GetAsString(ctx, name)
}

func (c *s3Cellar) StatObject(ctx context.Context, key string) (*image.Metadata, error) {
	resp, err := c.oss.HeadObject(ctx, key)
	if err != nil {
		return nil, err
	}
	return &image.Metadata{Created: resp.LastModified}, nil
}

type locallyCellar struct {
	oss     *xoss.ObjectStorageService
	storage *fsutil.Filesystem
}

func (c *locallyCellar) FindDirectories(_ context.Context, prefix string) ([]string, error) {
	return c.readDirAsList(prefix, func(info os.FileInfo) bool {
		return info.IsDir()
	})
}

func (c *locallyCellar) FindObjects(_ context.Context, prefix string) ([]string, error) {
	return c.readDirAsList(prefix, func(info os.FileInfo) bool {
		return !info.IsDir()
	})
}

func (c *locallyCellar) ReadAsString(_ context.Context, name string) (string, error) {
	return c.storage.ReadFileAsString(name)
}

// StatObject 本地的镜像修改时间应该从远程加载
func (c *locallyCellar) StatObject(ctx context.Context, key string) (*image.Metadata, error) {
	resp, err := c.oss.HeadObject(ctx, key)
	if err != nil {
		return nil, err
	}
	return &image.Metadata{Created: resp.LastModified}, nil
}

func (c *locallyCellar) readDirAsList(name string, predicate func(os.FileInfo) bool) ([]string, error) {
	f, err := c.storage.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	files, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var items []string
	for _, item := range files {
		if predicate(item) {
			items = append(items, item.Name())
		}
	}

	return items, nil
}
