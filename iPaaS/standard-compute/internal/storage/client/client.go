package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	openapi "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	iamclient "github.com/yuansuan/ticp/common/project-root-iam/iam-client"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/oshelp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/storage/client/assumerole"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/util"
)

const (
	defaultRetryTimes    = 5
	defaultRetryInterval = 10 * time.Second
)

type Client struct {
	cfg               *config.Config
	openapiClientPool *openapiClientPool
	iamClient         *iamclient.IamClient
	assumeRoleCache   *assumerole.Cache
}

func New(cfg *config.Config) *Client {
	c := &Client{
		cfg:               cfg,
		openapiClientPool: newOpenapiClientPool(cfg),
		assumeRoleCache:   assumerole.NewCache(),
	}

	// 此client用于角色扮演，以hpc角色，获取用户角色的ak pair，从而使用用户的个人存储
	c.iamClient = iamclient.NewClient(cfg.Iam.Endpoint, cfg.Iam.AppKey, cfg.Iam.AppSecret)
	if cfg.Iam.Proxy != "" {
		c.iamClient.SetProxy(cfg.Iam.Proxy)
	}

	return c
}

type openapiClientPool struct {
	pool          *sync.Pool //hpcPool
	poolWithProxy *sync.Pool //cloudPool
	cfg           *config.Config
}

// 一个叫pool, 一个叫poolWithProxy是不是更合适点？

func newOpenapiClientPool(cfg *config.Config) *openapiClientPool {
	cp := &openapiClientPool{
		pool:          &sync.Pool{},
		poolWithProxy: &sync.Pool{},
		cfg:           cfg,
	}

	cp.pool.New = func() interface{} {
		return createClient(cfg, false) // 不使用代理
	}

	cp.poolWithProxy.New = func() interface{} {
		return createClient(cfg, true) // 使用代理
	}
	return cp
}

func createClient(cfg *config.Config, useProxy bool) *openapi.Client {
	clientOpts := ensureClientOpts(cfg.OpenAPI, useProxy)
	openapiCli, err := openapi.NewClient(credential.NewCredential("", ""), clientOpts...)
	if err != nil {
		log.Errorf("new client failed, %v", err)
		return nil
	}
	return openapiCli
}

func storageRetryCondition(resp *http.Response, err error) bool {
	if err != nil && !util.IsFileMissingCausedError(err) {
		log.Warn(err)
		return true
	}

	if resp != nil && resp.StatusCode >= http.StatusInternalServerError {
		log.Warnf("StatusCode: %d", resp.StatusCode)
		return true
	}

	return false
}

func ensureClientOpts(cfg config.OpenAPIConfig, useProxy bool) []openapi.Option {
	clientOpts := make([]openapi.Option, 0)

	maxRetryTimes := defaultRetryTimes
	if cfg.MaxRetryTimes != 0 {
		maxRetryTimes = cfg.MaxRetryTimes
	}
	clientOpts = append(clientOpts, openapi.WithRetryTimes(maxRetryTimes))

	retryInterval := defaultRetryInterval
	if cfg.RetryInterval != 0 {
		retryInterval = cfg.RetryInterval
	}
	clientOpts = append(clientOpts, openapi.WithRetryInterval(retryInterval))

	if useProxy && cfg.Proxy != "" {
		clientOpts = append(clientOpts, openapi.WithProxy(cfg.Proxy))
	}

	clientOpts = append(clientOpts, openapi.WithRetryCondition(storageRetryCondition))

	return clientOpts
}

