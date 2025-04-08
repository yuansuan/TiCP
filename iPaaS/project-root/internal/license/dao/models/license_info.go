package models

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	license "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info/type"
	"sort"
	"time"
)

type LicenseInfo struct {
	Id             snowflake.ID `json:"id" db:"id" xorm:"pk BIGINT(20)"`
	ManagerId      snowflake.ID `json:"manager_id" db:"manager_id" xorm:"not null comment('manager id') BIGINT(20)"`
	Provider       string       `json:"provider" db:"provider" xorm:"not null comment('提供者') VARCHAR(255)"`
	LicenseServer  string       `json:"license_server" db:"license_server" xorm:"not null comment('许可证变量') VARCHAR(255)"`
	MacAddr        string       `json:"mac_addr" db:"mac_addr" xorm:"not null comment('Mac地址') VARCHAR(255)"`
	LicenseUrl     string       `json:"license_url" db:"license_url" xorm:"not null comment('许可证服务器') VARCHAR(255)"`
	LicensePort    int          `json:"license_port" db:"license_port" xorm:"not null default 0 comment('端口') INT(11)"`
	LicenseProxies string       `json:"license_proxies" db:"license_proxies" xorm:"not null default '' comment('HpcEndpoint对应的许可证服务器地址') VARCHAR(255)"`
	LicenseNum     string       `json:"license_num" db:"license_num" xorm:"not null comment('licenses许可证序列号') VARCHAR(255)"`
	Weight         int          `json:"weight" db:"weight" xorm:"not null default 0 comment('调度优先级') INT(11)"`
	BeginTime      time.Time    `json:"begin_time" db:"begin_time" xorm:"not null comment('使用有效期 开始') DATETIME"`
	EndTime        time.Time    `json:"end_time" db:"end_time" xorm:"not null comment('使用有效期 结束') DATETIME"`
	Auth           int          `json:"auth" db:"auth" xorm:"not null default 2 comment('是否被授权 1-授权 2-未授权') TINYINT(1)"`
	LicenseType    license.Type `json:"license_type" db:"license_type" xorm:"not null default 0 comment('供应商类型: 1-自有，2-外部，3-寄售') TINYINT(1)"`

	// 当license type为外部时，下面的参数是必须的
	ToolPath string `json:"tool_path" db:"tool_path" xorm:"not null comment('lmutil所在路径') VARCHAR(255)"`
	//  flex: "ABAQUS,Adams,Ansys Maxwell,CFX,Fluent,ANSYS HFSS,Icepak,Mechanical,SIwave,Workbench,Ansys,MAC Easy5,MSC Adams,MSC scFlowso,MSC Dytran,MSC Easy5,MSC Marc,MSC Nastran,STAR-CCM+"
	//  lsdyna: "LS-DYNA"
	//  altair: "OptiStruct,optistruct,FEKO"
	CollectorType         string    `json:"collector_type" db:"collector_type" xorm:"not null default '' comment('收集器类型: flex, lsdyna, altair, dsli') VARCHAR(255)"`
	HpcEndpoint           string    `json:"hpc_endpoint" db:"hpc_endpoint" xorm:"not null default '' comment('超算endpoint') VARCHAR(255)"`
	AllowableHpcEndpoints []string  `json:"allowable_hpc_endpoints" db:"allowable_hpc_endpoints" xorm:"not null default '' comment('支持的HpcEndpoint范围') VARCHAR(255)"`
	LicenseServerStatus   string    `json:"license_server_status" db:"license_server_status" xorm:"not null default 'abnormal' comment('license服务状态: normal-正常，abnormal-异常') VARCHAR(64)"`
	CreateTime            time.Time `json:"create_time" xorm:"created"`
	UpdateTime            time.Time `json:"update_time" xorm:"updated"`
}

const (
	LicenseServerStatusNormal   string = "normal"   // 正常
	LicenseServerStatusAbnormal string = "abnormal" // 异常
)

type LicenseProxy struct {
	// 许可证服务器地址
	Url string `json:"Url"`
	// 端口
	Port int `json:"Port"`
}

func (*LicenseInfo) TableName() string {
	return "license_info"
}

func (l *LicenseInfo) IsSelfLic() bool {
	return l.LicenseType.IsSelfLic()
}

