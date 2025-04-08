package rpc

import (
	"context"
	"sync"

	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
	"github.com/yuansuan/ticp/common/project-root-api/proto/idgen"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// Client ...
type Client struct {
	Base struct {
		Gen      idgen.IdGenClient               `grpc_client_inject:"idgen"`
		HydraLcp hydra_lcp.HydraLcpServiceClient `grpc_client_inject:"hydra_lcp"`
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

// GenID generate a snowflake id
func GenID(ctx context.Context) (snowflake.ID, error) {
	reply, err := GetInstance().Base.Gen.GenerateID(ctx, &idgen.GenRequest{})
	if err != nil {
		return snowflake.ID(0), err
	}
	return snowflake.ID(reply.Id), nil
}

func GenIDs(ctx context.Context, count int64) ([]snowflake.ID, error) {
	reply, err := GetInstance().Base.Gen.GenerateIDs(ctx, &idgen.GenerateIDsRequest{Count: count})
	if err != nil {
		return nil, err
	}
	ret := []snowflake.ID{}
	for _, v := range reply.Ids {
		ret = append(ret, snowflake.ID(v))
	}

	return ret, nil
}

func BatchCheckUserExist(ctx context.Context, users []string) ([]snowflake.ID, error) {
	resp, err := GetInstance().Base.HydraLcp.GetUserInfoBatch(ctx, &hydra_lcp.GetUserInfoBatchReq{
		Ysid: users,
	})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			if s.Code() == consts.ErrHydraLcpDBUserNotExist {
				return []snowflake.ID{}, nil
			}
		}
		return nil, err
	}

	res := make([]snowflake.ID, 0)
	for _, userInfo := range resp.GetUserInfo() {
		res = append(res, snowflake.MustParseString(userInfo.Ysid))
	}

	return res, nil
}
