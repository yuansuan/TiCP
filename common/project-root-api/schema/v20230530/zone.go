package v20230530

import "strings"

// Zone 分区
type Zone struct {
	// 超算域名
	HPCEndpoint string `json:"HPCEndpoint,omitempty" yaml:"hpc_endpoint"`
	// 云存储域名
	StorageEndpoint string `json:"StorageEndpoint,omitempty" yaml:"storage_endpoint"`
	// 启用3D云应用
	CloudAppEnable bool `json:"CloudAppEnable,omitempty" yaml:"cloud_app_enable"`
	// SyncRunner 域名
	SyncRunnerEndpoint string `json:"SyncRunnerEndpoint,omitempty" yaml:"sync_runner_endpoint"`
}

// Zones 分区列表
type Zones map[string]*Zone

// GetZoneByEndpoint 获取域名对应的分区
func (zs Zones) GetZoneByEndpoint(endpoint string) string {
	for name, z := range zs {
		if z.HPCEndpoint != "" && strings.HasPrefix(endpoint, z.HPCEndpoint) {
			return name
		}
		if z.StorageEndpoint != "" && strings.HasPrefix(endpoint, z.StorageEndpoint) {
			return name
		}
	}
	return "unknown"
}

// Exist 判断域名分区是否存在
func (zs Zones) Exist(endpoint string) bool {
	return zs.GetZoneByEndpoint(endpoint) != "unknown"
}

// List 获取分区列表
func (zs Zones) List() []string {
	var list []string
	for name := range zs {
		list = append(list, name)
	}
	return list
}

// IsZone 是否是列表中的分区
func (zs Zones) IsZone(zone string) bool {
	if _, ok := zs[zone]; ok {
		return true
	}
	return false
}
