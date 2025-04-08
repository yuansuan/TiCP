package image

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xurl"
)

// Metadata 镜像元数据
type Metadata struct {
	Created *time.Time
}

// Cellar 用于存放镜像的仓库
type Cellar interface {
	// FindDirectories 列出所有的文件夹
	FindDirectories(ctx context.Context, prefix string) ([]string, error)

	// FindObjects 列出所有的文件
	FindObjects(ctx context.Context, prefix string) ([]string, error)

	// ReadAsString 打开并返回其内容
	ReadAsString(ctx context.Context, name string) (string, error)

	// StatObject 获取对象的信息
	StatObject(ctx context.Context, key string) (*Metadata, error)
}

// FindOptions 查询镜像的参数配置
type FindOptions struct {
	Defaults bool
}

// FindLocators 查找符合标准下的
func FindLocators(ctx context.Context, cellar Cellar, pattern string, opts FindOptions) ([]*DefaultedLocator, error) {
	dirs, err := cellar.FindDirectories(ctx, "/")
	if err != nil {
		return nil, err
	}

	var locators []*DefaultedLocator
	for _, name := range dirs {
		if len(pattern) == 0 || strings.Contains(name, pattern) {
			if opts.Defaults {
				l := &DefaultedLocator{name: name}
				l.tag, _ = cellar.ReadAsString(ctx, fsutil.Join(name, _LockFile))
				if len(l.tag) != 0 {
					l.defaultedTag = true
					l.hash, _ = cellar.ReadAsString(ctx, fsutil.Join(name, l.tag, _LockFile))
					l.defaultedHash = len(l.hash) != 0
					if len(l.tag) != 0 && len(l.hash) != 0 {
						if md, err := cellar.StatObject(ctx, fsutil.Join(name, l.tag, l.hash)); err == nil {
							l.created = md.Created
						}
						locators = append(locators, l)
					}
				}
				continue
			}

			tags, err := cellar.FindDirectories(ctx, name+"/")
			if err != nil {
				return nil, err
			}

			dfTag, _ := cellar.ReadAsString(ctx, xurl.Join(name, _LockFile))
			for _, tag := range tags {
				dfHash, _ := cellar.ReadAsString(ctx, xurl.Join(name, tag, _LockFile))
				hashes, err := cellar.FindObjects(ctx, xurl.Join(name, tag)+"/")
				if err != nil {
					return nil, err
				}

				for _, hash := range hashes {
					if hash != _LockFile {
						l := &DefaultedLocator{name: name, tag: tag, hash: hash}
						l.defaultedTag = tag == dfTag
						l.defaultedHash = hash == dfHash
						locators = append(locators, l)

						if md, err := cellar.StatObject(ctx, fsutil.Join(name, l.tag, l.hash)); err == nil {
							l.created = md.Created
						}
					}
				}
			}
		}
	}
	return locators, nil
}

type DefaultedLocator struct {
	name    string
	tag     string
	hash    string
	created *time.Time

	defaultedTag  bool
	defaultedHash bool
}

// String 返回字符串形式标识
func (l *DefaultedLocator) String() string {
	return fmt.Sprintf("%s:%s@%s", l.name, l.tag, l.hash)
}

// ShortString 返回字符串形式的短标识
func (l *DefaultedLocator) ShortString() string {
	return fmt.Sprintf("%s:%s", l.name, l.tag)
}

// Name 返回需要定位的镜像的名称
func (l *DefaultedLocator) Name() string {
	return l.name
}

// Tag 返回需要定位的镜像的具体标签
func (l *DefaultedLocator) Tag() string {
	return l.tag
}

// Hash 返回需要定位的镜像的具体哈希值
func (l *DefaultedLocator) Hash() string {
	return l.hash
}

// Created 返回镜像的创建时间
func (l *DefaultedLocator) Created() *time.Time {
	return l.created
}

// IsDefaultTag 是否是默认的标签
func (l *DefaultedLocator) IsDefaultTag() bool {
	return l.defaultedTag
}

// IsDefaultedHash 是否是默认的哈希值
func (l *DefaultedLocator) IsDefaultedHash() bool {
	return l.defaultedHash
}
