package rpc

import (
	"context"
	"sync"

	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/company"
	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/idgen"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"

	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
)

// Client ...
type Client struct {
	Base struct {
		Gen idgen.IdGenClient `grpc_client_inject:"idgen"`
	}
	Company struct {
		Company company.CompanyServiceClient `grpc_client_inject:"company"`
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