func (l *LicenseInfo) IsConsignedLic() bool {
	return l.LicenseType.IsConsignedLic()
}

func (l *LicenseInfo) IsOthersLic() bool {
	return l.LicenseType.IsOthersLic()
}

func (l *LicenseInfo) GetLicenseProxies(hpcEndpoint string) (string, error) {
	if l.LicenseProxies == "" {
		if l.LicensePort > 0 {
			return fmt.Sprintf("%s=%d@%s", l.LicenseServer, l.LicensePort, l.LicenseUrl), nil
		} else {
			return fmt.Sprintf("%s=%s", l.LicenseServer, l.LicenseUrl), nil
		}
	}

	var licenseProxies map[string]LicenseProxy
	err := json.Unmarshal([]byte(l.LicenseProxies), &licenseProxies)
	if err != nil {
		logging.Default().Errorf("unmarshal license proxies fail, license proxies: %s, err: %v", l.LicenseProxies, err)
		return "", err
	}

	if licenseAddress, ok := licenseProxies[hpcEndpoint]; ok {
		if licenseAddress.Port > 0 {
			return fmt.Sprintf("%s=%d@%s", l.LicenseServer, licenseAddress.Port, licenseAddress.Url), nil
		} else {
			return fmt.Sprintf("%s=%s", l.LicenseServer, licenseAddress.Url), nil
		}
	}

	return "", fmt.Errorf("license proxies not found, license proxies: %s, hpc endpoint: %s", l.LicenseProxies, l.HpcEndpoint)
}

type JoinEntity struct {
	LicenseManager `xorm:"extends"`
	LicenseInfo    `xorm:"extends"`
	ModuleConfig   `xorm:"extends"`
}

func ToLicenseManagerExt(entities []*JoinEntity) []*LicenseManagerExt {
	saved := map[snowflake.ID][]*JoinEntity{}
	for _, en := range entities {
		if _, ok := saved[en.LicenseManager.Id]; !ok {
			saved[en.LicenseManager.Id] = []*JoinEntity{}
		}
		saved[en.LicenseManager.Id] = append(saved[en.LicenseManager.Id], en)
	}
	var res []*LicenseManagerExt
	for _, l := range saved {
		res = append(res, ToOneLicenseManagerExt(l))
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].LicenseManager.Id < res[j].LicenseManager.Id
	})
	return res
}

func ToOneLicenseManagerExt(entities []*JoinEntity) *LicenseManagerExt {
	if len(entities) == 0 {
		return nil
	}
	res := &LicenseManagerExt{}
	res.LicenseManager = entities[0].LicenseManager

	maps := lo.GroupBy(entities, func(entity *JoinEntity) string {
		return entity.LicenseInfo.Id.String()
	})
	for _, v := range maps {
		licInfo := &LicenseInfoExt{
			LicenseInfo: v[0].LicenseInfo,
		}
		if licInfo.Id == 0 {
			continue
		}
		for _, module := range v {
			m := module.ModuleConfig
			if m.Id == 0 {
				continue
			}
			licInfo.Modules = append(licInfo.Modules, &m)
		}
		sort.Slice(licInfo.Modules, func(i, j int) bool {
			return licInfo.Modules[i].Id < licInfo.Modules[j].Id
		})
		res.Licenses = append(res.Licenses, licInfo)
	}
	sort.Slice(res.Licenses, func(i, j int) bool {
		return res.Licenses[i].Id < res.Licenses[j].Id
	})
	return res
}

type LicenseManagerExt struct {
	LicenseManager
	Licenses []*LicenseInfoExt
}

type LicenseInfoExt struct {
	LicenseInfo
	Modules []*ModuleConfig
}

type AddOrUpdateEntity struct {
	Manager *LicenseManager
	Infos   []*LicenseInfo
	Modules []*ModuleConfig
}

func (JoinEntity) TableName() string {
	return "license_manager"
}

func (lExt *LicenseInfoExt) GetRemainingLicenseNum(moduleName string) (int, error) {
	for _, module := range lExt.Modules {
		if moduleName == module.ModuleName {
			return module.Used, nil
		}
	}

	return -1, errors.New(fmt.Sprintf("module name not found, moduleName: %s", moduleName))
}
