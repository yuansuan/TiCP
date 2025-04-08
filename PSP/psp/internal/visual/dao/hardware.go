package dao

import (
	"context"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/visual/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type hardwareDaoImpl struct{}

func NewHardwareDao() HardwareDao {
	return &hardwareDaoImpl{}
}

func (d *hardwareDaoImpl) ListHardware(ctx context.Context, IDs []snowflake.ID, name string, cpu, mem, gpu, offset, limit int) ([]*model.Hardware, int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var hardwares []*model.Hardware
	if len(IDs) > 0 {
		session.In("id", IDs)
	}
	if name != "" {
		session.Where("name like ?", "%"+name+"%")
	}
	if cpu > 0 {
		session.Where("cpu = ?", cpu)
	}
	if mem > 0 {
		session.Where("mem = ?", mem)
	}
	if gpu > 0 {
		session.Where("gpu = ?", gpu)
	}
	if offset > 0 {
		session.Limit(limit, (offset-1)*limit)
	} else {
		session.Limit(limit)
	}

	// 适配查询 gpu 个数等于 0
	if gpu == consts.NumberEqualZeroMark {
		session.Where("gpu = ?", 0)
	}
	// 适配查询 gpu 个数大于 0
	if gpu == consts.NumberGreaterThanZeroMark {
		session.Where("gpu > ?", 0)
	}

	total, err := session.Select("*").Desc("update_time").FindAndCount(&hardwares)
	if err != nil {
		return nil, 0, err
	}
	return hardwares, total, nil
}

func (d *hardwareDaoImpl) InsertHardware(ctx context.Context, hardware *model.Hardware) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.Insert(hardware)
	if err != nil {
		return err
	}
	return nil
}

func (d *hardwareDaoImpl) UpdateHardware(ctx context.Context, hardware *model.Hardware) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.ID(hardware.ID).Omit("id, out_hardware_id").AllCols().Update(hardware)
	if err != nil {
		return err
	}
	return nil
}

func (d *hardwareDaoImpl) DeleteHardware(ctx context.Context, hardwareID snowflake.ID) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.ID(hardwareID).Delete(&model.Hardware{})
	if err != nil {
		return err
	}
	return nil
}

func (d *hardwareDaoImpl) GetHardware(ctx context.Context, hardwareID snowflake.ID, name string) (*model.Hardware, bool, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	hardware := &model.Hardware{}
	if int64(hardwareID) > 0 {
		session.Where("id = ?", hardwareID)
	}
	if name != "" {
		session.Where("name = ?", name)
	}

	exist, err := session.Get(hardware)
	if err != nil {
		return nil, exist, err
	}
	return hardware, exist, nil
}
