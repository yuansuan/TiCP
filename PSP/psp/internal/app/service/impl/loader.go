package impl

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/PSP/psp/internal/app/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/util"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	openapiapp "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/app"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

type AppLoader struct {
	sid *snowflake.Node
	api *openapi.OpenAPI

	appDao dao.AppDao
}

func NewAppLoader() (*AppLoader, error) {
	node, err := snowflake.GetInstance()
	if err != nil {
		logging.Default().Errorf("new snowflake node err: %v", err)
		return nil, err
	}
	api, err := openapi.NewLocalAPI()
	if err != nil {
		logging.Default().Errorf("new openapi err: %v", err)
		return nil, err
	}
	return &AppLoader{
		sid:    node,
		api:    api,
		appDao: dao.NewAppDao(),
	}, nil
}

func (s *AppLoader) CheckInternalAppTemplate(ctx context.Context) {
	logger := logging.GetLogger(ctx)

	starccmApp := &model.App{ID: consts.InternalTemplateStarCCMId}
	exist, err := s.appDao.GetApp(ctx, starccmApp)
	if err != nil {
		logger.Errorf("get internal starccm app template err: [%v]", err)
		return
	}

	if !exist || starccmApp.OutAppID != "" {
		return
	}

	response, err := openapiapp.AddApp(s.api, starccmApp.Name, starccmApp.Type, starccmApp.VersionNum, starccmApp.Image, starccmApp.LicenseManagerId, starccmApp.ResidualLogParser, starccmApp.EnableResidual, starccmApp.EnableSnapshot, starccmApp.BinPath)
	if err != nil {
		logger.Errorf("add internal starccm app template err: [%v]", err)
		return
	}
	if response == nil || response.ErrorCode != "" {
		logger.Errorf("openapi response nil or response err: [%+v]", response.Response)
		return
	}
	starccmApp.OutAppID = response.Data.AppID

	err = s.appDao.UpdateApp(ctx, starccmApp)
	if err != nil {
		logger.Errorf("update internal starccm app template err: [%v]", err)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("internal starccm app: [%v] template compensate success, outAppID: [%v]", starccmApp.ID, starccmApp.OutAppID))
}

func (s *AppLoader) InitAppTemplates(ctx context.Context) {
	logger := logging.GetLogger(ctx)
	apps, templatePath, err := loadInternalTemplate(ctx)
	if err != nil {
		logger.Errorf("load internal template err: %v", err)
		return
	}
	if len(apps) == 0 {
		logger.Infof("no config app template in path: [%v]", templatePath)
		return
	}

	var appList []*model.App
	for _, v := range apps {
		if v.ComputeType == common.Cloud {
			continue
		}
		app, err := util.ConvertAppModel(v)
		if err != nil {
			logger.Errorf("convert app model err: %v", err)
			continue
		}
		app.ComputeType = common.Local
		appList = append(appList, app)
	}

	app := &model.App{ComputeType: common.Local}
	dbApps, err := s.appDao.ListApp(ctx, nil, app)
	if err != nil {
		logger.Errorf("get app list err: %v", err)
		return
	}

	newApps := compareAppTemplate(appList, dbApps)
	if len(newApps) == 0 {
		logger.Infof("no new app template")
		return
	}

	nameAppMap, err := getNameAppMap(ctx, s.api)
	if err != nil {
		logger.Errorf("get openapi app list err: %v", err)
		return
	}

	for _, v := range newApps {
		v.ID = s.sid.Generate()

		if appData, ok := nameAppMap[v.Name]; !ok {
			response, err := openapiapp.AddApp(s.api, v.Name, v.Type, v.VersionNum, "", "", "", false, false, nil)
			if err != nil {
				logger.Errorf("openapi add app err: %v, appName: [%v]", err, v.Name)
				continue
			}
			if response == nil || response.ErrorCode != "" {
				logger.Errorf("openapi response nil or response err: [%+v], ", response.Response)
				continue
			}
			if response.Data == nil {
				logger.Errorf("openapi response data nil")
				continue
			}

			v.OutAppID = response.Data.AppID
		} else {
			v.OutAppID = appData.AppID
			v.Image = appData.Image
			v.LicenseManagerId = appData.LicManagerId

			binPathMap, err := util.ConvertBinPathMap(appData.BinPath)
			if err != nil {
				logger.Errorf("convert bin path: [%v] to bin path map err: %v", appData.BinPath, err)
				continue
			}
			v.BinPath = binPathMap
		}

		err = s.appDao.AddApps(ctx, []*model.App{v})
		if err != nil {
			logger.Errorf("add app: [%v] err: %v", err, v.Name)
			continue
		}

		err = AddAppPermission(ctx, common.PermissionResourceTypeLocalApp, s, v)
		if err != nil {
			continue
		}
	}
}

