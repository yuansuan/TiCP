import { observable, action, computed } from 'mobx'
import moment from 'moment'

type ModuleConf = {
  module_name: string
  id: string
  free_num: number
  total: number // 实时总数量（监控统计的）
}

// License类型，跟供应商有关
export const LicenseProviderType = [
  // 自有
  'SELFOWNED',
  // 外部
  'OTHEROWNED',
  // 寄售
  'CONSIGNED'
]

export const CollectorTypeLabel = ['flex', 'lsdyna', 'altair']

type LicenseProviderType = 'SELFOWNED' | 'OTHEROWNED' | 'CONSIGNED'

interface ILicenseInfo {

  id: string
  provider: string
  manager_id: string
  // 许可证变量
  license_env_var: string
  // mac地址
  mac_addr: string
  // 路径
  tool_path: string
  // 许可证
  license_url: string
  // 端口
  port: number
  // license许可证序列号
  license_num: string
  // 模块配置列表
  module_config_infos: ModuleConf[]
  // 调度优先级
  weight: number
  // 使用有效期 开始
  begin_time: string
  // 使用有效期 结束
  end_time: string
  // 是否授权
  auth: boolean
 
  // 供应商类型
  license_type: string
  
}

export class LicenseInfo implements ILicenseInfo {
  @observable id: string
  @observable manager_id: string

  @observable provider: string
  // 许可证变量
  @observable license_env_var: string
  // mac地址
  @observable mac_addr: string
  // 路径
  @observable tool_path: string
  // 许可证
  @observable license_url: string
  // 端口
  @observable port: number
  // license许可证序列号
  @observable license_num: string
  // 模块配置列表
  @observable module_config_infos: ModuleConf[] = []
  // 计算规则
  @observable compute_rule: number
  @observable weight: number
  
  // 供应商类型
  @observable license_type: string
  @observable auth: boolean
  @observable begin_time: string
  @observable end_time: string
 

  constructor(obj: ILicenseInfo) {
    this.init(obj)
  }

  @computed
  get license_type_string(): string {
    return CollectorTypeLabel[this.license_type]
  }

  @action
  init(obj) {
    if (!obj) return
    Object.assign(this, {
      ...obj,
      begin_time: moment(obj.begin_time).format('YYYY-MM-DD HH:mm:ss'),
      end_time: moment(obj.end_time).format('YYYY-MM-DD HH:mm:ss'),
      module_config_infos: obj.module_config_infos
    })
  }
}
