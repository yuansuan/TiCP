/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { Perm } from './Perm'
import { Http } from '@/utils'

/**
 * 当前用户-当前企业下的权限列表
 */
export class PermList {
  @observable list: Perm[] = []

  @action
  update = ({ list, ...props }) => {
    Object.assign(this, props)

    if (list) {
      this.list = [...list].map(item => new Perm(item))
    }
  }

  fetch = async () => {
    const { data } = await Http.get('/company/user/permission')
    this.update({
      list: data.list,
    })
  }

  check = perm => this.list.findIndex(item => item.code === perm) > -1
}
