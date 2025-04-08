package iam

import (
	"errors"
	"fmt"
	"sync"

	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	iam_client "github.com/yuansuan/ticp/common/project-root-iam/iam-client"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/assumerole"
)

var once = &sync.Once{}
var _client *client

type client struct {
	base *iam_client.IamClient
}

func newClient() (*client, error) {
	conf := config.GetConfig()

	if conf.AK == "" {
		return nil, fmt.Errorf("openapi.access_key_id cannot be empty")
	}

	if conf.AS == "" {
		return nil, fmt.Errorf("openapi.access_key_secret cannot be empty")
	}

	if conf.OpenAPIEndpoint == "" {
		return nil, fmt.Errorf("openapi.endpoint cannot be empty")
	}

	return &client{
		base: iam_client.NewClient(conf.OpenAPIEndpoint, conf.AK, conf.AS),
	}, nil
}

// Client iam client
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

// IsYsProductUser 是否是远算产品用户
func (c *client) IsYsProductUser(userID snowflake.ID) (bool, error) {
	resp, err := c.base.IsYSProductAccount(&iam_api.IsYSProductAccountRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return false, fmt.Errorf("call iam server IsYSProductAccount failed, %w", err)
	}

	return resp.IsYSProductAccount, nil
}

// AssumeRole 角色扮演
func (c *client) AssumeRole(userID string) (assumerole.Value, error) {
	value := assumerole.Value{}
	// set roleName empty for now
	assumeRoleResp, err := c.base.AssumeRoleDefault(userID, "")
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