func (cp *openapiClientPool) Get(endpoint, appKey, appSecret string) (*openapi.Client, error) {
	isHpcStorage := strings.HasPrefix(endpoint, cp.cfg.HpcStorageAddress)
	var cli interface{}
	if isHpcStorage {
		cli = cp.pool.Get()
		log.Debugf("Using HPC Storage client, No Proxy, Endpoint: %s", endpoint)
	} else {
		cli = cp.poolWithProxy.Get()
		log.Debugf("Using Cloud Storage client, With Proxy, Endpoint: %s", endpoint)
	}
	openapiClient, ok := cli.(*openapi.Client)
	if !ok {
		return nil, fmt.Errorf("get from pool cannot convert to *openapi.Client")
	}
	// 动态设置baseUrl/credential
	openapiClient.SetBaseUrl(endpoint)
	if err := openapiClient.SetCredential(credential.NewCredential(appKey, appSecret)); err != nil {
		return nil, fmt.Errorf("openapi client set credentail failed, %w", err)
	}

	return openapiClient, nil
}

func (cp *openapiClientPool) Put(openapiClient *openapi.Client, isHpcStorage bool) {
	if openapiClient == nil {
		return
	}

	// 丢回pool时重新至空
	openapiClient.SetBaseUrl("")
	_ = openapiClient.SetCredential(credential.NewCredential("", ""))

	if isHpcStorage {
		cp.pool.Put(openapiClient)
	} else {
		cp.poolWithProxy.Put(openapiClient)
	}
}

type FileStat struct {
	name    string
	size    int64
	mode    uint32
	modTime time.Time
	isDir   bool
}

func (fs *FileStat) Name() string {
	return fs.name
}

func (fs *FileStat) Size() int64 {
	return fs.size
}

func (fs *FileStat) Mode() os.FileMode {
	return os.FileMode(fs.mode)
}

func (fs *FileStat) ModTime() time.Time {
	return fs.modTime
}

func (fs *FileStat) IsDir() bool {
	return fs.isDir
}

func (fs *FileStat) Sys() interface{} {
	return nil
}

func (fs *FileStat) String() string {
	return fmt.Sprintf("[name: %s, size: %d, mode: %s, modTime: %s, isDir: %v]",
		fs.Name(),
		fs.Size(),
		fs.Mode(),
		fs.ModTime(),
		fs.IsDir())
}

func (c *Client) Stat(endpoint, path string) (os.FileInfo, error) {
	log.Debugf("call storage stat api, endpoint: %s, path: %s", endpoint, path)
	userId, err := parseUserIdFromPath(path)
	if err != nil {
		return nil, fmt.Errorf("parse userId from path failed, %w", err)
	}

	cli, err := c.getOpenAPIClientAfterAssumeRole(endpoint, userId)
	if err != nil {
		return nil, fmt.Errorf("get openapi client after assume role failed, %w", err)
	}
	defer func() {
		isHpcStorage := strings.HasPrefix(endpoint, c.cfg.HpcStorageAddress)
		c.openapiClientPool.Put(cli, isHpcStorage)
	}()

	resp, err := cli.Storage.Stat(cli.Storage.Stat.Path(path))
	if err != nil {
		if resp.ErrorCode == "PathNotFound" {
			return nil, os.ErrNotExist
		}
		return nil, fmt.Errorf("call api stat failed, %w", err)
	}
	if resp.Data == nil || resp.Data.File == nil {
		return nil, fmt.Errorf("resp.Data or resp.Data.File is nil")
	}

	fs := &FileStat{
		name:    resp.Data.File.Name,
		size:    resp.Data.File.Size,
		mode:    resp.Data.File.Mode,
		modTime: resp.Data.File.ModTime,
		isDir:   resp.Data.File.IsDir,
	}

	log.Debugf("call storage stat api success, fileStat: %s", fs.String())
	return fs, nil
}

func (c *Client) RealPath(endpoint, path string) (string, error) {
	log.Debugf("call storage real path api, endpoint %s, path: %s", endpoint, path)
	cli, err := c.openapiClientPool.Get(endpoint, c.cfg.Iam.AppKey, c.cfg.Iam.AppSecret)
	if err != nil {
		return "", fmt.Errorf("get openapi client from pool failed, %w", err)
	}
	defer func() {
		isHpcStorage := strings.HasPrefix(endpoint, c.cfg.HpcStorageAddress)
		c.openapiClientPool.Put(cli, isHpcStorage)
	}()

	resp, err := cli.Storage.Realpath(cli.Storage.Realpath.RelativePath(path))
	if err != nil {
		return "", fmt.Errorf("call real path api failed, %w", err)
	}
	if resp.Data == nil {
		return "", fmt.Errorf("resp.Data is nil")
	}

	return resp.Data.RealPath, nil
}

