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
	"xorm.io/xorm"
)

// AppAllowSrv 应用配额
type AppAllowSrv interface {
	GetAllow(ctx context.Context, AppID snowflake.ID) (*schema.ApplicationAllow, error)
	AddAllow(ctx context.Context, AppID snowflake.ID) (*schema.ApplicationAllow, error)
	DeleteAllow(ctx context.Context, AppID snowflake.ID) error
}

type applicationAllowService struct {
	store store.FactoryNew
	idGen snowflake.IDGen
}

var _ AppAllowSrv = (*applicationAllowService)(nil)

func newAppAllowService(srv *service, idGen snowflake.IDGen) *applicationAllowService {
	return &applicationAllowService{
		store: srv.store,
		idGen: idGen,
	}
}

func (a *applicationAllowService) GetAllow(ctx context.Context, appID snowflake.ID) (*schema.ApplicationAllow, error) {
	logger := logging.GetLogger(ctx).With("func", "GetAllow", "AppID", appID)
	itemDao := a.store.ApplicationAllow()
	item, err := itemDao.GetByAppId(ctx, nil, appID)
	if err != nil {
		if errors.Is(err, xorm.ErrNotExist) {
			return nil, common.ErrAppAllowNotFound
		}
		logger.Warnf("get item error: %v", err)
		return nil, err
	}

	return convertApplicationAllowToSchema(item), nil
}

func (a *applicationAllowService) AddAllow(ctx context.Context, appID snowflake.ID) (*schema.ApplicationAllow, error) {
	logger := logging.GetLogger(ctx).With("func", "AddAllow", "AppID", appID)
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

		id, err := a.idGen.GenID(ctx)
		if err != nil {
			logger.Warnf("gen id error: %v", err)
			return nil, err
		}

		item := models.ApplicationAllow{
			ID:            id,
			ApplicationID: appID,
		}

		_, err = session.Insert(&item)
		if err != nil {
			if db.IsDuplicatedError(err) {
				return nil, common.ErrAppAllowAlreadyExist
			}
			logger.Warnf("insert item error: %v", err)
			return nil, err
		}
		return item, nil
	})

	if err != nil {
		if errors.Is(err, xorm.ErrNotExist) {
			return nil, common.ErrAppIDNotFound
		}
		return nil, err
	}

	item, ok := result.(models.ApplicationAllow)
	if !ok {
		logger.Warnf("invalid result, not *models.ApplicationAllow")
		return nil, errors.New("invalid result")
	}

	return convertApplicationAllowToSchema(&item), nil
}

func (a *applicationAllowService) DeleteAllow(ctx context.Context, appID snowflake.ID) error {
	logger := logging.GetLogger(ctx).With("func", "DeleteAllow", "AppID", appID)
	item := models.ApplicationAllow{
		ApplicationID: appID,
	}
	session := a.store.Engine().Context(ctx)
	defer session.Close()

	cnt, err := session.Delete(&item)
	if err != nil {
		logger.Warnf("delete item error: %v", err)
		return err
	}
	if cnt == 0 {
		return common.ErrAppAllowNotFound
	}
	return nil
}

func convertApplicationAllowToSchema(item *models.ApplicationAllow) *schema.ApplicationAllow {
	return &schema.ApplicationAllow{
		ID:    item.ID.String(),
		AppID: item.ApplicationID.String(),
	}
}
