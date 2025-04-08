/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import { Department } from './Department'
import { departmentServer } from '@/server'
import { env } from '@/domain'

export class BaseList {
  @observable page_index: number = 1
  @observable page_size: number = 10
  @observable name: string = null
  @observable status: number = 1
  @observable list: Department[] = []
  @observable total: number = 10
}

export class DepartmentList extends BaseList {
  @action
  setName = name => {
    this.name = name
  }
  @action
  setStatus = status => {
    this.status = status
  }
  @action
  setPage = (current, pageSize) => {
    this.page_index = current
    this.page_size = pageSize
  }

  @action
  update = ({ list, total }) => {
    if (list) {
      this.list = [...list].map(item => new Department(item))
      this.total = total
    }
  }

  fetch = async () => {
    const { data } = await departmentServer.getList({
      company_id: env.company.id,
      page_index: this.page_index,
      page_size: this.page_size,
      name: this.name,
      status: this.status
    })

    runInAction(() => {
      this.update({
        list: data.list,
        total: data.page_ctx.total
      })
    })
  }
}
