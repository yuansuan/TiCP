package impl

import (
	"context"
	"fmt"
	openapiapp "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/app"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"

	"github.com/yuansuan/ticp/PSP/psp/internal/app/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/util"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/license"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/monitor"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/maputil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/serializeutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"google.golang.org/grpc/status"
)

type AppService struct {
	sid    *snowflake.Node
	api    *openapi.OpenAPI
	appDao dao.AppDao
}

func NewAppService() (*AppService, error) {
	node, err := snowflake.GetInstance()
	if err != nil {
		logging.Default().Errorf("new snowflake node err: %v", err)
		return nil, err
	}
	api, err := openapi.NewLocalAPI()
	if err != nil {
		logging.Default().Errorf("new local openapi err: %v", err)
		return nil, err
	}
	return &AppService{
		sid:    node,
		api:    api,
		appDao: dao.NewAppDao(),
	}, nil
}

func (s *AppService) ListApp(ctx context.Context, userId int64, computeType, state string, hasPermission, desktop bool) ([]*dto.App, error) {
	resourceIds := make([]snowflake.ID, 0)
	if hasPermission {
		var resourceTypeList []string
		switch computeType {
		case common.Local:
			resourceTypeList = append(resourceTypeList, common.PermissionResourceTypeLocalApp)
		case common.Cloud:
			resourceTypeList = append(resourceTypeList, common.PermissionResourceTypeAppCloudApp)
		case common.Empty:
			resourceTypeList = append(resourceTypeList, common.PermissionResourceTypeLocalApp, common.PermissionResourceTypeAppCloudApp)
		}

		permissions, err := client.GetInstance().RBAC.Permission.ListObjectResources(ctx, &rbac.ListObjectResourcesRequest{
			Id: &rbac.ObjectID{
				Id:   snowflake.ID(userId).String(),
				Type: rbac.ObjectType_USER,
			},
			ResourceType: resourceTypeList,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range permissions.Perms {
			resourceIds = append(resourceIds, snowflake.ID(v.ResourceId))
		}

		tracelog.Info(ctx, fmt.Sprintf("user: [%d] has permission app, resourceIds: %+v", userId, resourceIds))

		if len(resourceIds) == 0 {
			return []*dto.App{}, nil
		}
	}

	app := &model.App{ComputeType: computeType, State: state}
	apps, err := s.appDao.ListApp(ctx, resourceIds, app)
	if err != nil {
		return nil, err
	}

	var queueNames []string
	var queueInfos []*dto.QueueInfo
	var licenses []*dto.LicenseInfo
	// 桌面请求时不请求 Paas 数据
	if !desktop {
		queueInfos, queueNames, err = getQueueNames(ctx)
		if err != nil {
			return nil, err
		}

		tracelog.Info(ctx, fmt.Sprintf("list queueInfo: [%v], queueNames: [%+v]", queueInfos, queueNames))

		licenses, _, err = getLicenseInfo(ctx)
		if err != nil {
			return nil, err
		}

		tracelog.Info(ctx, fmt.Sprintf("list license: [%v]", serializeutil.GetStringForTraceLog(licenses)))
	}

	appInfoList := make([]*dto.App, 0, len(apps))
	for _, v := range apps {
		appInfo, err := util.ConvertAppDto(v, queueInfos, licenses)
		if err != nil {
			return nil, err
		}

		// 桌面请求时不请求 Paas 数据
		if !desktop {
			// 队列名称不为空时, 获取真实的队列核数信息
			if hasPermission && v.ComputeType == common.Local && len(v.QueueNames) != 0 {
				queueInfosSelected, err := getQueueNamesWithCPUNumber(ctx, v.QueueNames)
				if err != nil {
					return nil, err
				}

				if len(queueInfos) != 0 {
					appInfo.Queues = queueInfosSelected
				}
			}

			// 队列名称为空时, 需要在提交作业时显示所有的队列核数信息
			if hasPermission && v.ComputeType == common.Local && len(v.QueueNames) == 0 {
				queueInfosAll, err := getQueueNamesWithCPUNumber(ctx, queueNames)
				if err != nil {
					return nil, err
				}

				if len(queueInfosAll) != 0 {
					appInfo.Queues = queueInfosAll
				}
			}
		}

		appInfoList = append(appInfoList, appInfo)
	}

	return appInfoList, nil
}

func (s *AppService) GetAppInfo(ctx context.Context, req *dto.GetAppInfoServiceRequest) (*dto.App, error) {
	app := &model.App{ID: req.ID, OutAppID: req.OutAppID, Name: req.Name, Type: req.AppType, VersionNum: req.Version, ComputeType: req.ComputeType}
	exist, err := s.appDao.GetApp(ctx, app)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, status.Errorf(errcode.ErrAppTemplateNotExist, "app template not exist")
	}

	// 检测本地应用 OutAppID 为空的数据并进行数据补偿
	if app.ComputeType == common.Local && app.OutAppID == "" {
		getApp, exist, err := checkAndGetApp(ctx, s.api, app.Name)
		if err != nil {
			return nil, err
		}
		if !exist {
			response, err := openapiapp.AddApp(s.api, app.Name, app.Type, app.VersionNum, "", "", "", false, false, nil)
			if err != nil {
				return nil, err
			}
			if response == nil || response.ErrorCode != "" {
				return nil, fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
			}
			app.OutAppID = response.Data.AppID
		} else {
			app.OutAppID = getApp.AppID
			app.Image = getApp.Image
			app.LicenseManagerId = getApp.LicManagerId

			binPathMap, err := util.ConvertBinPathMap(getApp.BinPath)
			if err != nil {
				return nil, fmt.Errorf("convert bin path: [%v] to bin path map err: %v", getApp.BinPath, err)
			}
			app.BinPath = binPathMap
		}

		_ = s.appDao.UpdateApp(ctx, app)

		tracelog.Info(ctx, fmt.Sprintf("local app: [%v] compensate out_app_id info: [%v]", app.ID, app.OutAppID))
	}

	queueInfos, _, err := getQueueNames(ctx)
	if err != nil {
		return nil, err
	}
	licenses, _, err := getLicenseInfo(ctx)
	if err != nil {
		return nil, err
	}

	appInfo, err := util.ConvertAppDto(app, queueInfos, licenses)
	if err != nil {
		return nil, err
	}

	return appInfo, nil
}

func (s *AppService) GetAppTotalNum(ctx context.Context) (int64, error) {
	count, err := s.appDao.GetAppTotalNum(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *AppService) AddApp(ctx context.Context, req *dto.AddAppServiceRequest) error {
	newName := fmt.Sprintf("%v%v%v", req.NewType, common.Blank, req.NewVersion)
	app := &model.App{Name: newName, ComputeType: req.ComputeType}
	newAppExist, err := s.appDao.GetApp(ctx, app)
	if err != nil {
		return err
	}
	if newAppExist {
		return status.Errorf(errcode.ErrAppTemplateHasExist, "app template has exist")
	}
	if req.ComputeType == common.Cloud {
		return status.Errorf(errcode.ErrAppUnableCreateNewCloudApp, "unable create new cloud app ")
	}

	baseApp := &model.App{}
	if req.BaseName != "" {
		baseApp = &model.App{Name: req.BaseName, ComputeType: req.ComputeType}
		baseAppExist, err := s.appDao.GetApp(ctx, baseApp)
		if err != nil {
			return err
		}
		if !baseAppExist {
			return status.Errorf(errcode.ErrAppBaseTemplateNotExist, "base template not exist")
		}
	}

	newApp := baseApp
	newApp.ID = s.sid.Generate()
	newApp.Name = newName
	newApp.Type = req.NewType
	newApp.VersionNum = req.NewVersion
	newApp.ComputeType = common.Local
	newApp.QueueNames = getSelectedQueueNames(req.Queues)
	newApp.LicenseManagerId = getSelectedLicenseInfo(req.Licenses)
	newApp.State = common.Unpublished
	newApp.Description = req.Description
	newApp.EnableResidual = req.EnableResidual
	newApp.ResidualLogParser = req.ResidualLogParser
	newApp.EnableSnapshot = req.EnableSnapshot
	newApp.Image = req.Image
	newApp.Icon = req.Icon
	newApp.BinPath = util.GetKeyValueMap(req.BinPath)
	newApp.SchedulerParam = util.GetKeyValueMap(req.SchedulerParam)

	if req.CloudOutAppId != "" {
		tmpApp := &model.App{OutAppID: req.CloudOutAppId, ComputeType: common.Cloud}
		tmpAppExist, err := s.appDao.GetApp(ctx, tmpApp)
		if err != nil {
			return err
		}
		if !tmpAppExist {
			return status.Errorf(errcode.ErrAppBindCloudTemplateNotExist, "need bind app template not exist")
		}
		newApp.OutID = req.CloudOutAppId
		newApp.OutName = tmpApp.Name
	}

	getApp, exist, err := checkAndGetApp(ctx, s.api, newApp.Name)
	if err != nil {
		return err
	}
	if !exist {
		response, err := openapiapp.AddApp(s.api, newApp.Name, newApp.Type, newApp.VersionNum, newApp.Image, newApp.LicenseManagerId, newApp.ResidualLogParser, newApp.EnableResidual, newApp.EnableSnapshot, util.GetKeyValueMap(req.BinPath))
		if err != nil {
			return err
		}
		if response == nil || response.ErrorCode != "" {
			return fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
		}
		newApp.OutAppID = response.Data.AppID
	} else {
		newApp.OutAppID = getApp.AppID

		if checkAppChangedForLocal(getApp, newApp) {
			keyValue, err := util.ConvertBinPathMap(getApp.BinPath)
			if err != nil {
				return err
			}

			response, err := openapiapp.UpdateApp(s.api, newApp.OutAppID, newApp.Name, newApp.Type, newApp.VersionNum, newApp.Image, newApp.LicenseManagerId, newApp.ResidualLogParser, newApp.EnableResidual, newApp.EnableSnapshot, keyValue, update.Published)
			if err != nil {
				return err
			}
			if response == nil || response.ErrorCode != "" {
				return fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
			}

			tracelog.Info(ctx, fmt.Sprintf("in add app operate, update app: [%v] to paas data, related outAppId: [%v]", newApp.ID, newApp.OutAppID))
		}
	}

	err = s.appDao.AddApps(ctx, []*model.App{newApp})
	if err != nil {
		return err
	}

	err = AddAppPermission(ctx, common.PermissionResourceTypeLocalApp, s, newApp)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("add app: [%v]", serializeutil.GetStringForTraceLog(util.GetSimpleAppValue(newApp))))

	return nil
}

func (s *AppService) UpdateApp(ctx context.Context, appDto *dto.App, baseName string) error {
	newName := fmt.Sprintf("%v%v%v", appDto.Type, common.Blank, appDto.Version)
	if appDto.ComputeType == common.Cloud {
		newName = appDto.Name
	}

	app := &model.App{Name: newName, ComputeType: appDto.ComputeType}
	exist, err := s.appDao.GetApp(ctx, app)
	if err != nil {
		return err
	}
	if baseName == "" && !exist {
		return status.Errorf(errcode.ErrAppTemplateNotExist, "app template not exist")
	}
	if baseName != "" && exist {
		return status.Errorf(errcode.ErrAppTemplateHasExist, "app template has exist")
	}
	if app.State == common.Published {
		return status.Errorf(errcode.ErrAppTemplateHasPublished, "app template: [%v] has published", newName)
	}
	if app.ComputeType != appDto.ComputeType {
		return status.Errorf(errcode.ErrAppUnableUpdateAppCompeteType, "unable update app compute type")
	}

	if baseName != "" {
		baseApp := &model.App{Name: baseName, ComputeType: appDto.ComputeType}
		_, err = s.appDao.GetApp(ctx, baseApp)
		if err != nil {
			return err
		}
		if app.ComputeType == common.Cloud {
			return status.Errorf(errcode.ErrAppUnableCreateNewAPPBaseCloudApp, "unable base on cloud app to create new app ")
		}

		app = baseApp
	}

	originApp := &model.App{}
	*originApp = *app

	app.Name = newName
	app.Type = appDto.Type
	app.VersionNum = appDto.Version
	app.ComputeType = appDto.ComputeType
	app.QueueNames = getSelectedQueueNames(appDto.Queues)
	app.LicenseManagerId = getSelectedLicenseInfo(appDto.Licenses)
	app.State = common.Unpublished
	app.Description = appDto.Description
	app.EnableResidual = appDto.EnableResidual
	app.ResidualLogParser = appDto.ResidualLogParser
	app.EnableSnapshot = appDto.EnableSnapshot
	app.Image = appDto.Image
	app.BinPath = util.GetKeyValueMap(appDto.BinPath)
	app.SchedulerParam = util.GetKeyValueMap(appDto.SchedulerParam)

	if appDto.CloudOutAppID != "" {
		tmpApp := &model.App{OutAppID: appDto.CloudOutAppID, ComputeType: common.Cloud}
		tmpAppExist, err := s.appDao.GetApp(ctx, tmpApp)
		if err != nil {
			return err
		}
		if !tmpAppExist {
			return status.Errorf(errcode.ErrAppRelationTemplateNotExist, "relation app template not exist")
		}
		app.OutID = appDto.CloudOutAppID
		app.OutName = tmpApp.Name
	} else {
		app.OutID = ""
		app.OutName = ""
	}

	if appDto.Icon != "" {
		app.Icon = appDto.Icon
	}
	app.Script = fmt.Sprintf("%v%v%v", appDto.Type, common.Dot, consts.Cmd)
	if appDto.HelpDoc != nil {
		app.DocType = appDto.HelpDoc.Type
		app.DocContent = appDto.HelpDoc.Value
	}
	content, err := yaml.Marshal(appDto.SubForm)
	if err != nil {
		return err
	}
	app.Content = string(content)
	app.Script = appDto.Script

	if baseName != "" {
		getApp, exist, err := checkAndGetApp(ctx, s.api, app.Name)
		if err != nil {
			return err
		}
		if !exist {
			response, err := openapiapp.AddApp(s.api, app.Name, app.Type, app.VersionNum, app.Image, app.LicenseManagerId, app.ResidualLogParser, app.EnableResidual, app.EnableSnapshot, util.GetKeyValueMap(appDto.BinPath))
			if err != nil {
				return err
			}
			if response == nil || response.ErrorCode != "" {
				return fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
			}
			app.OutAppID = response.Data.AppID

			tracelog.Info(ctx, fmt.Sprintf("in update app operate, add app: [%v] to paas data, related outAppId: [%v]", app.ID, app.OutAppID))
		} else {
			app.OutAppID = getApp.AppID

			if checkAppChangedForLocal(getApp, app) {
				response, err := openapiapp.UpdateApp(s.api, app.OutAppID, app.Name, app.Type, app.VersionNum, app.Image, app.LicenseManagerId, app.ResidualLogParser, app.EnableResidual, app.EnableSnapshot, util.GetKeyValueMap(appDto.BinPath), update.Published)
				if err != nil {
					return err
				}
				if response == nil || response.ErrorCode != "" {
					return fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
				}
			}
		}

		app.ID = s.sid.Generate()
		err = s.appDao.AddApps(ctx, []*model.App{app})
		if err != nil {
			return err
		}

		err = AddAppPermission(ctx, common.PermissionResourceTypeLocalApp, s, app)
		if err != nil {
			return err
		}
	} else {
		if checkAppChangedSelfForLocal(originApp, app) {
			response, err := openapiapp.UpdateApp(s.api, app.OutAppID, app.Name, app.Type, app.VersionNum, app.Image, app.LicenseManagerId, app.ResidualLogParser, app.EnableResidual, app.EnableSnapshot, util.GetKeyValueMap(appDto.BinPath), update.Published)
			if err != nil {
				return err
			}
			if response == nil || response.ErrorCode != "" {
				return fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
			}
		}

		err = s.appDao.UpdateApp(ctx, app)
		if err != nil {
			return err
		}
	}

	tracelog.Info(ctx, fmt.Sprintf("update app: [%v]", serializeutil.GetStringForTraceLog(util.GetSimpleAppValue(app))))

	return nil
}

func (s *AppService) DeleteApp(ctx context.Context, name, computeType string) error {
	app := &model.App{Name: name, ComputeType: computeType}
	exist, err := s.appDao.GetApp(ctx, app)
	if err != nil {
		return err
	}
	if !exist {
		return status.Errorf(errcode.ErrAppTemplateNotExist, "app template not exist")
	}
	if app.State == common.Published {
		return status.Errorf(errcode.ErrAppTemplateHasPublished, "app template [%v] has published", name)
	}
	if app.ComputeType == common.Cloud {
		return status.Errorf(errcode.ErrAppUnableOperateCloudAppData, "cloud app template [%v] data can't operate", name)
	}

	_ = DeleteAppPermission(ctx, common.PermissionResourceTypeLocalApp, s, app)
	_, _ = openapiapp.DeleteApp(s.api, app.OutAppID)

	err = s.appDao.DeleteApp(ctx, app)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("by name: [%v], computeType: [%v], delete app: [%v]", name, computeType, serializeutil.GetStringForTraceLog(util.GetSimpleAppValue(app))))

	return nil
}