func (c *Client) Download(endpoint, path string, beginOffset int64, endOffset int64) ([]byte, error) {
	log.Debugf("call storage download api, endpoint %s, path: %s, beginOffset: %d, endOffset: %d", endpoint, path, beginOffset, endOffset)
	userId, err := parseUserIdFromPath(path)
	if err != nil {
		return nil, fmt.Errorf("parse userId from path failed, %w", err)
	}

	cli, err := c.getOpenAPIClientAfterAssumeRole(endpoint, userId)
	if err != nil {
		return nil, fmt.Errorf("get openapi client after assume role failed, %w", err)
	}
	defer func() {
		isHpcStorage := strings.HasPrefix(endpoint, c.cfg.HpcStorageAddress)
		c.openapiClientPool.Put(cli, isHpcStorage)
	}()

	resp, err := cli.Storage.Download(
		cli.Storage.Download.Path(path),
		cli.Storage.Download.Range(beginOffset, endOffset),
	)
	if err != nil {
		return nil, fmt.Errorf("call download api failed, %w", err)
	}

	return resp.Data, nil
}

func (c *Client) DownloadByStream(ctx context.Context, endpoint, path string, dst *os.File) error {
	log.Debugf("call storage download api, endpoint %s, path: %s", endpoint, path)
	userId, err := parseUserIdFromPath(path)
	if err != nil {
		return fmt.Errorf("parse userId from path failed, %w", err)
	}

	cli, err := c.getOpenAPIClientAfterAssumeRole(endpoint, userId)
	if err != nil {
		return fmt.Errorf("get openapi client after assume role failed, %w", err)
	}
	defer func() {
		isHpcStorage := strings.HasPrefix(endpoint, c.cfg.HpcStorageAddress)
		c.openapiClientPool.Put(cli, isHpcStorage)
	}()
	_, err = cli.Storage.Download(
		cli.Storage.Download.Path(path),
		cli.Storage.Download.WithResolver(func(resp *http.Response) error {
			defer resp.Body.Close()

			username := config.GetConfig().BackendProvider.SchedulerCommon.SubmitSysUser
			opts := make([]oshelp.Option, 0)
			if username != "" {
				opts = append(opts, oshelp.WithChown(username))
			}

			err = oshelp.CopyToFile(ctx, dst, resp.Body)
			return err
		}),
	)
	if err != nil {
		return fmt.Errorf("call api download failed, %w", err)
	}

	return nil
}

func (c *Client) Ls(endpoint string, path string) ([]os.FileInfo, error) {
	log.Debugf("call storage ls api, endpoint: %s, path: %s", endpoint, path)
	userId, err := parseUserIdFromPath(path)
	if err != nil {
		return nil, fmt.Errorf("parse userId from path failed, %w", err)
	}

	cli, err := c.getOpenAPIClientAfterAssumeRole(endpoint, userId)
	if err != nil {
		return nil, fmt.Errorf("get openapi client after assume role failed, %w", err)
	}
	defer func() {
		isHpcStorage := strings.HasPrefix(endpoint, c.cfg.HpcStorageAddress)
		c.openapiClientPool.Put(cli, isHpcStorage)
	}()

	filesInfo, err := c.lsTotal(cli, path)
	if err != nil {
		return nil, fmt.Errorf("ls total failed, %w", err)
	}

	log.Debugf("call storage ls api success, filesInfo: %s", filesInfo)
	return filesInfo, nil
}

