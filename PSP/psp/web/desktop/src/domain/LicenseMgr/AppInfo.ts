import { observable, action, computed } from 'mobx'
import { LicenseInfo } from './LicenseInfo'
import moment from 'moment'
interface IAppInfo {
  id: string
  create_time: string
  os: number
  desc: string
  compute_rule: string
  license_type: string
  app_type: string
  license_infos: LicenseInfo
}

/*
enum OS {
  UNKNOWNOS = 0;
  LINUX = 1;
  WIN = 2;
}

// 发布状态
enum Status {
  UNKNOWNSTATUS = 0;
  PUBLISHED = 1;
  UNPUBLISHED = 2;
}
*/

export const OS_NAME = ['UNKNOWN', 'Linux', 'Windows']
export const STATUS_NAME = ['UNKNOWN', '已发布', '未发布']
export const TYPE_NAME = {
  own: '自有'
}

export class AppInfo implements IAppInfo {
  @observable id
  @observable license_type: string
  @observable os
  @observable create_time
  @observable compute_rule
  @observable desc
  @observable app_type
  @observable license_infos: LicenseInfo


  constructor(obj: IAppInfo) {
    this.init(obj)
  }

  @computed
  get os_name() {
    return OS_NAME[this.os] || '--'
  }

  // @computed
  // get status_name() {
  //   return STATUS_NAME[this.status] || '--'
  // }

  @computed
  get type_name() {
    return TYPE_NAME[this.license_type] || '--'
  }
  
  @action
  init(obj) {
    if (!obj) return
    Object.assign(this, {
      ...obj,
      create_time:moment(obj.create_time).format('YYYY-MM-DD HH:mm:ss'),
      license_infos: obj.license_infos?.map(item => {
        item.license_type = obj.app_type
        return  new LicenseInfo(item)
      }) || []
    })
  }
}