func (s *AppService) SyncAppContent(ctx context.Context, baseAppId string, syncAppIds []string) error {
	app := &model.App{ID: snowflake.MustParseString(baseAppId)}
	exist, err := s.appDao.GetApp(ctx, app)
	if err != nil {
		return err
	}
	if !exist {
		return status.Errorf(errcode.ErrAppTemplateNotExist, "app template not exist")
	}

	for _, v := range syncAppIds {
		syncApp := &model.App{ID: snowflake.MustParseString(v)}
		exist, err := s.appDao.GetApp(ctx, syncApp)
		if err != nil {
			return err
		}
		if !exist {
			return status.Errorf(errcode.ErrAppTemplateNotExist, "app template not exist")
		}

		syncApp.Content = app.Content

		_ = s.appDao.UpdateApp(ctx, syncApp)
	}

	tracelog.Info(ctx, fmt.Sprintf("sync baseAppId: [%v] content to syncAppIds: [%+v]", baseAppId, syncAppIds))

	return nil
}

func (s *AppService) PublishApp(ctx context.Context, names []string, computeType, state string) error {
	err := s.appDao.UpdateAppsState(ctx, names, computeType, state)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("publish app param: names: [%+v], computeType: [%v], state: [%v]", names, computeType, state))

	return nil
}