const (
	defaultPageSize    = 1000
	defaultMaxCapacity = 200
	pageEndNextMarker  = -1
)

// lsWithPage 200次连续调用，pageSize = 1000，如果还不结束，认为此处有异常，报错抛出
func (c *Client) lsTotal(cli *openapi.Client, path string) ([]os.FileInfo, error) {
	filesInfo := make([]os.FileInfo, 0)

	var offset int64 = 0
	capacity := 0
	for {
		if capacity >= defaultMaxCapacity {
			return nil, fmt.Errorf("already call more than %d times lsWithPage api with pageSize = %d, should check it manually", defaultMaxCapacity, defaultPageSize)
		}
		capacity++

		resp, err := cli.Storage.LsWithPage(
			cli.Storage.LsWithPage.Path(path),
			cli.Storage.LsWithPage.PageSize(defaultPageSize),
			cli.Storage.LsWithPage.PageOffset(offset),
		)
		if err != nil {
			return nil, fmt.Errorf("call storage lsWithPage api failed, %w", err)
		}
		if resp.Data == nil {
			return nil, fmt.Errorf("response.Data is nil")
		}

		for _, fs := range resp.Data.Files {
			fileStat := &FileStat{
				name:    fs.Name,
				size:    fs.Size,
				mode:    fs.Mode,
				modTime: fs.ModTime,
				isDir:   fs.IsDir,
			}

			filesInfo = append(filesInfo, fileStat)
		}

		if resp.Data.NextMarker == pageEndNextMarker {
			break
		}

		offset = resp.Data.NextMarker
	}

	return filesInfo, nil
}

const (
	pathSeparator = "/"
)

// path format like: /{userId}/xxx
func parseUserIdFromPath(path string) (string, error) {
	fields := strings.Split(path, pathSeparator)
	// at least /{userId}/
	if len(fields) < 3 {
		return "", fmt.Errorf("invalid path [%s], not contain userId", path)
	}

	if fields[1] == "" {
		return "", fmt.Errorf("invalid path [%s], userId is empty", path)
	}

	return fields[1], nil
}

const (
	refreshAssumeRoleIfLessThanTime = 5 * time.Minute
)

// get openapi client for storage
func (c *Client) getOpenAPIClientAfterAssumeRole(openapiEndpoint, userId string) (*openapi.Client, error) {
	if userId == "" {
		return nil, errors.New("userId is empty")
	}
	var err error

	// first get from cache
	assumeRoleValue, exist := c.assumeRoleCache.Get(userId)
	if (!exist) || (exist && time.Now().Add(refreshAssumeRoleIfLessThanTime).After(assumeRoleValue.ExpiredTime)) {
		// insert/refresh assumeRoleValue to cache
		assumeRoleValue, err = c.assumeRole(userId)
		if err != nil {
			return nil, fmt.Errorf("assume role failed, %w", err)
		}
		c.assumeRoleCache.Set(userId, assumeRoleValue)
	}

	openapiClient, err := c.openapiClientPool.Get(openapiEndpoint,
		assumeRoleValue.AccessKeyId,
		assumeRoleValue.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("get openapi client from pool failed, %w", err)
	}

	return openapiClient, nil
}

func (c *Client) assumeRole(userId string) (assumerole.Value, error) {
	value := assumerole.Value{}
	// set roleName empty for now
	assumeRoleResp, err := c.iamClient.AssumeRoleDefault(userId, "")
	if err != nil {
		return value, fmt.Errorf("call assume role api failed, %w", err)
	}
	if assumeRoleResp.Credentials == nil {
		return value, errors.New("assumeRoleResp.Credentials is nil")
	}

	value.AccessKeyId = assumeRoleResp.Credentials.AccessKeyId
	value.AccessKeySecret = assumeRoleResp.Credentials.AccessKeySecret
	value.Token = assumeRoleResp.Credentials.SessionToken
	value.ExpiredTime = assumeRoleResp.ExpireTime

	return value, nil
}