func loadInternalTemplate(ctx context.Context) ([]*dto.App, string, error) {
	logger := logging.GetLogger(ctx)

	appTemplatePath, err := GetAppTemplatePath("", "", "", common.Local)
	if err != nil {
		logger.Errorf("get app template path err: %v, computeType: [%v]", err, common.Local)
		return nil, "", err
	}

	dirs, err := os.ReadDir(appTemplatePath)
	if err != nil {
		logger.Errorf("read app template dir err: %v", err)
		return nil, "", err
	}

	apps := make([]*dto.App, 0)
	for _, dir := range dirs {
		if dir.IsDir() {
			dirName := dir.Name()
			underlineIndex := strings.LastIndex(dirName, common.Underline)
			appFilePath := filepath.Join(appTemplatePath, dirName)
			app, err := readApplicationConfig(ctx, appFilePath, dirName[:underlineIndex])
			if err != nil {
				logger.Errorf("find app err: %v, params: appFilePath: [%v], appName: [%v]", err, appFilePath, dirName[:underlineIndex])
				return nil, "", err
			}
			app.Icon = loadAppIcon(ctx, appFilePath, dirName[:underlineIndex])
			app.Version = strings.TrimSpace(dirName[underlineIndex+1:])
			apps = append(apps, app)
		}
	}
	return apps, appTemplatePath, nil
}

func compareAppTemplate(appList, dbApps []*model.App) []*model.App {
	if len(dbApps) == 0 {
		return appList
	}
	dbListMap := make(map[string]*model.App)
	for _, app := range dbApps {
		dbListMap[app.Name] = app
	}

	newApps := make([]*model.App, 0)
	for _, app := range appList {
		if _, ok := dbListMap[app.Name]; ok {
			delete(dbListMap, app.Name)
		} else {
			newApps = append(newApps, app)
		}
	}

	return newApps
}

func loadAppIcon(ctx context.Context, templatePath, appName string) string {
	logger := logging.GetLogger(ctx)
	iconPath := filepath.Join(templatePath, fmt.Sprintf("%v%v%v", appName, common.Dot, common.Png))
	_, err := os.Stat(iconPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Infof("template icon not exist: [%v]", iconPath)
			return ""
		}
		logger.Errorf("get template icon err: %v, iconFile: [%v]", err, iconPath)
		return ""
	}

	content, err := os.ReadFile(iconPath)
	if err != nil {
		logger.Infof("read template icon err: %v, param: iconFile: [%v]", err, iconPath)
		return ""
	}
	return fmt.Sprintf("%v%v", consts.IconDataBase64Prefix, base64.StdEncoding.EncodeToString(content))
}

func readApplicationConfig(ctx context.Context, configPath, configFileName string) (*dto.App, error) {
	logger := logging.GetLogger(ctx)
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configFileName)
	v.SetConfigType(common.Yaml)
	err := v.ReadInConfig()
	if err != nil {
		logger.Errorf("read app template config err: %v", err)
		return nil, err
	}

	md := mapstructure.Metadata{}
	app := &dto.App{}
	err = v.Unmarshal(app, func(config *mapstructure.DecoderConfig) {
		config.TagName = "json"
		config.Metadata = &md
	})
	if err != nil {
		logger.Errorf("unmarshal app template config err: %v", err)
		return nil, err
	}
	return app, nil
}

func getNameAppMap(ctx context.Context, api *openapi.OpenAPI) (map[string]*schema.Application, error) {
	logger := logging.GetLogger(ctx)

	nameAppMap := make(map[string]*schema.Application)
	response, err := openapiapp.ListApp(api)
	if err != nil {
		logger.Errorf("openapi list app err: %v", err)
		return nil, err
	}
	if response == nil || response.ErrorCode != "" {
		logger.Errorf("openapi response nil or response err: [%+v], ", response.Response)
		return nil, fmt.Errorf("openapi response nil or response err: [%+v], ", response.Response)
	}
	if response.Data == nil {
		logger.Errorf("openapi response data nil")
		return nameAppMap, nil
	}

	for _, v := range *response.Data {
		appName := fmt.Sprintf("%v%v%v", v.Type, common.Blank, v.Version)
		nameAppMap[appName] = &v
	}

	return nameAppMap, nil
}
