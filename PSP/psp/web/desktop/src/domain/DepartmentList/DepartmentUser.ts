/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'

export class BaseDepartmentUser {
  @observable user_id: string
  // 姓名
  @observable real_name: string
  // 电话
  @observable phone: string
  // email
  @observable email: string

  // 用户名
  @observable user_name: string
  // 显示用户名
  @observable display_user_name: string
}

export class DepartmentUser extends BaseDepartmentUser {
  constructor(props?: Partial<BaseDepartmentUser>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update(data: Partial<BaseDepartmentUser>) {
    Object.assign(this, data)
  }
}
