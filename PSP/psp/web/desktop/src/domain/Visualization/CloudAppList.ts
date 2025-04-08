/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import CloudApp from './CloudApp'
import { Http } from '@/utils'
export default class CloudAppList {
  @observable list: Array<CloudApp> = []
  @action
  async fetch() {
    const res = await Http.get('/visual/app/list')
    const l = res.data.map((item: any) => {
      return new CloudApp(item)
    })
    runInAction(() => {
      this.list = l
    })
  }
}
