package openapi

import (
	"errors"
	"fmt"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"

	baseschema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	openys "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
)

var (
	ErrRequestNil        = errors.New("request is nil")
	ErrResourceIdInvalid = errors.New("invalid resourceId")
	ErrUserIdInvalid     = errors.New("invalid userId")
)

type Client struct {
	base *openys.Client
}

func NewClient(cfg config.OpenAPI) (*Client, error) {
	cli, err := openys.NewClient(
		credential.NewCredential(cfg.AccessKeyId, cfg.AccessKeySecret),
		openys.WithBaseURL(cfg.Endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("new openapi client failed, %w", err)
	}

	c := &Client{
		base: cli,
	}

	return c, nil
}

func (c *Client) GetAccountByUserId(userId snowflake.ID) (*baseschema.AccountDetail, error) {
	if userId == snowflake.ID(0) {
		return nil, ErrUserIdInvalid
	}

	resp, err := c.base.Account.ByYsIDGet(
		c.base.Account.ByYsIDGet.UserID(userId.String()),
	)
	if err != nil {
		return nil, fmt.Errorf("call account byYsIDGet api failed, %w", err)
	}
	if resp == nil || resp.Data == nil || resp.Data.AccountDetail == nil {
		return nil, common.ErrAccountResponse
	}

	return resp.Data.AccountDetail, nil
}

func (c *Client) GetCloudStorageEndpointByZone(zone string) (string, error) {
	zones, err := c.base.Job.ZoneList()
	if err != nil {
		return "", fmt.Errorf("call zone list api failed, %w", err)
	}
	if zones == nil || zones.Data == nil {
		return "", fmt.Errorf("invalid list zone response")
	}

	for z, zopt := range zones.Data.Zones {
		if z == zone {
			if zopt != nil && zopt.StorageEndpoint != "" {
				return zopt.StorageEndpoint, nil
			}
		}
	}

	return "", fmt.Errorf("zone [%s] not found by list zone api", zone)
}

func (c *Client) CreateShareDirectory(userId snowflake.ID, shareDirectory string) (*schema.SharedDirectory, error) {
	resp, err := c.base.StorageSharedDirectory.Create(
		c.base.StorageSharedDirectory.Create.IgnoreExisting(true),
		c.base.StorageSharedDirectory.Create.Paths([]string{fmt.Sprintf("/%s/%s", userId, shareDirectory)}),
	)
	if err != nil {
		return nil, fmt.Errorf("call create storage shared directory failed, %w", err)
	}
	if resp == nil || resp.Data == nil || len(*resp.Data) == 0 {
		return nil, fmt.Errorf("invalid create shared directory response, %w", err)
	}

	return (*resp.Data)[0], nil
}
