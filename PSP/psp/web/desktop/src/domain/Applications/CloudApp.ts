/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, computed } from 'mobx'
import { Http } from '@/utils'
import App from './App'

export default class CloudAppList {
  @observable list: App[] = []

  async fetch() {
    const res = await Http.get('/app/list')
    this.list = res.data?.apps?.map(item => new App(item))
    // hack enough
    localStorage.setItem(
      'FLAG_ENTERTAINMENT',
      JSON.stringify(this.publishedAppList)
    )
  }
  @computed
  get publishedAppList() {
    return this.list && this.list?.filter(item => item.state === 'published')
  }
}
