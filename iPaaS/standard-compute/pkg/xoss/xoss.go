package xoss

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/sbs"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xurl"
)

const (
	// _ViewSizeLimited 查看对象内容的限制大小
	_ViewSizeLimited = 1024 * 1024 // 1Mi
)

var (
	// ErrObjectTooLarge 文件太大不能载入内存
	ErrObjectTooLarge = errors.New("oss: object too large")
)

// ObjectStorageService 基于AWS-S3兼容的对象存储服务
type ObjectStorageService struct {
	s3 *s3.S3

	prefix string
	bucket *string
}

// PutObject 上传文件
func (oss *ObjectStorageService) PutObject(ctx context.Context, key string, r io.ReadSeeker) (*s3.PutObjectOutput, error) {
	return oss.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: oss.bucket,
		Key:    oss.key(key),
		Body:   r,
	})
}

// UploadOption 上传选项配置
type UploadOption = func(u *s3manager.Uploader)

// UploadObject 上传文件
func (oss *ObjectStorageService) UploadObject(ctx context.Context, key string, r io.Reader, options ...UploadOption) (*s3manager.UploadOutput, error) {
	u := s3manager.NewUploaderWithClient(oss.s3, options...)
	return u.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: oss.bucket,
		Key:    oss.key(key),
		Body:   r,
	})
}

// CopyObject 复制对象
func (oss *ObjectStorageService) CopyObject(ctx context.Context, from, to string) (*s3.CopyObjectOutput, error) {
	return oss.s3.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Bucket:     oss.bucket,
		Key:        oss.key(to),
		CopySource: aws.String(xurl.Join(*oss.bucket, *oss.key(from))),
	})
}

// DeleteObject 删除对象
func (oss *ObjectStorageService) DeleteObject(ctx context.Context, key string) (*s3.DeleteObjectOutput, error) {
	return oss.s3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: oss.bucket,
		Key:    oss.key(key),
	})
}

// HeadObject 获取对象的元数据
func (oss *ObjectStorageService) HeadObject(ctx context.Context, key string) (*s3.HeadObjectOutput, error) {
	return oss.s3.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: oss.bucket,
		Key:    oss.key(key),
	})
}

// GetObject 获取一个对象
func (oss *ObjectStorageService) GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	return oss.s3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: oss.bucket,
		Key:    oss.key(key),
	})
}

// GetAsString 获取一个对象的内容并返回
func (oss *ObjectStorageService) GetAsString(ctx context.Context, key string) (string, error) {
	resp, err := oss.GetObject(ctx, key)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if *resp.ContentLength > _ViewSizeLimited {
		return "", ErrObjectTooLarge
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return sbs.String(data), nil
}

// ListObjects 列出当前目录下的所有文件
func (oss *ObjectStorageService) ListObjects(ctx context.Context, prefix string) (*s3.ListObjectsOutput, error) {
	var k = oss.key(prefix)
	if strings.HasSuffix(prefix, "/") {
		s := *k + "/"
		k = &s
	}

	return oss.s3.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket:    oss.bucket,
		Delimiter: aws.String("/"),
		Prefix:    k,
	})
}

// ListDirectories 列出所有目录
func (oss *ObjectStorageService) ListDirectories(ctx context.Context, prefix string) ([]string, error) {
	resp, err := oss.ListObjects(ctx, prefix)
	if err != nil {
		return nil, err
	}

	var directories []string
	for _, item := range resp.CommonPrefixes {
		dir := xurl.Trim((*item.Prefix)[len(*resp.Prefix):])
		if len(dir) != 0 {
			directories = append(directories, dir)
		}
	}

	return directories, nil
}

// ListFiles 列出所有的文件
func (oss *ObjectStorageService) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	resp, err := oss.ListObjects(ctx, prefix)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, item := range resp.Contents {
		filename := xurl.Trim((*item.Key)[len(*resp.Prefix):])
		if len(filename) != 0 {
			files = append(files, filename)
		}
	}

	return files, nil
}

// key 返回带前缀的对象名称
func (oss *ObjectStorageService) key(k string) *string {
	s := xurl.Trim(xurl.Join(oss.prefix, k))
	return &s
}

// New 创建基于OSS的文件服务
func New(cfg *config.ObjectStorageService) (*ObjectStorageService, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Endpoint:    aws.String(cfg.Endpoint),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.AccessSecret, ""),
	})

	if err != nil {
		return nil, err
	}

	return &ObjectStorageService{s3: s3.New(sess), bucket: &cfg.Bucket, prefix: cfg.PathPrefix}, nil
}

// IsObjectNotExists 检查错误是否是对象不存在
func IsObjectNotExists(err error) bool {
	if err != nil {
		if ae, ok := err.(awserr.Error); ok {
			if ae.Code() == s3.ErrCodeNoSuchKey {
				return true
			}
		}
		if ae, ok := err.(awserr.RequestFailure); ok {
			if ae.StatusCode() == http.StatusNotFound {
				return true
			}
		}
	}
	return false
}