func (s *AppService) ListZone(ctx context.Context) ([]string, error) {
	response, err := job.ZoneList(s.api)
	if err != nil {
		return nil, err
	}
	if response == nil || response.ErrorCode != "" {
		return nil, fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	zoneList := make([]string, 0)
	if response.Data != nil && response.Data.Zones != nil {
		for i := range response.Data.Zones {
			zoneList = append(zoneList, i)
		}
	}

	return zoneList, nil
}

func (s *AppService) ListQueue(ctx context.Context, appId string) ([]*dto.QueueInfo, error) {
	queueListResponse, err := client.GetInstance().Monitor.QueueList(ctx, &monitor.QueueListRequest{})
	if err != nil {
		return nil, err
	}
	if queueListResponse == nil {
		return make([]*dto.QueueInfo, 0), nil
	}

	if strutil.IsEmpty(appId) {
		return util.ConvertQueueInfos(queueListResponse), nil
	}

	id := snowflake.MustParseString(appId)
	app := &model.App{ID: id}
	exist, err := s.appDao.GetApp(ctx, app)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, status.Errorf(errcode.ErrAppTemplateNotExist, "app template not exist")
	}
	if app.ComputeType == common.Cloud {
		return nil, status.Errorf(errcode.ErrAppCloudCompeteTypeNotSupported, "cloud compute type not support")
	}

	queryQueueNames := queueListResponse.QueueNames
	if len(app.QueueNames) != 0 {
		queryQueueNames = app.QueueNames
	}

	queueInfos, err := getQueueNamesWithCPUNumber(ctx, queryQueueNames)
	if err != nil {
		return nil, err
	}

	return queueInfos, nil
}

