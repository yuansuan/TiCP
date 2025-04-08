package rpc

import (
	"context"
	grpcboot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/proto/idgen"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/errors"
	"sync"
)

// Client ...
type Client struct {
	IDGen idgen.IdGenClient `grpc_client_inject:"idgen"`
}

var instance *Client
var mtx sync.Mutex

// GetInstance get instance
func GetInstance() *Client {
	mtx.Lock()
	defer mtx.Unlock()
	if instance == nil {
		grpcboot.RegisterClient("idgen", idgen.NewIdGenClient)
		instance = &Client{}
		grpcboot.InjectAllClient(instance)
	}
	return instance
}

func GenID(ctx context.Context) (snowflake.ID, error) {
	resp, err := GetInstance().IDGen.GenerateID(ctx, &idgen.GenRequest{})
	if err != nil {
		logging.Default().Errorf("generate snowflake id err, err: %v", err)
		return snowflake.ID(0), errors.ErrSnowflakeGeneration
	}

	return snowflake.ID(resp.Id), nil
}
