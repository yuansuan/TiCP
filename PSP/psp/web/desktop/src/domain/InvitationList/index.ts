/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import { Invitation } from './Invitation'
import { PageCtx } from '../common'
import { companyServer, userServer } from '@/server'

export class InvitationList {
  @observable list: Invitation[] = []
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
      .getInviteList({
        page_index: 1,
        page_size: 1,
        status: 1,
      })
      .then(res => {
        let { page_ctx } = (res as any).data
        page_ctx = page_ctx || { total: 0 }
        this.updateCount(page_ctx.total)
        return res
      })

  fetch = async (params: {
    status?: number
    page_index: number
    page_size: number
  }) => {
    const { data } = await userServer.getInviteList(params)

    runInAction(() => {
      this.updateList(data.list.map(item => new Invitation(item)))
      this.pageCtx.update(data.page_ctx)
    })
  }

  fetchLast = () => {
    this.fetchUnhandledCount()

    return this.fetch({
      page_index: 1,
      page_size: 5,
      status: 1,
    })
  }

  fetchByCompany = async (params: {
    company_id: string
    status?: number
    page_index: number
    page_size: number
  }) => {
    const { data } = await companyServer.getInviteList(params)

    runInAction(() => {
      this.updateList(data.list.map(item => new Invitation(item)))
      this.pageCtx.update(data.page_ctx)
    })
  }
}
