package account

import (
	"fmt"
	"sync"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	openapiWrap "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/openapi"

	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	openapi "github.com/yuansuan/ticp/common/openapi-go"
)

var once = &sync.Once{}
var _client *client

type client struct {
	openapiClient *openapi.Client
}

func newClient() (*client, error) {
	openapiClient, err := openapiWrap.NewClient()
	if err != nil {
		return nil, err
	}
	c := &client{
		openapiClient: openapiClient,
	}
	return c, nil
}

func Client() *client {
	once.Do(func() {
		cli, err := newClient()
		if err != nil {
			panic(fmt.Sprintf("new merchandise client failed, %v", err))
		}
		_client = cli
	})
	return _client
}

func (c *client) GetAccountByUserId(userId snowflake.ID) (*v20230530.AccountDetail, error) {
	if userId == snowflake.ID(0) {
		return nil, common.ErrInvalidUserID
	}
	resp, err := c.openapiClient.Account.ByYsIDGet(
		c.openapiClient.Account.ByYsIDGet.UserID(userId.String()),
	)
	if err != nil {
		return nil, fmt.Errorf("call account byYsIDGet api failed, %w", err)
	}
	if resp == nil || resp.Data == nil || resp.Data.AccountDetail == nil {
		return nil, common.ErrAccountResponse
	}
	return resp.Data.AccountDetail, nil
}
