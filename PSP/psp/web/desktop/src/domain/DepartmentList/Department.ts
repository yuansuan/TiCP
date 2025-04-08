/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { DepartmentUser } from './DepartmentUser'

export class BaseDepartment {
  @observable id: string
  // 企业ID
  @observable company_id: string
  // 部门状态
  @observable status: number
  // 部门名
  @observable name: string
  @observable remark: string
  @observable users: DepartmentUser[] = []
}

export class Department extends BaseDepartment {
  constructor(props?: Partial<BaseDepartment>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update(data: Partial<BaseDepartment>) {
    Object.assign(this, data)
  }
}
