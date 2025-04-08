package rpc

import (
	"context"
	"sync"

	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/proto/hydra_lcp"
	"github.com/yuansuan/ticp/common/project-root-api/proto/idgen"
)

// Client ...
type Client struct {
	Base struct {
		IDGen idgen.IdGenClient `grpc_client_inject:"idgen"`
	}
	HydraLcp struct {
		SsoClient hydra_lcp.HydraLcpServiceClient `grpc_client_inject:"hydra_lcp"`
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
		grpc_boot.RegisterClient("hydra_lcp", hydra_lcp.NewHydraLcpServiceClient)
		grpc_boot.RegisterClient("idgen", idgen.NewIdGenClient)
		//grpc_boot.RegisterClient("account_bill", account_bill.NewAccountServiceClient)
		Instance = &Client{}
		grpc_boot.InjectAllClient(Instance)
	}
	return Instance
}

// GenID generate snowflakes id
func (c *Client) GenID(ctx context.Context) (snowflake.ID, error) {
	reply, err := c.Base.IDGen.GenerateID(ctx, &idgen.GenRequest{})
	if err != nil {
		return snowflake.ID(0), err
	}

	return snowflake.ID(reply.Id), nil
}

func (c *Client) GenIDs(ctx context.Context, count int64) ([]snowflake.ID, error) {
	reply, err := c.Base.IDGen.GenerateIDs(ctx, &idgen.GenerateIDsRequest{Count: count})
	if err != nil {
		return nil, err
	}
	ret := []snowflake.ID{}
	for _, v := range reply.Ids {
		ret = append(ret, snowflake.ID(v))
	}

	return ret, nil
}

// GetSSOUserByID 通过UserID获取SSO用户信息
func (c *Client) GetSSOUserByID(ctx context.Context, ID string) (*hydra_lcp.UserInfo, error) {
	reply, err := c.HydraLcp.SsoClient.GetUserInfo(ctx,
		&hydra_lcp.GetUserInfoReq{
			Ysid: ID,
		})
	if err != nil {
		logging.GetLogger(ctx).Error("err_hydra_lcp_service_on_GetUserInfo", "error", err)
		return nil, err
	}
	return reply, nil
}