func (s *AppService) ListLicense(ctx context.Context) ([]*dto.LicenseInfo, error) {
	licenses, _, err := getLicenseInfo(ctx)
	if err != nil {
		return nil, err
	}

	return licenses, nil
}

func (s *AppService) CheckLicenseManagerIdUsed(ctx context.Context, licenseManagerId string) (bool, error) {
	apps, err := s.appDao.ListApp(ctx, nil, &model.App{LicenseManagerId: licenseManagerId})
	if err != nil {
		return false, err
	}

	if len(apps) > 0 {
		return true, nil
	}

	return false, nil
}

func (s *AppService) GetSchedulerResourceKey(ctx context.Context) ([]string, error) {
	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, config.GetConfig().SchedulerResourcePath)
	resources, err := config.GetResources(path)
	if err != nil {
		return nil, err
	}

	resource := make([]string, 0)
	for k := range resources {
		resource = append(resource, k)
	}

	return resource, nil
}

func (s *AppService) GetSchedulerResourceValue(ctx context.Context, appId, resourceType, resourceSubType string) ([]*dto.Item, error) {
	logger := logging.GetLogger(ctx)

	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, config.GetConfig().SchedulerResourcePath)
	resources, err := config.GetResources(path)
	if err != nil {
		return nil, err
	}

	resource := resources[resourceType]
	if resource == nil {
		return nil, status.Errorf(errcode.ErrAppGetSchedulerResourceKeyNotFound, "resource type not found")
	}

	items := make([]*dto.Item, 0)
	if resourceSubType == "" {
		resourceParamMap := make(map[string]any)
		if consts.ResourceDataTypeDynamic == resource.Type {
			err = loadSchedulerResourceParam(ctx, resourceType, resourceParamMap)
			if err != nil {
				return nil, err
			}
		}

		app := &model.App{ID: snowflake.MustParseString(appId)}
		exist, _ := s.appDao.GetApp(ctx, app)
		switch {
		case resourceType == consts.ResourceResolverTypeQueue && !exist:
			logger.Errorf("when get queue resource, app template not exist, appId: [%v]", appId)
		}
		appQueueMap := make(map[string]struct{}, 0)
		for _, v := range app.QueueNames {
			appQueueMap[v] = struct{}{}
		}

		for _, v := range resource.Data {
			item := &dto.Item{Value: v.Key}
			if consts.ResourceDataTypeDynamic == resource.Type {
				if _, ok := appQueueMap[v.Key]; resourceType == consts.ResourceResolverTypeQueue && len(appQueueMap) != 0 && !ok {
					continue
				}
				tmpSuffix := resource.KeySuffix
				err = fillSchedulerResourceParam(resourceType, resourceParamMap, v.Key, &tmpSuffix)
				if err != nil {
					return nil, err
				}
				item.Suffix = tmpSuffix
			}
			items = append(items, item)
		}
	} else {
		for _, v := range resource.Data {
			if resourceSubType == v.Key {
				for _, v := range v.Value {
					items = append(items, &dto.Item{Value: v})
				}
				break
			}
		}
	}

	return items, nil
}

