package handler_rpc

import (
	"fmt"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/proto/idgen"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/idgen/config"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"golang.org/x/net/context"
)

type server struct {
	node *snowflake.Node

	idgen.UnimplementedIdGenServer
}

const prefix = "/idgen/"

// GenerateID GenerateID
func (s *server) GenerateID(ctx context.Context, in *idgen.GenRequest) (*idgen.GenReply, error) {
	// Generate a snowflake ID.
	id := s.node.Generate()
	return &idgen.GenReply{Id: int64(id)}, nil
}

func (s *server) GenerateIDs(ctx context.Context, in *idgen.GenerateIDsRequest) (*idgen.GenerateIDsReply, error) {
	// Generate a snowflake ID.
	ids := []int64{}
	for i := int64(0); i < in.Count; i++ {
		ids = append(ids, s.node.Generate().Int64())
	}
	return &idgen.GenerateIDsReply{Ids: ids}, nil
}

var onexit func()

// OnShutdown OnShutdown
func OnShutdown(drv *http.Driver) {
	onexit()
}

// InitGRPCServer InitGRPCServer
func InitGRPCServer(drv *http.Driver) {
	svr := &server{}
	ctx, cancel := context.WithCancel(context.Background())
	logger := logging.Default()

	id := -1
	if config.GetConfig().NodeID <= 0 {
		client := boot.MW.DefaultEtcd()
		s, err := concurrency.NewSession(client, concurrency.WithTTL(10))
		util.ChkErr(err)

		onexit = func() {
			logger.Infof("close here")
			cancel()
			s.Close()
		}

		for i := 0; i < (1 << snowflake.NodeBits); i++ {
			key := prefix + fmt.Sprintf("%v", i)
			client := s.Client()

			cmp := v3.Compare(v3.CreateRevision(key), "=", 0)

			// put self in lock waiters via myKey; oldest waiter holds lock
			put := v3.OpPut(key, "", v3.WithLease(s.Lease()))

			resp, err := client.Txn(ctx).If(cmp).Then(put).Commit()
			if err != nil {
				continue
			}
			if err == nil && resp.Succeeded {
				id = i
				break
			}
		}
	} else {
		id = config.GetConfig().NodeID
	}
	if id == -1 {
		util.ChkErr(fmt.Errorf("error in set node id"))
	}
	logger.Infof("start new node with id %v", id)
	node, err := snowflake.NewNode(int64(id))
	util.ChkErr(err)
	svr.node = node

	ss, err := boot.GRPC.DefaultServer()
	util.ChkErr(err)

	idgen.RegisterIdGenServer(ss.Driver(), svr)
}
