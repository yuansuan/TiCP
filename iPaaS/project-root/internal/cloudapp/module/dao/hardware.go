package dao

import (
	"context"

	"github.com/pkg/errors"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils/exorm"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	db2 "github.com/yuansuan/ticp/iPaaS/project-root/pkg/db"
)

// ListHardwareParams 根据查询条件找到对应硬件的参数
type ListHardwareParams struct {
	Name       string
	Cpu        int
	Mem        int
	Gpu        int
	PageOffset int
	PageSize   int
	Zone       zone.Zone
}

// ListHardware 根据查询条件找到对应硬件
func ListHardware(ctx context.Context, params *ListHardwareParams) (list []*models.Hardware, total int64, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		total, err = exorm.New(db).
			Cond(params.Name != "", "name like ?", "%"+params.Name+"%").
			Cond(params.Zone != "", "zone = ?", params.Zone).
			Cond(params.Cpu != 0, "cpu = ?", params.Cpu).
			Cond(params.Mem != 0, "mem = ?", params.Mem).
			Cond(params.Gpu != 0, "gpu = ?", params.Gpu).
			Raw().
			Limit(params.PageSize, params.PageOffset).
			OrderBy("id desc").
			FindAndCount(&list)
		return errors.Wrap(err, "dao")
	})
	return
}

func ListHardwareByUser(ctx context.Context, params *ListHardwareParams, userId snowflake.ID) (list []*models.Hardware, total int64, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		db.Table(_Hardware.TableName()).Alias("h").
			Join("INNER", []string{_HardwareUser.TableName(), "hu"}, "h.id = hu.hardware_id").
			Where("hu.user_id = ?", userId)

		total, err = exorm.New(db).
			Cond(params.Name != "", "h.name like ?", "%"+params.Name+"%").
			Cond(params.Zone != "", "h.zone = ?", params.Zone).
			Cond(params.Cpu != 0, "h.cpu = ?", params.Cpu).
			Cond(params.Mem != 0, "h.mem = ?", params.Mem).
			Cond(params.Gpu != 0, "h.gpu = ?", params.Gpu).
			Raw().
			Limit(params.PageSize, params.PageOffset).
			OrderBy("h.id desc").
			FindAndCount(&list)
		return errors.Wrap(err, "dao")
	})
	return
}

// GetHardware 根据ID找到对应硬件
func GetHardware(ctx context.Context, id snowflake.ID) (hardware *models.Hardware, exist bool, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		hardware = &models.Hardware{Id: id}
		exist, err = db.ID(id).Get(hardware)
		if err != nil {
			return errors.Wrap(err, "dao")
		}
		return nil
	})

	if err != nil {
		return nil, false, err
	}

	return hardware, exist, nil
}

func GetHardwareByUser(ctx context.Context, hardwareId, userId snowflake.ID) (*models.Hardware, bool, error) {
	exist, hardware := false, new(models.Hardware)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		var e error
		exist, e = db.Table(_Hardware.TableName()).Alias("h").
			Join("INNER", []string{_HardwareUser.TableName(), "hu"}, "h.id = hu.hardware_id").
			Where("h.id = ?", hardwareId).
			Where("hu.user_id = ?", userId).
			Get(hardware)
		return e
	})
	if err != nil {
		return nil, false, err
	}

	return hardware, exist, nil
}

// AddHardware 添加硬件
func AddHardware(ctx context.Context, hardware *models.Hardware) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Insert(hardware)
		return errors.Wrap(err, "insert")
	})
}

// UpdateHardware 更新硬件
func UpdateHardware(ctx context.Context, id snowflake.ID, hardware *models.Hardware, mustCol ...string) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(id).Cols(mustCol...).Update(hardware)
		return errors.Wrap(err, "edit")
	})
}

// UpdateHardwareAllCol 更新硬件(所有字段)
func UpdateHardwareAllCol(ctx context.Context, id snowflake.ID, hardware *models.Hardware) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(id).
			Cols("zone", "name", "desc", "instance_type", "instance_family", "network", "cpu", "cpu_model", "mem", "gpu", "gpu_model", "update_time").
			Update(hardware)
		return errors.Wrap(err, "edit")
	})
}

// DeleteHardware 删除硬件
func DeleteHardware(ctx context.Context, id snowflake.ID) (count int64, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		count, err = db.ID(id).Delete(&models.Hardware{})
		return errors.Wrap(err, "delete")
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func BatchCheckHardwaresExist(ctx context.Context, hardwares []snowflake.ID) ([]snowflake.ID, error) {
	existList := make([]*models.Hardware, 0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		return db.Cols("id").
			In("id", hardwares).
			Find(&existList)
	})
	if err != nil {
		return nil, errors.Wrap(err, "dao")
	}

	res := make([]snowflake.ID, 0)
	for _, v := range existList {
		res = append(res, v.Id)
	}

	return res, nil
}

func BatchAddHardwareUsers(ctx context.Context, hardwares, users []snowflake.ID) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		hardwareUsers := make([]*models.HardwareUser, 0)
		for _, hardware := range hardwares {
			for _, user := range users {
				hardwareUsers = append(hardwareUsers, &models.HardwareUser{
					HardwareId: hardware,
					UserId:     user,
				})
			}
		}

		_, err := db.InsertMulti(&hardwareUsers)
		if db2.IsDuplicatedError(err) {
			return db2.ErrDuplicatedEntry
		}

		return err
	})
}

func BatchDeleteHardwareUsersByHardwareId(ctx context.Context, hardwares []snowflake.ID) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.In("hardware_id", hardwares).Delete(new(models.HardwareUser))
		return err
	})
}

func BatchDeleteHardwareUsers(ctx context.Context, hardwares, users []snowflake.ID) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.In("hardware_id", hardwares).
			In("user_id", users).
			Delete(new(models.HardwareUser))
		return err
	})
}

func HardwareUserExist(ctx context.Context, hardwareId, userId snowflake.ID) (bool, error) {
	exist := false
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		var e error
		exist, e = db.Where("hardware_id = ?", hardwareId).
			Where("user_id = ?", userId).
			Exist(new(models.HardwareUser))
		return e
	})

	return exist, err
}

func GetHardwareUserByUsers(ctx context.Context, userIds []snowflake.ID) ([]models.HardwareUser, error) {
	list := make([]models.HardwareUser, 0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		return db.In("user_id", userIds).Find(&list)
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}