func fillSchedulerResourceParam(resourceType string, resourceParamMap map[string]any, key string, keySuffix *string) error {
	switch resourceType {
	case consts.ResourceResolverTypePlatform:
		if v, ok := resourceParamMap[key]; ok {
			platformCore := v.(*monitor.PlatformCore)
			*keySuffix = fmt.Sprintf(*keySuffix, platformCore.AvailableCores, platformCore.TotalCores)
		} else {
			*keySuffix = fmt.Sprintf(*keySuffix, 0, 0)
		}
	case consts.ResourceResolverTypeQueue:
		if v, ok := resourceParamMap[key]; ok {
			queueInfo := v.(*monitor.QueueCoreInfo)
			*keySuffix = fmt.Sprintf(*keySuffix, queueInfo.AvailableCores, queueInfo.TotalCores)
		} else {
			*keySuffix = fmt.Sprintf(*keySuffix, 0, 0)
		}
	default:
		return fmt.Errorf("resource type: [%v] not match when fill scheduler resource param", resourceType)
	}

	return nil
}

func loadSchedulerResourceParam(ctx context.Context, resolver string, resourceParamMap map[string]any) error {
	switch resolver {
	case consts.ResourceResolverTypePlatform:
		platformCores, err := client.GetInstance().Monitor.GetPlatformCores(ctx, &monitor.GetPlatformCoresRequest{})
		if err != nil {
			return err
		}
		for _, v := range platformCores.PlatformCores {
			resourceParamMap[v.PlatformName] = v
		}
	case consts.ResourceResolverTypeQueue:
		queueInfos, err := client.GetInstance().Monitor.GetQueueCoreInfos(ctx, &monitor.GetQueueCoreInfosRequest{})
		if err != nil {
			return err
		}
		for _, v := range queueInfos.QueueCoreInfos {
			resourceParamMap[v.QueueName] = v
		}
	default:
		return fmt.Errorf("resource type: [%v] not match when load scheduler resource param", resolver)
	}

	return nil
}

