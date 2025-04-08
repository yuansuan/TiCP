package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"go.uber.org/multierr"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/domain"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/state"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
)

const (
	collectInstanceInterval  = 10 * time.Second
	waitCloseSessionInterval = 10 * time.Second
	leakInstanceInterval     = 3 * time.Minute
	recentInstanceDuration   = 5 * time.Minute
	defaultChannelSize       = 1024
)

// Scheduler schedule instance
type Scheduler struct {
	ctx context.Context
	wg  *sync.WaitGroup

	cfg config.CustomT
	log *logging.Logger

	postPaidBillMsgChan chan PostPaidBillSliceMessage

	state *state.State
}

func (s *Scheduler) Start() {
	s.wg.Add(1)
	// 从所有可视化集群中获取所有实例（部分状态的实例获取不到）并更新数据库中的状态：10sec
	go s.collectAndUpdateInstance(collectInstanceInterval)

	// 按分区获取所有用户请求关闭的会话并调用对应可视化集群接口关闭实例：10sec
	for zone := range s.cfg.CloudApp.Zones {
		s.wg.Add(1)
		go s.collectAndCloseWaitingCloseSessionsByZone(waitCloseSessionInterval, zone.String())
	}

	s.wg.Add(1)
	// 获取所有长时间未更新且未结束的会话，并尝试恢复或关闭它：5min
	go s.collectAndUpdateUnterminatedInstances(leakInstanceInterval)

	s.wg.Wait()
}

// collectAndUpdateInstance 从所有可视化集群中获取所有运行中的实例，并将实例状态更新到数据库中
func (s *Scheduler) collectAndUpdateInstance(interval time.Duration) {
	defer s.wg.Done()

	for {
		s.log.Debug("collect instances from cloud ...")
		if instances, err := s.state.Cloud.DescribeInstances(); err != nil {
			s.log.Warn("query instances: %s", err)
		} else {
			for _, instance := range instances {
				err = s.updateInstance(instance)
				if err != nil {
					s.log.Warnf("update instance [%s] error: %s", instance.ID(), err)
				}
			}
		}

		time.Sleep(interval)
	}
}

// updateInstance 将可视化集群中的实例状态同步到数据库中
func (s *Scheduler) updateInstance(instance domain.Instance) error {
	s.log.Infow("update instance from cloud", "instance", instance)
	return with.DefaultTransaction(s.ctx, func(ctx context.Context) error {
		sd, err := dao.GetSessionDetailByRawInstanceId(ctx, instance.ID())
		if err != nil {
			return errors.Wrap(err, "session not found in database")
		}

		bootVolumeId, err := s.state.Cloud.ParseBootVolumeId(sd.Session.Zone, instance.Raw())
		if err != nil {
			// log warn but not interrupt update below
			s.log.Warnf("parse boot volume id failed, %v", err)
		}

		if err = dao.UpdateInstance(ctx, &models.Instance{
			Id:           sd.Instance.Id,
			InstanceData: instance.Raw(),
			BootVolumeId: bootVolumeId,
		}, []string{"instance_data", "boot_volume_id"}); err != nil {
			return fmt.Errorf("update instance failed, %w", err)
		}

		switch instance.Status() {
		case models.InstanceRunning:
			err = multierr.Append(
				dao.UpdateInstanceStatus(ctx, sd.Instance.Id, models.InstanceRunning),
				dao.UpdateSessionStatus(ctx, sd.Session.Id,
					[]schema.SessionStatus{schema.SessionStarting, schema.SessionRebooting, schema.SessionPoweringOn, schema.SessionPowerOff},
					schema.SessionStarted),
			)
			if err != nil {
				break
			}
		case models.InstanceRebooting:
			err = dao.UpdateInstanceStatus(ctx, sd.Instance.Id, models.InstanceRebooting)
		case models.InstanceStopping:
			err = multierr.Append(
				dao.UpdateInstanceStatus(ctx, sd.Instance.Id, models.InstanceStopping),
				dao.UpdateSessionStatus(ctx, sd.Session.Id, []schema.SessionStatus{schema.SessionStarted}, schema.SessionPoweringOff),
			)
		case models.InstanceStopped:
			err = multierr.Append(
				dao.UpdateInstanceStatus(ctx, sd.Instance.Id, models.InstanceStopped),
				dao.UpdateSessionStatus(ctx, sd.Session.Id, []schema.SessionStatus{schema.SessionPoweringOff, schema.SessionStarted}, schema.SessionPowerOff),
			)
		case models.InstanceStarting:
			err = multierr.Append(
				dao.UpdateInstanceStatus(ctx, sd.Instance.Id, models.InstanceStarting),
				dao.UpdateSessionStatus(ctx, sd.Session.Id, []schema.SessionStatus{schema.SessionPowerOff}, schema.SessionPoweringOn),
			)
		case models.InstanceLaunchFailed: // 创建失败
			err = s.shallowCloseSessionWithReason(ctx, sd.Session.Id, sd.Instance.Id, "launch failure")
		case models.InstanceShutdown, models.InstanceTerminating: // 停止待销毁, 销毁中
			err = s.shallowCloseSessionWithReason(ctx, sd.Session.Id, sd.Instance.Id, "instance shutdown")
		}
		return err
	})
}

