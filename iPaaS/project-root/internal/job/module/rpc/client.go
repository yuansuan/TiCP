package rpc

import (
	"context"
	"sync"

	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"

	"github.com/yuansuan/ticp/common/project-root-api/proto/license"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"

	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	idgen "github.com/yuansuan/ticp/common/project-root-api/proto/idgen"
)

// Client ...
type Client struct {
	Idgen struct {
		Gen idgen.IdGenClient `grpc_client_inject:"idgen"`
	}
	HydraLcp struct {
		HydraLcp hydra_lcp.HydraLcpServiceClient `grpc_client_inject:"hydra_lcp"`
	}
	License struct {
		LicenseServer license.LicenseManagerServiceClient `grpc_client_inject:"license_server"`
	}
}

// Instance ...
var Instance *Client

var mtx sync.Mutex

// GetInstance get instance
func GetInstance() *Client {
	mtx.Lock()
	defer mtx.Unlock()

	if Instance == nil {
		Instance = &Client{}
		grpc_boot.InjectAllClient(Instance)
	}
	return Instance
}

// GenID generate a snowslake id
func (c *Client) GenID(ctx context.Context) (snowflake.ID, error) {
	reply, err := GetInstance().Idgen.Gen.GenerateID(ctx, &idgen.GenRequest{})
	if err != nil {
		return snowflake.ID(0), err
	}
	return snowflake.ID(reply.Id), nil
}

// GetUser get user info
func (c *Client) GetUser(ctx context.Context, userID string) (*hydra_lcp.UserInfo, error) {
	reply, err := GetInstance().HydraLcp.HydraLcp.GetUserInfo(ctx, &hydra_lcp.GetUserInfoReq{
		Ysid: userID,
	})

	if err != nil {
		trace.GetLogger(ctx).Warnf("get user info failed, err: %v", err)
		if s, ok := status.FromError(err); ok && s.Code() == consts.ErrHydraLcpDBUserNotExist {
			return nil, nil
		}
		return nil, err
	}

	return reply, nil
}