func saveAppYaml(ctx context.Context, appName, appType, appVersion, computeType string, content []byte) error {
	yamlFilePath, err := GetAppYamlPath(appName, appType, appVersion, computeType)
	if err != nil {
		return fmt.Errorf("get app script path err: %v, computeType: [%v]", err, computeType)
	}

	err = os.WriteFile(yamlFilePath, content, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("write app template yaml err: %v", err)
	}

	return nil
}

func saveAppScript(ctx context.Context, appName, appType, appVersion, computeType, content string) error {
	scriptFilePath, err := GetAppScriptPath(appName, appType, appVersion, computeType)
	if err != nil {
		return fmt.Errorf("get app script path err: %v, computeType: [%v]", err, computeType)
	}

	err = os.WriteFile(scriptFilePath, []byte(content), os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("write app template script err: %v", err)
	}

	return nil
}

func getAppScript(ctx context.Context, appName, appType, appVersion, computeType string) (string, error) {
	scriptFilePath, err := GetAppScriptPath(appName, appType, appVersion, computeType)
	if err != nil {
		return "", fmt.Errorf("get app script path err: %v, computeType: [%v]", err, computeType)
	}

	content, err := os.ReadFile(scriptFilePath)
	if err != nil {
		return "", fmt.Errorf("read app template script err: %v", err)
	}

	return string(content), nil
}

