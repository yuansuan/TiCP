package storage

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	openapi "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/ls"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/readAt"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/stat"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/assumerole"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/iam"
)

func DefaultCred(cfg config.CustomT) *credential.Credential {
	return credential.NewCredential(cfg.GetAK(), cfg.GetAS())
}

var once = &sync.Once{}
var _client *client

type client struct {
	assumeRoleCache   *assumerole.Cache
	OpenapiClientPool *openapiClientPool
}

// Client storage client
func Client() *client {
	once.Do(func() {
		_client = New()
	})

	return _client
}

func New() *client {
	c := &client{
		OpenapiClientPool: newOpenapiClientPool(),
		assumeRoleCache:   assumerole.NewCache(),
	}

	return c
}

type openapiClientPool struct {
	pool *sync.Pool
}

func newOpenapiClientPool() *openapiClientPool {
	cp := &openapiClientPool{
		pool: &sync.Pool{},
	}

	cp.pool.New = func() interface{} {
		opts := []openapi.Option{
			openapi.WithBaseURL(openapi.DefaultBaseURL),
			openapi.WithTimeout(module.DefaultTimeout),
			openapi.WithRetryTimes(module.DefaultRetryTimes),
			openapi.WithRetryInterval(module.DefaultRetryInterval),
		}
		openapiClient, err := openapi.NewClient(credential.NewCredential("", ""), opts...)
		if err != nil {
			logging.Default().Errorf("new client failed, %v", err)
			return nil
		}
		return openapiClient
	}
	return cp
}

func (cp *openapiClientPool) Get(endpoint string, timeout time.Duration, cred *credential.Credential) (*openapi.Client, error) {
	cli, ok := cp.pool.Get().(*openapi.Client)
	if !ok {
		return nil, fmt.Errorf("get from pool cannot convert to *openapi.Client")
	}

	// 动态设置baseUrl/timeout/credential
	cli.SetBaseUrl(endpoint)
	cli.SetTimeout(timeout)
	if err := cli.SetCredential(cred); err != nil {
		return nil, fmt.Errorf("openapi client set credentail failed, %w", err)
	}

	return cli, nil
}

func (cp *openapiClientPool) Put(openapiClient *openapi.Client) {
	if openapiClient == nil {
		return
	}

	// 丢回pool时重新至空
	openapiClient.SetBaseUrl("")
	openapiClient.SetTimeout(module.DefaultTimeout)
	_ = openapiClient.SetCredential(credential.NewCredential("", ""))
	cp.pool.Put(openapiClient)
}

const (
	refreshAssumeRoleIfLessThanTime = 5 * time.Minute
)

// get openapi client for storage
func (c *client) GetOpenAPIClientAfterAssumeRole(openapiEndpoint string, timeout time.Duration, userId string) (*openapi.Client, error) {
	if userId == "" {
		return nil, errors.New("userId is empty")
	}
	var err error

	// first get from cache
	assumeRoleValue, exist := c.assumeRoleCache.Get(userId)
	if !exist || time.Now().Add(refreshAssumeRoleIfLessThanTime).After(assumeRoleValue.ExpiredTime) {
		// insert/refresh assumeRoleValue to cache
		assumeRoleValue, err = iam.Client().AssumeRole(userId)
		if err != nil {
			return nil, fmt.Errorf("assume role failed, %w", err)
		}
		c.assumeRoleCache.Set(userId, assumeRoleValue)
	}

	openapiClient, err := c.OpenapiClientPool.Get(openapiEndpoint, timeout,
		credential.NewCredential(assumeRoleValue.AccessKeyId, assumeRoleValue.AccessKeySecret),
	)
	if err != nil {
		return nil, fmt.Errorf("get openapi client from pool failed, %w", err)
	}

	return openapiClient, nil
}

type ClientParams struct {
	Endpoint string
	UserID   string
	Timeout  time.Duration // Default 0 for no timeout
	AdminAPI bool
}

type LsParams struct {
	Offset    int64
	Lspath    string
	RegxpList []string
}

func (c *client) LsWithPage(clientParams ClientParams, lsParams LsParams) (*ls.Response, error) {
	if clientParams.AdminAPI {
		cli, err := c.OpenapiClientPool.Get(clientParams.Endpoint, clientParams.Timeout, DefaultCred(config.GetConfig()))
		if err != nil {
			return nil, err
		}
		defer c.OpenapiClientPool.Put(cli)
		api := cli.Storage.AdminLsWithPage
		return api(
			api.FilterRegexpList(lsParams.RegxpList),
			api.Path(lsParams.Lspath),
			api.PageOffset(lsParams.Offset),
		)
	} else {
		cli, err := c.GetOpenAPIClientAfterAssumeRole(clientParams.Endpoint, clientParams.Timeout, clientParams.UserID)
		if err != nil {
			return nil, err
		}
		defer c.OpenapiClientPool.Put(cli)
		api := cli.Storage.LsWithPage
		return api(
			api.FilterRegexpList(lsParams.RegxpList),
			api.Path(lsParams.Lspath),
			api.PageOffset(lsParams.Offset),
		)
	}
}

func (c *client) Stat(clientParams ClientParams, statpath string) (*stat.Response, error) {
	if clientParams.AdminAPI {
		cli, err := c.OpenapiClientPool.Get(clientParams.Endpoint, clientParams.Timeout, DefaultCred(config.GetConfig()))
		if err != nil {
			return nil, err
		}
		defer c.OpenapiClientPool.Put(cli)
		api := cli.Storage.AdminStat
		return api(
			api.Path(statpath),
		)
	} else {
		cli, err := c.GetOpenAPIClientAfterAssumeRole(clientParams.Endpoint, clientParams.Timeout, clientParams.UserID)
		if err != nil {
			return nil, err
		}
		defer c.OpenapiClientPool.Put(cli)
		api := cli.Storage.Stat
		return api(
			api.Path(statpath),
		)
	}
}

type ReadAtParams struct {
	Readpath string
	Length   int64
	Offset   int64
	Resolver xhttp.ResponseResolver
}

func (c *client) ReadAt(clientParams ClientParams, readAtParams ReadAtParams) (*readAt.Response, error) {
	if clientParams.AdminAPI {
		cli, err := c.OpenapiClientPool.Get(clientParams.Endpoint, clientParams.Timeout, DefaultCred(config.GetConfig()))
		if err != nil {
			return nil, err
		}
		defer c.OpenapiClientPool.Put(cli)
		api := cli.Storage.AdminReadAt
		return api(
			api.Path(readAtParams.Readpath),
			api.Length(readAtParams.Length),
			api.Offset(readAtParams.Offset),
			api.WithResolver(readAtParams.Resolver),
		)
	} else {
		cli, err := c.GetOpenAPIClientAfterAssumeRole(clientParams.Endpoint, clientParams.Timeout, clientParams.UserID)
		if err != nil {
			return nil, err
		}
		defer c.OpenapiClientPool.Put(cli)
		api := cli.Storage.ReadAt
		return api(
			api.Path(readAtParams.Readpath),
			api.Length(readAtParams.Length),
			api.Offset(readAtParams.Offset),
			api.WithResolver(readAtParams.Resolver),
		)
	}
}
