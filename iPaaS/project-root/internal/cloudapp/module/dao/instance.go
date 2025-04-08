package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
)

var (
	ErrIbvInstanceNotFound = errors.New("instance not found")
)

func GetInstance(ctx context.Context, id snowflake.ID) (*models.Instance, error) {
	var instance models.Instance
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		if exists, err := db.ID(id).Get(&instance); err != nil {
			return errors.Wrap(err, "dao")
		} else if !exists {
			return errors.Wrap(ErrIbvInstanceNotFound, "dao")
		}
		return nil
	})
	return &instance, err
}

func GetInstanceByInstanceId(ctx context.Context, instanceId string) (*models.Instance, error) {
	var instance models.Instance
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		exists, err := db.Where("instance_id = ?", instanceId).Get(&instance)
		if err != nil {
			return errors.Wrap(err, "dao")
		} else if !exists {
			return errors.Wrap(ErrIbvInstanceNotFound, "dao")
		}
		return nil
	})
	return &instance, err
}

func ListRecentUnterminatedInstance(ctx context.Context, t time.Duration) (list []*models.Instance, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		err = db.
			Where("update_time < ?", time.Now().Add(-t)).
			Where("instance_status != ?", models.InstanceTerminated).
			Limit(100).
			Find(&list)
		return errors.Wrap(err, "dao")
	})
	return
}

func NewInstance(ctx context.Context, instance *models.Instance) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Insert(instance)
		return errors.Wrap(err, "dao")
	})
}

func InstanceCreated(ctx context.Context, id snowflake.ID, instanceId string) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(id).
			Where("instance_status = ?", models.InstancePending).
			Cols("instance_id", "start_time", "instance_status").
			Update(&models.Instance{
				InstanceId:     instanceId,
				StartTime:      utils.PNow(),
				InstanceStatus: models.InstanceCreated,
			})
		return errors.Wrap(err, "dao")
	})
}

func InstanceTerminated(ctx context.Context, id snowflake.ID) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(id).
			Cols("instance_status", "end_time").
			Update(&models.Instance{
				InstanceStatus: models.InstanceTerminated,
				EndTime:        utils.PNow(),
			})
		return errors.Wrap(err, "dao")
	})
}

func UpdateInstanceStatus(ctx context.Context, id snowflake.ID, status models.InstanceStatus) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(id).
			Cols("instance_status").
			Update(&models.Instance{
				InstanceStatus: status,
			})
		return errors.Wrap(err, "dao")
	})
}

func UpdateInstanceData(ctx context.Context, id snowflake.ID, data string) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(id).Cols("instance_data").
			Update(&models.Instance{
				InstanceData: data,
			})
		return errors.Wrap(err, "dao")
	})
}

func UpdateInstance(ctx context.Context, instance *models.Instance, cols []string) error {
	if instance == nil || instance.Id == snowflake.ID(0) {
		return fmt.Errorf("invalid instance")
	}

	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(instance.Id).
			Cols(cols...).
			Update(instance)

		return errors.Wrap(err, "dao")
	})
}

func IsBootVolumeOccupied(ctx context.Context, bootVolumeId string) (bool, snowflake.ID, error) {
	var exist bool
	sessionId := snowflake.ID(0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		var e error
		sessWithIns := new(models.SessionWithInstance)
		exist, e = db.Table(_Instance.TableName()).Alias("i").
			Join("INNER", []string{_Session.TableName(), "s"}, "i.id = s.instance_id").
			Cols("s.id").
			Where("i.boot_volume_id = ?", bootVolumeId).
			Where("i.instance_status != ?", models.InstanceTerminated).
			Where("s.status != ?", schema.SessionClosed).
			Get(sessWithIns)
		if e != nil {
			return fmt.Errorf("dao, %w", e)
		}
		if !exist {
			return nil
		}

		sessionId = sessWithIns.Session.Id
		return nil
	})
	if err != nil {
		return false, snowflake.ID(0), err
	}

	return exist, sessionId, nil
}
