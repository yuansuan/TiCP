/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import { Record } from './record'
import { PageCtx } from '../common'
import { userServer } from '@/server'
import {Http} from '@/utils'
export class RecordList {
  @observable list: Record[] = []
  @observable pageCtx: PageCtx = new PageCtx()
  @observable unhandledCount = 0
  @action
  updateList = list => {
    this.list = [...list]
  }
  @action
  updateCount = count => {
    this.unhandledCount = count
  }

  fetchUnhandledCount = () =>
    userServer
      .getShareUnreadCount()
      .then(res => {
        this.updateCount(res.data)
        return res
      })

  fetch = async (params: {
    index: number
    size: number
  }) => {
    const { data } = await userServer.getShareList(params)

    runInAction(() => {
      this.updateList((data.list|| []).map(item => new Record(item)))
      this.pageCtx.update(data.page)
    })
  }

  fetchLast = () => {
    this.fetchUnhandledCount()

    return this.fetch({
      index: 1,
      size: 5,
    })
  }

  readAll = async () => {
    await Http.put('/storage/share/readAll')
    runInAction(() => {
      this.list.forEach(item => {
        item.state = 2
      })
    })

     // 设置未读
     this.updateCount(0)
  }

}
