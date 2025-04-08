package iam

import (
	"fmt"

	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	iam_client "github.com/yuansuan/ticp/common/project-root-iam/iam-client"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
)

type Client struct {
	base *iam_client.IamClient
}

func NewClient(cfg config.OpenAPI) (*Client, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}

	if cfg.AccessKeyId == "" {
		return nil, fmt.Errorf("access_key_id cannot be empty")
	}

	if cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("access_key_secret cannot be empty")
	}

	return &Client{
		base: iam_client.NewClient(cfg.Endpoint, cfg.AccessKeyId, cfg.AccessKeySecret),
	}, nil
}

func (c *Client) IsYsProductUser(userId snowflake.ID) (bool, error) {
	resp, err := c.base.IsYSProductAccount(&iam_api.IsYSProductAccountRequest{
		UserId: userId.String(),
	})
	if err != nil {
		return false, fmt.Errorf("call iam server IsYSProductAccount failed, %w", err)
	}

	return resp.IsYSProductAccount, nil
}

func (c *Client) AssumeRole(assumeUserId snowflake.ID) (*iam_api.Credentials, error) {
	assumeResp, err := c.base.AssumeRoleDefault(assumeUserId.String(), "")
	if err != nil {
		return nil, fmt.Errorf("call assume role api failed, %w", err)
	}

	if assumeResp == nil || assumeResp.Credentials == nil {
		return nil, fmt.Errorf("invalid assume role response")
	}

	if assumeResp.Credentials.AccessKeyId == "" {
		return nil, fmt.Errorf("accessKeyID empty after assuming role")
	}

	if assumeResp.Credentials.AccessKeySecret == "" {
		return nil, fmt.Errorf("accessKeySecret empty after assuming role")
	}

	return assumeResp.Credentials, nil
}
