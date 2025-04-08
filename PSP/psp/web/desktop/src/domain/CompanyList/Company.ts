/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, computed } from 'mobx'

export class BaseCompany {
  @observable id: string
  @observable name: string
  @observable account_id: string
  @observable roles: { id: string; name: string }[] = []
  @observable cloud_type: string
  @observable live_chat_id: string
  @observable box
  @observable is_open_department_manage: number
  @observable new_job_mgr: string
  @observable max_projects: number
  @observable new_file_mgr: string
  @observable zone_list: string
  @observable job_zone_list: string
}

export class Company extends BaseCompany {
  constructor(props?: Partial<BaseCompany>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @computed
  get isOpenDepMgr() {
    return this.is_open_department_manage === 1
  }

  // TODO: 新老板文件管理标识
  @computed
  get isOpenNewFileMgr() {
    return this.new_file_mgr === 'enable'
  }

  @computed
  get isOpenStandardJobMgr() {
    return this.new_job_mgr === 'enable'
  }

  @computed
  get isOpenNewJobMgr() {
    return this.new_job_mgr === 'enableNew'
  }

  @computed
  get isMixed() {
    return this.cloud_type === 'mixed'
  }

  @action
  update(data: Partial<BaseCompany>) {
    Object.assign(this, data)
  }
}
