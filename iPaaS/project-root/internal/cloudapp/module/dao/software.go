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

// ListSoftwareParams 列举所有软件的参数
type ListSoftwareParams struct {
	Name       string
	Platform   string
	PageOffset int
	PageSize   int
	Zone       zone.Zone
}

// ListSoftware 根据查询条件找到对应软件
func ListSoftware(ctx context.Context, params *ListSoftwareParams) (list []*models.Software, total int64, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		total, err = exorm.New(db).
			Cond(params.Name != "", "name like ?", "%"+params.Name+"%").
			Cond(params.Zone != "", "zone = ?", params.Zone).
			Cond(params.Platform != "", "platform = ?", params.Platform).
			Raw().
			Limit(params.PageSize, params.PageOffset).
			OrderBy("id desc").
			FindAndCount(&list)
		return errors.Wrap(err, "list")
	})
	return
}

func ListSoftwareByUser(ctx context.Context, params *ListSoftwareParams, userId snowflake.ID) (list []*models.Software, total int64, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		db.Table(_Software.TableName()).Alias("s").
			Join("INNER", []string{_SoftwareUser.TableName(), "su"}, "s.id = su.software_id").
			Where("su.user_id = ?", userId)

		total, err = exorm.New(db).
			Cond(params.Name != "", "s.name like ?", "%"+params.Name+"%").
			Cond(params.Zone != "", "s.zone = ?", params.Zone).
			Cond(params.Platform != "", "s.platform = ?", params.Platform).
			Raw().
			Limit(params.PageSize, params.PageOffset).
			OrderBy("s.id desc").
			FindAndCount(&list)
		return errors.Wrap(err, "list")
	})
	return
}

// GetSoftware 获取一个软件
func GetSoftware(ctx context.Context, id snowflake.ID) (software *models.Software, exist bool, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		software = &models.Software{Id: id}
		exist, err = db.ID(id).Get(software)
		if err != nil {
			return errors.Wrap(err, "dao")
		}
		return nil
	})

	if err != nil {
		return nil, false, err
	}

	return software, exist, nil
}

func GetSoftwareByUser(ctx context.Context, softwareId, userId snowflake.ID) (*models.Software, bool, error) {
	exist := false
	software := new(models.Software)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		var e error
		exist, e = db.Table(_Software.TableName()).Alias("s").
			Join("INNER", []string{_SoftwareUser.TableName(), "su"}, "s.id = su.software_id").
			Where("s.id = ?", softwareId).
			Where("su.user_id = ?", userId).
			Get(software)
		return e
	})
	if err != nil {
		return nil, false, err
	}

	return software, exist, nil
}

// AddSoftware 添加软件
func AddSoftware(ctx context.Context, software *models.Software) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Insert(software)
		return errors.Wrap(err, "insert")
	})
}

// UpdateSoftware 更新软件
func UpdateSoftware(ctx context.Context, software *models.Software, mustCol ...string) (bool, error) {
	var count int64
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		var e error
		count, e = db.ID(software.Id).Cols(mustCol...).Update(software)
		return errors.Wrap(e, "update")
	})

	return count > 0, err
}

// UpdateSoftwareAllCol 更新软件(所有字段)
func UpdateSoftwareAllCol(ctx context.Context, software *models.Software) (bool, error) {
	var count int64
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		var e error
		count, e = db.ID(software.Id).Cols("zone", "name", "desc", "icon", "platform", "image_id", "init_script", "gpu_desired", "update_time").Update(software)
		return errors.Wrap(e, "update")
	})

	return count > 0, err
}

// DeleteSoftware 删除软件
func DeleteSoftware(ctx context.Context, id snowflake.ID) (count int64, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		count, err = db.ID(id).Delete(&models.Software{})
		return errors.Wrap(err, "delete")
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func BatchCheckSoftwaresExist(ctx context.Context, softwares []snowflake.ID) ([]snowflake.ID, error) {
	existList := make([]*models.Software, 0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		e := db.Cols("id").
			In("id", softwares).
			Find(&existList)
		return e
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

func BatchAddSoftwareUsers(ctx context.Context, softwares, users []snowflake.ID) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		softwareUsers := make([]*models.SoftwareUser, 0)
		for _, software := range softwares {
			for _, user := range users {
				softwareUsers = append(softwareUsers, &models.SoftwareUser{
					SoftwareId: software,
					UserId:     user,
				})
			}
		}

		_, err := db.InsertMulti(&softwareUsers)
		if db2.IsDuplicatedError(err) {
			return db2.ErrDuplicatedEntry
		}

		return err
	})
}

func BatchDeleteSoftwareUsersBySoftwareId(ctx context.Context, softwares []snowflake.ID) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.In("software_id", softwares).Delete(new(models.SoftwareUser))
		return err
	})
}

func BatchDeleteSoftwareUsers(ctx context.Context, softwares, users []snowflake.ID) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.In("software_id", softwares).
			In("user_id", users).
			Delete(new(models.SoftwareUser))
		return err
	})
}

func SoftwareUserExist(ctx context.Context, softwareId, userId snowflake.ID) (bool, error) {
	exist := false
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		var e error
		exist, e = db.Where("software_id = ?", softwareId).
			Where("user_id = ?", userId).
			Exist(new(models.SoftwareUser))
		return e
	})

	return exist, err
}

func GetSofwareUserByUsers(ctx context.Context, userIds []snowflake.ID) ([]models.SoftwareUser, error) {
	list := make([]models.SoftwareUser, 0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		return db.In("user_id", userIds).Find(&list)
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}
