package zone

import (
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/zonelist"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
)

// List 获取所有区域
// 一个获取所有分区信息的接口(返回包括计算区域和存储区域以及3D云应用区域)
func List(cfg config.CustomT) zonelist.Data {
	data := zonelist.Data{}
	data.Zones = make(map[string]*schema.Zone)

	// 从配置文件中读取区域信息
	data.Zones = cfg.Zones

	return data
}
