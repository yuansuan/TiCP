/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import { companyList } from '@/domain'
import { visualServer } from '@/server'

export class BaseVisualConfig {
  @observable isOpen: boolean
  @observable activeTerminal: number
  @observable maxTerminal: number
}

export default class VisualConfig extends BaseVisualConfig {
  constructor(props?) {
    super()
    !!props && Object.assign(this, props)
  }

  @action
  update(props?) {
    Object.assign(this, props)
  }

  @action
  async fetch() {
    const { data } = await visualServer.fetch()
    runInAction(() => {
      this.update(data)
    })
  }

  @action
  async use() {
    visualServer.bindCurrentUser()
  }

  get showVisualizeApp() {
    /**
     * 在导航栏中显示3D云应用的条件：
     *   1。企业用户
     *   2。在OMS中开启了可视化应用
     */
    return !!companyList.current && !!this.isOpen
  }
}
