package impl

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/cmd/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/service/client"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

type SysConfigServiceImpl struct {
	sid                   *snowflake.Node
	sysConfigDao          dao.SysConfigDao
	alertManagerConfigDao dao.AlertManagerConfigDao
}

func NewSysConfigService() (service.SysConfigService, error) {
	node, err := snowflake.GetInstance()
	if err != nil {
		logging.Default().Errorf("new snowflake node err: %v", err)
		return nil, err
	}
	notificationDao := dao.NewAlertNotificationDao()
	if err != nil {
		logging.Default().Errorf("new alert notification dao err: %v", err)
		return nil, err
	}
	return &SysConfigServiceImpl{
		sid:                   node,
		sysConfigDao:          dao.NewAppDao(),
		alertManagerConfigDao: notificationDao,
	}, nil
}

func (s *SysConfigServiceImpl) GetGlobalSysConfig(ctx context.Context) (*dto.GetGlobalSysConfigResponse, error) {
	mainConfig := config.Custom.Main

	competeTypeNameMap := mainConfig.ComputeTypeNames
	computeTypeNameList := make([]*dto.ComputeTypeName, 0)
	for i, v := range competeTypeNameMap {
		computeTypeNameList = append(computeTypeNameList, &dto.ComputeTypeName{ComputeType: i, ShowName: v})
	}

	return &dto.GetGlobalSysConfigResponse{
		EnableVisual: mainConfig.EnableVisual,
		ComputeTypes: computeTypeNameList,
	}, nil
}

func (s *SysConfigServiceImpl) GetJobConfig(ctx context.Context) (*dto.GetJobConfigResponse, error) {
	jobConfig := &dto.GetJobConfigResponse{}

	sysConfig, exist, err := s.sysConfigDao.Get(ctx, consts.JobConfig)
	if err != nil {
		return nil, err
	}
	if !exist {
		return jobConfig, nil
	}

	err = json.Unmarshal([]byte(sysConfig.Value), jobConfig)
	if err != nil {
		return nil, err
	}

	return jobConfig, nil
}

func (s *SysConfigServiceImpl) SetJobConfig(ctx context.Context, req *dto.SetJobConfigRequest) error {
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	id := s.sid.Generate()
	err = s.sysConfigDao.Set(ctx, id, consts.JobConfig, string(bytes))
	if err != nil {
		return err
	}

	return nil
}

func (s *SysConfigServiceImpl) GetJobBurstConfig(ctx context.Context) (*dto.GetJobBurstConfigResponse, error) {
	jobBurstConfig := &dto.GetJobBurstConfigResponse{}

	sysConfig, exist, err := s.sysConfigDao.Get(ctx, consts.JobBurstConfig)
	if err != nil {
		return nil, err
	}
	if !exist {
		return jobBurstConfig, nil
	}

	err = json.Unmarshal([]byte(sysConfig.Value), jobBurstConfig)
	if err != nil {
		return nil, err
	}

	return jobBurstConfig, nil
}

func (s *SysConfigServiceImpl) SetJobBurstConfig(ctx context.Context, req *dto.SetJobBurstConfigRequest) error {
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	id := s.sid.Generate()
	err = s.sysConfigDao.Set(ctx, id, consts.JobBurstConfig, string(bytes))
	if err != nil {
		return err
	}

	return nil
}

func (s *SysConfigServiceImpl) GetRBACDefaultRoleId(ctx context.Context) (int64, error) {
	defaultRoleId := int64(0)

	sysConfig, exist, err := s.sysConfigDao.Get(ctx, consts.RBACDefaultRoleId)
	if err != nil {
		return defaultRoleId, err
	}
	if !exist {
		return defaultRoleId, nil
	}

	err = json.Unmarshal([]byte(sysConfig.Value), &defaultRoleId)
	if err != nil {
		return defaultRoleId, err
	}

	return defaultRoleId, nil
}

func (s *SysConfigServiceImpl) SetRBACDefaultRoleId(ctx context.Context, defaultRoleId int64) error {
	defaultRoleIdStr := strconv.FormatInt(defaultRoleId, 10)

	id := s.sid.Generate()
	err := s.sysConfigDao.Set(ctx, id, consts.RBACDefaultRoleId, defaultRoleIdStr)
	if err != nil {
		return err
	}

	return nil
}

func (s *SysConfigServiceImpl) GetThreePersonManagementConfig(ctx context.Context) (*dto.GetThreePersonConfigResponse, error) {
	defaultSafeUserId := &dto.GetThreePersonConfigResponse{}

	sysConfig, exist, err := s.sysConfigDao.Get(ctx, consts.RBACDefaultSafeUser)
	if err != nil {
		return defaultSafeUserId, err
	}
	if !exist {
		return defaultSafeUserId, nil
	}

	err = json.Unmarshal([]byte(sysConfig.Value), &defaultSafeUserId)
	if err != nil {
		return defaultSafeUserId, err
	}
	if strutil.IsNotEmpty(defaultSafeUserId.DefSafeUserID) {
		rsp, err := client.GetInstance().Perm.CheckResourcesPerm(ctx, &rbac.CheckResourcesPermRequest{
			Id: &rbac.ObjectID{
				Id:   defaultSafeUserId.DefSafeUserID,
				Type: rbac.ObjectType_USER,
			},
			Resources: []*rbac.ResourceIdentity{&rbac.ResourceIdentity{
				Identity: &rbac.ResourceIdentity_Name{
					Name: &rbac.ResourceName{
						Type: common.PermissionResourceTypeSystem,
						Name: common.ResourceSecurityManagerName,
					},
				},
			}},
		})
		if err != nil {
			return nil, err
		}

		if !rsp.Pass {
			if err = s.SetThreePersonManagementConfig(ctx, &dto.SetThreePersonConfigRequest{
				DefSafeUserID:   "",
				DefSafeUserName: "",
			}); err != nil {
				return nil, err
			}
			return &dto.GetThreePersonConfigResponse{}, nil
		}
	}

	return defaultSafeUserId, nil
}

func (s *SysConfigServiceImpl) SetThreePersonManagementConfig(ctx context.Context, req *dto.SetThreePersonConfigRequest) error {
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	id := s.sid.Generate()
	err = s.sysConfigDao.Set(ctx, id, consts.RBACDefaultSafeUser, string(bytes))
	if err != nil {
		return err
	}

	return nil
}
