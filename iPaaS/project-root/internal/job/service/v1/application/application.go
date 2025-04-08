package application

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/add"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/rpc"
	"xorm.io/xorm"
)

type AppSrv interface {
	GetApp(ctx context.Context, id snowflake.ID) (*models.Application, error)
	ListApps(ctx context.Context, userID snowflake.ID, publishStatus consts.PublishStatus) ([]*models.Application, int64, error)
	AddApp(ctx context.Context, appInfo *add.Request) (snowflake.ID, error)
	UpdateApp(ctx context.Context, appInfo *update.Request) error
	DeleteApp(ctx context.Context, id snowflake.ID) error
}

type applicationService struct {
	store store.FactoryNew
}

var _ AppSrv = (*applicationService)(nil)

func newAppService(srv *service) *applicationService {
	return &applicationService{store: srv.store}
}

func (s *applicationService) GetApp(ctx context.Context, id snowflake.ID) (*models.Application, error) {
	return s.store.Applications().GetApp(ctx, id)
}

// outside to inside
func convertAddApplication2App(application *add.Request) *models.Application {
	licMId := snowflake.ID(0)
	if application.LicManagerId != "" {
		licMId = snowflake.MustParseString(application.LicManagerId)
	}
	var q string
	if len(application.SpecifyQueue) > 0 {
		jsonBytes, _ := json.Marshal(application.SpecifyQueue)
		q = string(jsonBytes)
	}

	return &models.Application{
		Name:               application.Name,
		Type:               application.Type,
		Version:            application.Version,
		AppParamsVersion:   application.AppParamsVersion,
		Image:              application.Image,
		Endpoint:           application.Endpoint,
		Command:            application.Command,
		Description:        application.Description,
		IconUrl:            application.IconUrl,
		CoresMaxLimit:      application.CoresMaxLimit,
		CoresPlaceholder:   application.CoresPlaceholder,
		FileFilterRule:     application.FileFilterRule,
		ResidualEnable:     application.ResidualEnable,
		ResidualLogRegexp:  application.ResidualLogRegexp,
		ResidualLogParser:  application.ResidualLogParser,
		MonitorChartEnable: application.MonitorChartEnable,
		MonitorChartRegexp: application.MonitorChartRegexp,
		MonitorChartParser: application.MonitorChartParser,
		SnapshotEnable:     application.SnapshotEnable,
		BinPath:            application.BinPath,
		ExtentionParams:    application.ExtentionParams,
		LicManagerId:       licMId,
		NeedLimitCore:      application.NeedLimitCore,
		SpecifyQueue:       q,
	}
}

func (s *applicationService) ListApps(ctx context.Context, userID snowflake.ID, publishStatus consts.PublishStatus) ([]*models.Application, int64, error) {
	apps, size, err := s.store.Applications().ListApps(ctx, userID, publishStatus)
	if err != nil {
		return nil, 0, errors.Wrap(err, "db list apps failed")
	}
	return apps, size, nil
}

func (s *applicationService) AddApp(ctx context.Context, appInfo *add.Request) (snowflake.ID, error) {
	id, _, err := GenAppID(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "generate app id failed")
	}
	app := convertAddApplication2App(appInfo)
	app.ID = id
	app.PublishStatus = string(update.Unpublished)
	err = s.store.Applications().AddApp(ctx, app)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GenAppID(ctx context.Context) (snowflake.ID, string, error) {
	logger := logging.GetLogger(ctx).With("func", "jobcreate.genJobID")

	appID, err := rpc.GetInstance().GenID(ctx)
	if err != nil {
		logger.Errorf("generate a snowslake id fail")
		return 0, "", err
	}
	return appID, appID.String(), nil
}

func (s *applicationService) UpdateApp(ctx context.Context, appInfo *update.Request) error {
	app := convertAddApplication2App(&appInfo.Request)
	app.ID = snowflake.MustParseString(appInfo.AppID)
	app.PublishStatus = string(appInfo.PublishStatus)
	err := s.store.Applications().UpdateApp(ctx, app)
	return err
}

func (s *applicationService) DeleteApp(ctx context.Context, id snowflake.ID) error {
	engine := s.store.Engine()
	_, err := engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		session = session.Context(ctx)

		deleteRow, err := session.ID(id).Delete(&models.Application{ID: id})
		if err != nil {
			return nil, err
		}
		if deleteRow == 0 {
			return nil, xorm.ErrNotExist
		}

		// 删除关联所有应用配额
		_, err = session.Where("application_id = ?", id).Delete(&models.ApplicationQuota{
			ApplicationID: id,
		})
		return nil, err
	})

	return err
}
