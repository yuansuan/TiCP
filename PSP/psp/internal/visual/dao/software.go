package dao

import (
	"context"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"strconv"
	. "xorm.io/builder"
)

type softwareDaoImpl struct{}

func NewSoftwareDao() SoftwareDao {
	return &softwareDaoImpl{}
}

func (d *softwareDaoImpl) ListSoftware(ctx context.Context, IDs []snowflake.ID, name, platform, state string, offset, limit int) ([]*model.Software, int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var softwares []*model.Software
	if len(IDs) > 0 {
		session.In("id", IDs)
	}
	if name != "" {
		session.Where("name like ?", "%"+name+"%")
	}
	if platform != "" {
		session.Where("platform = ?", platform)
	}
	if state != "" {
		session.Where("state = ?", state)
	}
	if offset > 0 {
		session.Limit(limit, (offset-1)*limit)
	} else {
		session.Limit(limit)
	}

	total, err := session.Select("*").Desc("create_time").FindAndCount(&softwares)
	if err != nil {
		return nil, 0, err
	}
	return softwares, total, nil
}

func (d *softwareDaoImpl) InsertSoftware(ctx context.Context, software *model.Software) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.Insert(software)
	if err != nil {
		return err
	}
	return nil
}

func (d *softwareDaoImpl) UpdateSoftware(ctx context.Context, software *model.Software) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	_, err := session.ID(software.ID).Omit("id, out_software_id, state").AllCols().Update(software)
	if err != nil {
		return err
	}

	return nil
}

func (d *softwareDaoImpl) DeleteSoftware(ctx context.Context, softwareID snowflake.ID) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.ID(softwareID).Delete(&model.Software{})
	if err != nil {
		return err
	}
	return nil
}

func (d *softwareDaoImpl) PublishSoftware(ctx context.Context, id snowflake.ID, state string) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	app := &model.Software{State: state}

	_, err := session.ID(id).Cols("state").Update(app)
	if err != nil {
		return err
	}
	return nil
}

func (d *softwareDaoImpl) GetSoftware(ctx context.Context, softwareID snowflake.ID, name string) (*model.Software, bool, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	software := &model.Software{}
	if int64(softwareID) > 0 {
		session.Where("id = ?", softwareID)
	}
	if name != "" {
		session.Where("name = ?", name)
	}

	exist, err := session.Get(software)
	if err != nil {
		return nil, exist, err
	}
	return software, exist, nil
}

type UsingStatusesResponse struct {
	*model.Software `xorm:"extends"`
	*model.Session  `xorm:"extends"`
}

func (d *softwareDaoImpl) UsingStatuses(ctx context.Context, username string, softwareIds []snowflake.ID) ([]*UsingStatusesResponse, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var response []*UsingStatusesResponse

	subSqlBuilder := Select("t.software_id, max(t.update_time) as max_datetime").From("visual_session t").GroupBy("t.software_id")
	if len(softwareIds) != 0 {
		softwareIdStrList := make([]string, 0, len(softwareIds))
		for _, v := range softwareIds {
			id := strconv.FormatInt(int64(v), 10)
			softwareIdStrList = append(softwareIdStrList, id)
		}
		subSqlBuilder.Where(In("t.software_id", softwareIdStrList))
	}
	if username != "" {
		subSqlBuilder.Where(Eq{"t.user_name": username})
	}
	subSql, err := subSqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	session.Table(model.SoftwareTableName).Alias("s1").
		Join("left", common.LeftParentheses+subSql+common.RightParentheses+" as sub", "s1.id = sub.software_id").
		Join("left", model.SessionTableName+" as s2", "s1.id = s2.software_id and sub.max_datetime = s2.update_time").
		Where("s1.state = ?", common.Published)

	err = session.Select("s1.*, s2.*").Find(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *softwareDaoImpl) ListRemoteApp(ctx context.Context, offset, limit int) ([]*model.RemoteApp, int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	remoteApps := []*model.RemoteApp{}
	if offset > 0 {
		session.Limit(int(limit), int((offset-1)*limit))
	} else {
		session.Limit(int(limit))
	}
	total, err := session.Select("*").FindAndCount(&remoteApps)
	if err != nil {
		return nil, 0, err
	}
	return remoteApps, total, nil
}

func (d *softwareDaoImpl) InsertRemoteApp(ctx context.Context, remoteApp *model.RemoteApp) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.Insert(remoteApp)
	if err != nil {
		return err
	}
	return nil
}

func (d *softwareDaoImpl) UpdateRemoteApp(ctx context.Context, remoteApp *model.RemoteApp) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.ID(remoteApp.ID).UseBool().Update(remoteApp)
	if err != nil {
		return err
	}
	return nil
}

func (d *softwareDaoImpl) DeleteRemoteApp(ctx context.Context, remoteAppID snowflake.ID) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.ID(remoteAppID).Delete(&model.RemoteApp{})
	if err != nil {
		return err
	}
	return nil
}

func (d *softwareDaoImpl) DeleteRemoteAppWithSoftwareID(ctx context.Context, softwareID snowflake.ID) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.Where("software_id = ?", softwareID).Delete(&model.RemoteApp{})
	if err != nil {
		return err
	}
	return nil
}

func (d *softwareDaoImpl) GetRemoteApp(ctx context.Context, remoteAppID snowflake.ID) (*model.RemoteApp, bool, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	remoteApp := &model.RemoteApp{}
	exist, err := session.Where("id = ?", remoteAppID).Get(remoteApp)
	if err != nil {
		return nil, exist, err
	}
	return remoteApp, exist, nil
}

func (d *softwareDaoImpl) GetRemoteApps(ctx context.Context, softwareID snowflake.ID) ([]*model.RemoteApp, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	var remoteApps []*model.RemoteApp
	err := session.Where("software_id = ?", softwareID).Find(&remoteApps)
	if err != nil {
		return nil, err
	}
	return remoteApps, nil
}

func (d *softwareDaoImpl) InsertSoftwarePresets(ctx context.Context, presets []*model.SoftwarePreset) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.Insert(presets)
	if err != nil {
		return err
	}
	return nil
}

func (d *softwareDaoImpl) DeleteSoftwarePresets(ctx context.Context, softwareID snowflake.ID) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.Where("software_id = ?", softwareID).Delete(&model.SoftwarePreset{})
	if err != nil {
		return err
	}
	return nil
}

func (d *softwareDaoImpl) GetSoftwarePresets(ctx context.Context, softwareID snowflake.ID) ([]*model.SoftwarePreset, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var presets []*model.SoftwarePreset
	if softwareID != 0 {
		session.Where("software_id = ?", softwareID)
	}

	err := session.Find(&presets)
	if err != nil {
		return nil, err
	}

	return presets, nil
}
