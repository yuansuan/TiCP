package application

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/db"
	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
	"xorm.io/xorm"
)

// AppQuotaSrv 应用配额
type AppQuotaSrv interface {
	GetQuota(ctx context.Context, AppID, UserID snowflake.ID) (*schema.ApplicationQuota, error)
	AddQuota(ctx context.Context, AppID, UserID snowflake.ID) (*schema.ApplicationQuota, error)
	DeleteQuota(ctx context.Context, AppID, UserID snowflake.ID) error
}

// Helper 辅助函数
type Helper interface {
	snowflake.IDGen
	UserGeter
}

// UserGeter 获取用户信息
type UserGeter interface {
	GetUser(ctx context.Context, UserID string) (*hydra_lcp.UserInfo, error)
}

type applicationQuotaService struct {
	store  store.FactoryNew
	helper Helper
}

var _ AppQuotaSrv = (*applicationQuotaService)(nil)

func newAppQuotaService(srv *service, helper Helper) *applicationQuotaService {
	return &applicationQuotaService{
		store:  srv.store,
		helper: helper,
	}
}

func (a *applicationQuotaService) GetQuota(ctx context.Context, appID snowflake.ID, userID snowflake.ID) (*schema.ApplicationQuota, error) {
	logger := logging.GetLogger(ctx).With("func", "GetQuota", "AppID", appID, "UserID", userID)
	quotaDao := a.store.ApplicationQuota()
	quota, err := quotaDao.GetByUser(ctx, nil, appID, userID, false)
	if err != nil {
		if errors.Is(err, xorm.ErrNotExist) {
			return nil, common.ErrAppQuotaNotFound
		}
		logger.Warnf("get quota error: %v", err)
		return nil, err
	}

	return convertApplicationQuotaToSchema(quota), nil
}

func (a *applicationQuotaService) AddQuota(ctx context.Context, appID snowflake.ID, userID snowflake.ID) (*schema.ApplicationQuota, error) {
	logger := logging.GetLogger(ctx).With("func", "AddQuota", "AppID", appID, "UserID", userID)
	engine := a.store.Engine()
	result, err := engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		session = session.Context(ctx)

		app := &models.Application{}
		get, err := session.ID(appID).ForUpdate().Get(app)
		if err != nil {
			logger.Warnf("get application error: %v", err)
			return nil, err
		}
		if !get {
			return nil, xorm.ErrNotExist
		}

		// 检查userID
		user, err := a.helper.GetUser(ctx, userID.String())
		if err != nil {
			logger.Warnf("get user info error: %v", err)
			return nil, err
		}
		if user == nil {
			return nil, common.ErrUserNotExists
		}

		id, err := a.helper.GenID(ctx)
		if err != nil {
			logger.Warnf("gen id error: %v", err)
			return nil, err
		}

		quota := models.ApplicationQuota{
			ID:            id,
			ApplicationID: appID,
			YsID:          userID,
		}

		_, err = session.Insert(&quota)
		if err != nil {
			if db.IsDuplicatedError(err) {
				return nil, common.ErrAppQuotaAlreadyExist
			}
			logger.Warnf("insert quota error: %v", err)
			return nil, err
		}
		return quota, nil
	})

	if err != nil {
		if errors.Is(err, xorm.ErrNotExist) {
			return nil, common.ErrAppIDNotFound
		}
		return nil, err
	}

	quota, ok := result.(models.ApplicationQuota)
	if !ok {
		logger.Warnf("invalid result, not *models.ApplicationQuota")
		return nil, errors.New("invalid result")
	}

	return convertApplicationQuotaToSchema(&quota), nil
}

func (a *applicationQuotaService) DeleteQuota(ctx context.Context, appID snowflake.ID, userID snowflake.ID) error {
	logger := logging.GetLogger(ctx).With("func", "DeleteQuota", "AppID", appID, "UserID", userID)
	quota := models.ApplicationQuota{
		ApplicationID: appID,
		YsID:          userID,
	}
	session := a.store.Engine().Context(ctx)
	defer session.Close()

	cnt, err := session.Delete(&quota)
	if err != nil {
		logger.Warnf("delete quota error: %v", err)
		return err
	}
	if cnt == 0 {
		return common.ErrAppQuotaNotFound
	}
	return nil
}

func convertApplicationQuotaToSchema(quota *models.ApplicationQuota) *schema.ApplicationQuota {
	return &schema.ApplicationQuota{
		ID:     quota.ID.String(),
		AppID:  quota.ApplicationID.String(),
		UserID: quota.YsID.String(),
	}
}