func GetAppYamlPath(appName, appType, appVersion, computeType string) (string, error) {
	path, err := GetAppTemplatePath(appName, appType, appVersion, computeType)
	if err != nil {
		return "", err
	}

	yamlPath := fmt.Sprintf("%v%v%v", appType, common.Dot, common.Yaml)

	return filepath.Join(path, yamlPath), nil
}

func GetAppScriptPath(appName, appType, appVersion, computeType string) (string, error) {
	path, err := GetAppTemplatePath(appName, appType, appVersion, computeType)
	if err != nil {
		return "", err
	}

	scriptPath := fmt.Sprintf("%v%v%v", appType, common.Dot, consts.Cmd)

	return filepath.Join(path, scriptPath), nil
}

// GetAppTemplatePath 根据不同的计算类型获取对应的计算应用模版路径
func GetAppTemplatePath(appName, appType, appVersion, computeType string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get pwd err: %v", err)
	}

	appTemplatePath, appDirName := "", ""

	switch computeType {
	case common.Local:
		appDirName = fmt.Sprintf("%v%v%v", appType, common.Underline, appVersion)
		appTemplatePath = filepath.Join(pwd, consts.LocalAppTemplateDir)
	case common.Cloud:
		appDirName = fmt.Sprintf("%v%v%v", appType, common.Underline, appVersion)
		appTemplatePath = filepath.Join(pwd, consts.CloudAppTemplateDir)
	default:
		return "", fmt.Errorf("compute type [%v] not support", computeType)
	}

	// 只需要获取计算应用模版的存储路径
	if appName == "" && appType == "" && appVersion == "" {
		// 检测路径是否存在, 不存在则创建
		err = checkDirOrDefaultCreate(appTemplatePath)
		if err != nil {
			return "", err
		}

		return appTemplatePath, nil
	}

	// 需要获取某个计算应用存储的全路径
	appTemplatePath = filepath.Join(appTemplatePath, appDirName)

	// 检测路径是否存在, 不存在则创建
	err = checkDirOrDefaultCreate(appTemplatePath)
	if err != nil {
		return "", err
	}

	return appTemplatePath, nil
}

func checkDirOrDefaultCreate(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(path, 0755); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("app template path err: %v", err)
		}
	}
	return nil
}

func getSelectedQueueNames(queueInfos []*dto.QueueInfo) []string {
	queueNames := make([]string, 0)

	for _, v := range queueInfos {
		if v.Select {
			queueNames = append(queueNames, v.QueueName)
		}
	}

	return queueNames
}

func getSelectedLicenseInfo(licenseInfos []*dto.LicenseInfo) string {
	licenseManagerId := ""

	for _, v := range licenseInfos {
		if v.Select {
			licenseManagerId = v.Id
		}
	}

	return licenseManagerId
}

