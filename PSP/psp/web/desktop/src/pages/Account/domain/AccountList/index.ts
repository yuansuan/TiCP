/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import { Detail, DetailRequest } from './Detail'
import { accountServer } from '@/server'

export class BaseDetailList {
  @observable list: Detail[] = []
  @observable page_ctx: { index: number; size: number; total: number } = {
    index: 1,
    size: 10,
    total: 0,
  }
}

type Request = Omit<BaseDetailList, 'list'> & {
  list: DetailRequest[]
}

export class DetailList extends BaseDetailList {
  constructor(props?: Partial<Request>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update({ list, ...props }: Partial<Request>) {
    Object.assign(this, props)

    if (list) {
      this.list = list.map(item => new Detail(item))
    }
  }

  fetch = async (params: Parameters<typeof accountServer.getDetailList>[0]) => {
    const { data } = await accountServer.getDetailList(params)

    runInAction(() => {
      this.update({
        list: data.list,
        page_ctx: data.page_ctx,
      })
    })
  }
}
