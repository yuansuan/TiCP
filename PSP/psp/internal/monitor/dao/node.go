package dao

import (
	"context"
	"fmt"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"google.golang.org/grpc/status"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/monitor"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type nodeDaoImpl struct{}

func NewNodeDao() NodeDao {
	return &nodeDaoImpl{}
}

func (d *nodeDaoImpl) AddNodes(ctx context.Context, session *xorm.Session, nodes []*model.NodeInfo) error {
	if session == nil {
		session = boot.MW.DefaultSession(ctx)
		defer session.Close()
	}

	_, err := session.Insert(nodes)
	if err != nil {
		return err
	}
	return nil
}
func (d *nodeDaoImpl) UpdateNode(ctx context.Context, session *xorm.Session, node *model.NodeInfo) error {
	if session == nil {
		session = boot.MW.DefaultSession(ctx)
		defer session.Close()
	}

	_, err := session.ID(node.Id).
		Cols("node_name", "node_type", "queue_name", "platform_name", "scheduler_status", "status", "total_core_num", "used_core_num", "free_core_num",
			"total_mem", "used_mem", "free_mem", "available_mem", "update_time").Update(node)
	if err != nil {
		return err
	}
	return nil
}

func (d *nodeDaoImpl) DeleteNotIds(ctx context.Context, session *xorm.Session, ids []snowflake.ID) error {
	if session == nil {
		session = boot.MW.DefaultSession(ctx)
		defer session.Close()
	}

	_, err := session.NotIn("id", ids).Delete(&model.NodeInfo{})
	if err != nil {
		return err
	}
	return nil
}

func (d *nodeDaoImpl) GetNodes(ctx context.Context, nodeName string, index, size int64) ([]*model.NodeInfo, int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var nodes []*model.NodeInfo

	total, err := session.Where("node_name Like ?", "%"+nodeName+"%").Count(&model.NodeInfo{})
	if err != nil {
		msg := fmt.Sprintf("query node failed %v", err)
		return nil, 0, status.Error(errcode.ErrNodeQueryFailed, msg)
	}

	session.Where("node_name Like ?", "%"+nodeName+"%").Asc("node_name")

	err = session.Limit(int(size), int((index-1)*size)).Find(&nodes)
	if err != nil {
		msg := fmt.Sprintf("query node failed %v", err)
		return nil, 0, status.Error(errcode.ErrNodeQueryFailed, msg)
	}
	return nodes, total, nil
}

func (d *nodeDaoImpl) NodeList(ctx context.Context) ([]*model.NodeInfo, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var nodes []*model.NodeInfo

	err := session.Asc("node_name").Find(&nodes)
	if err != nil {
		msg := fmt.Sprintf("query node failed %v", err)
		return nil, status.Error(errcode.ErrNodeQueryFailed, msg)
	}
	return nodes, nil
}

func (d *nodeDaoImpl) GetNodeByNames(ctx context.Context, nodeNames []string) ([]*model.NodeInfo, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var nodes []*model.NodeInfo
	if err := session.In("node_name", nodeNames).Find(&nodes); err != nil {
		msg := fmt.Sprintf("query nodes failed %v", err)
		return nil, status.Error(errcode.ErrNodeQueryFailed, msg)
	}
	return nodes, nil
}

type Statistics struct {
	TotalNum int64 `xorm:"total_num"`
	FreeNum  int64 `xorm:"free_num"`
}

func (d *nodeDaoImpl) StatisticCoreNum(ctx context.Context, nodeNames []string) (*Statistics, error) {
	session := boot.MW.DefaultSession(ctx)

	var statistics []*Statistics
	session = session.Table(model.TableName)
	if nodeNames != nil && len(nodeNames) > 0 {
		session.In("node_name", nodeNames)
	}
	session.Select("sum(total_core_num) as total_num, sum(free_core_num) as free_num")
	if err := session.Find(&statistics); err != nil {
		return nil, err
	}
	if len(statistics) == 0 {
		return nil, status.Error(errcode.ErrNodeNotExist, errcode.MsgNodeNotExist)
	}

	return statistics[0], nil
}

func (d *nodeDaoImpl) QueueList(ctx context.Context) ([]string, error) {
	session := boot.MW.DefaultSession(ctx)

	var queues []string
	session = session.Table(model.TableName)
	err := session.Distinct("queue_name").Asc("queue_name").Find(&queues)
	if err != nil {
		return nil, err
	}
	return queues, nil
}

func (d *nodeDaoImpl) GetQueueAvailableCores(ctx context.Context, queueNames []string) ([]*pb.QueueCore, error) {
	session := boot.MW.DefaultSession(ctx)

	var queueCores []*pb.QueueCore
	session = session.Table(model.TableName)
	if queueNames != nil && len(queueNames) > 0 {
		session.In("queue_name", queueNames)
	}
	session.Where("status = ?", consts.Idle)
	session.GroupBy("queue_name")
	session.Select("sum(free_core_num) as core_num, queue_name")
	if err := session.Find(&queueCores); err != nil {
		return nil, err
	}

	return queueCores, nil
}

func (d *nodeDaoImpl) GetPlatformCores(ctx context.Context) ([]*pb.PlatformCore, error) {
	session := boot.MW.DefaultSession(ctx)

	var platformCore []*pb.PlatformCore
	session = session.Table(model.TableName)
	session.Where("status = ?", consts.Idle)
	session.GroupBy("platform_name")
	session.Select("sum(free_core_num) as available_cores,sum(total_core_num) as total_cores, platform_name")
	if err := session.Find(&platformCore); err != nil {
		return nil, err
	}

	return platformCore, nil
}

func (d *nodeDaoImpl) GetQueueCoreInfos(ctx context.Context) ([]*pb.QueueCoreInfo, error) {
	session := boot.MW.DefaultSession(ctx)

	var queueCoreInfos []*pb.QueueCoreInfo
	session = session.Table(model.TableName)
	session.Where("status = ?", consts.Idle)
	session.GroupBy("queue_name")
	session.Select("sum(free_core_num) as available_cores,sum(total_core_num) as total_cores, queue_name")
	if err := session.Find(&queueCoreInfos); err != nil {
		return nil, err
	}

	return queueCoreInfos, nil
}