func getQueueNames(ctx context.Context) ([]*dto.QueueInfo, []string, error) {
	queueListResponse, err := client.GetInstance().Monitor.QueueList(ctx, &monitor.QueueListRequest{})
	if err != nil {
		return nil, nil, err
	}
	if queueListResponse == nil {
		return make([]*dto.QueueInfo, 0), nil, nil
	}

	return util.ConvertQueueInfos(queueListResponse), queueListResponse.QueueNames, nil
}

func getLicenseInfo(ctx context.Context) ([]*dto.LicenseInfo, []*license.LicenseType, error) {
	licenseListResponse, err := client.GetInstance().License.QueueLicenseTypeList(ctx, &license.QueueLicenseTypeListRequest{})
	if err != nil {
		return nil, nil, err
	}
	if licenseListResponse == nil || len(licenseListResponse.LicenseTypes) == 0 {
		return make([]*dto.LicenseInfo, 0), make([]*license.LicenseType, 0), nil
	}

	return util.ConvertLicenseInfos(licenseListResponse), licenseListResponse.LicenseTypes, nil
}

func getQueueNamesWithCPUNumber(ctx context.Context, queryQueueNames []string) ([]*dto.QueueInfo, error) {
	queueInfoListResponse, err := client.GetInstance().Monitor.GetQueueAvailableCores(ctx, &monitor.GetQueueAvailableCoresRequest{QueueNames: queryQueueNames})
	if err != nil {
		return nil, err
	}

	queueInfos := make([]*dto.QueueInfo, 0, len(queueInfoListResponse.QueueCores))
	for _, v := range queueInfoListResponse.QueueCores {
		queueInfos = append(queueInfos, &dto.QueueInfo{
			QueueName: v.QueueName,
			CPUNumber: v.CoreNum,
		})
	}

	return queueInfos, nil
}

func checkAndGetApp(ctx context.Context, api *openapi.OpenAPI, name string) (*schema.Application, bool, error) {
	logger := logging.GetLogger(ctx)

	response, err := openapiapp.ListApp(api)
	if err != nil {
		logger.Errorf("openapi list app err: %v", err)
		return nil, false, err
	}
	if response == nil || response.ErrorCode != "" {
		logger.Errorf("openapi response nil or response err: [%+v], ", response.Response)
		return nil, false, fmt.Errorf("openapi response nil or response err: [%+v], ", response.Response)
	}
	if response.Data == nil {
		logger.Errorf("openapi response data nil")
		return nil, false, nil
	}

	for _, v := range *response.Data {
		appName := fmt.Sprintf("%v%v%v", v.Type, common.Blank, v.Version)
		if name == appName {
			return &v, true, nil
		}
	}

	return nil, false, nil
}

func checkAppChangedForLocal(getApp *schema.Application, app *model.App) bool {
	binPathMap, err := util.ConvertBinPathMap(getApp.BinPath)
	if err != nil {
		logging.Default().Errorf("convert bin path: [%v] to bin path map err: %v", getApp.BinPath, err)
		return false
	}

	return getApp.Image != app.Image || getApp.LicManagerId != app.LicenseManagerId || !maputil.EqualMaps(binPathMap, app.BinPath) ||
		getApp.ResidualEnable != app.EnableResidual || getApp.SnapshotEnable != app.EnableSnapshot || getApp.ResidualLogParser != app.ResidualLogParser
}

func checkAppChangedSelfForLocal(originApp, app *model.App) bool {
	return originApp.Image != app.Image || originApp.LicenseManagerId != app.LicenseManagerId || !maputil.EqualMaps(originApp.BinPath, app.BinPath) ||
		originApp.EnableResidual != app.EnableResidual || originApp.EnableSnapshot != app.EnableSnapshot || originApp.ResidualLogParser != app.ResidualLogParser
}

func checkAppChangedForCloud(localApp *model.App, v *schema.Application) bool {
	return localApp.VersionNum == v.Version && localApp.Type == v.Type && localApp.EnableResidual == v.ResidualEnable && localApp.EnableSnapshot == v.SnapshotEnable &&
		localApp.ResidualLogParser == v.ResidualLogParser
}