// collectAndCloseWaitingCloseSessionsByZone 从数据库中获取指定分区待关闭的会话(用户请求)
func (s *Scheduler) collectAndCloseWaitingCloseSessionsByZone(interval time.Duration, zoneName string) {
	defer s.wg.Done()

	for {
		time.Sleep(interval)

		s.log.Debugf("collect zone %s waiting close session from database ...", zoneName)
		if sessions, err := dao.ListWaitingCloseSessionByZone(s.ctx, zoneName); err != nil {
			s.log.Warnf("collect zone %s waitingCloseSession error: %s", zoneName, err)
		} else {
			for _, session := range sessions {
				err = s.closeWaitingCloseSession(session)
				if err != nil {
					s.log.Warnf("close zone %s waitingCloseSession error: %s", zoneName, err)
				}
			}
		}
	}
}

// closeWaitingCloseSession 调用三方云接口关闭用户指定的实例
func (s *Scheduler) closeWaitingCloseSession(session *models.SessionWithDetail) error {
	var (
		instance *models.Instance
		err      error
	)

	err = with.DefaultTransaction(s.ctx, func(ctx context.Context) error {
		instance, err = dao.GetInstance(ctx, session.Session.InstanceId)
		if err != nil {
			return errors.Wrap(err, "instance not found in database")
		}

		// already closed
		if instance.InstanceStatus == models.InstanceTerminated {
			return dao.SessionClosed(ctx, session.Session.Id, "user closed[already closed]")
		}

		if err = s.state.Cloud.TerminateInstance(instance.Zone, instance.InstanceId); err != nil {
			// 无效的实例ID可能是实例启动过程中断导致
			if err == consts.ErrMalformedInstanceID {
				s.log.Errorw("malformed instance id", "instance", instance)
				return s.shallowCloseSessionWithReason(ctx, session.Session.Id, session.Session.InstanceId, "user closed[malformed instance id]")
			}

			// 找不到实例可能是各种时间点凑到边界上导致
			if err == consts.ErrInstanceNotFound {
				return s.shallowCloseSessionWithReason(ctx, session.Session.Id, session.Session.Id, "user closed[instance not found]")
			}

			return err
		}

		return s.shallowCloseSessionWithReason(ctx, session.Session.Id, session.Session.InstanceId, "user closed")
	})

	return err
}

// collectAndUpdateUnterminatedInstances 从数据库中获取所有最近一段时间未关闭的实例
func (s *Scheduler) collectAndUpdateUnterminatedInstances(interval time.Duration) {
	defer s.wg.Done()

	for {
		s.log.Debug("collect unterminated instance from database ...")
		// 查询更新时间在 5min 之前的所有未关闭的实例，并检查这些实例的状态
		if instances, err := dao.ListRecentUnterminatedInstance(s.ctx, recentInstanceDuration); err != nil {
			s.log.Warnf("collect Unterminated instance error: %s", err)
		} else {
			for _, instance := range instances {
				err = s.updateUnterminatedInstance(instance)
				if err != nil {
					s.log.Warnf("update Unterminated instance [%s] error: %s", instance.Id, err)
				}
			}
		}

		time.Sleep(interval)
	}
}

// updateUnterminatedInstance 处理所有长时间未更新的会话信息，尝试获取实例并检查其状态
func (s *Scheduler) updateUnterminatedInstance(instance *models.Instance) error {
	s.log.Debugw("consume leakInstance", "instance", instance)

	session, err := dao.GetSessionByInstancePk(s.ctx, instance.Id)
	if err != nil {
		s.log.Errorw("session not found", "error", err)
		return err
	}

	vm, err := s.state.Cloud.DescribeInstance(session.Zone, instance.InstanceId)
	if err != nil {
		if err == consts.ErrMalformedInstanceID {
			err = s.shallowCloseSessionWithReason(s.ctx, session.Id, instance.Id, "leak[malformed instance id]")
		} else if errors.Is(err, consts.ErrInstanceNotFound) {
			err = s.shallowCloseSessionWithReason(s.ctx, session.Id, instance.Id, "leak[instance not found]")
		}

		s.log.Warnf("update Unterminated instance error: %s", err)
		return err
	}

	if vm != nil {
		s.log.Debugf("instance(%s) alive: session => %s", instance.InstanceId, session.Id)
		// 强制更新实例状态
		return s.updateInstance(vm)
	}

	// unreached
	if err = s.shallowCloseSessionWithReason(s.ctx, session.Id, instance.Id, "leak session"); err != nil {
		s.log.Warnf("consume instance error: %s", err)
		return err
	}

	return nil
}

// shallowCloseSessionWithReason 关闭一个会话并更新关闭原因，注意这个方式仅修改数据库！
func (s *Scheduler) shallowCloseSessionWithReason(ctx context.Context, sessionID, instanceID snowflake.ID, reason string) error {
	return multierr.Append(
		dao.InstanceTerminated(ctx, instanceID),
		dao.SessionClosed(ctx, sessionID, reason),
	)
}

// NewScheduler 创建调度器用于关闭会话、同步状态以及计费
func NewScheduler(st *state.State) (*Scheduler, error) {
	if st == nil {
		return nil, fmt.Errorf("state is nil")
	}

	s := &Scheduler{
		ctx:                 context.Background(),
		wg:                  &sync.WaitGroup{},
		cfg:                 config.GetConfig(),
		log:                 logging.Default(),
		postPaidBillMsgChan: make(chan PostPaidBillSliceMessage, defaultChannelSize),
		state:               st,
	}

	return s, nil
}
