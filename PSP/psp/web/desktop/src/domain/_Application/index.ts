/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, computed } from 'mobx'
import { Http } from '@/utils'
import App from './App'
import { currentUser } from '@/domain'
export default class ApplicationList {
  @observable list: App[] = []
  
  async fetch(isDesktop: boolean = false) {
    const res = await Http.get('/app/list', {
      params: {
        has_permission: true,
        state: 'published',
        desktop:isDesktop
      }
    })

    this.list = res.data?.apps?.map(item => new App(item))

    localStorage.setItem(
      'FLAG_ENTERTAINMENT',
      JSON.stringify(this.publishedAppList)
    )
  }
  @computed
  get publishedAppList() {
    return this.list
  }
}
