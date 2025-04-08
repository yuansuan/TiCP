/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Http } from '@/utils'
import { action, computed, observable, runInAction } from 'mobx'

type WebConfigData = {
  statistics?: boolean
  livechat_id?: string
}

export class WebConfig {
  @observable data: WebConfigData = {}

  @computed
  get statisticsEnabled() {
    return !!this.data.statistics
  }

  @computed
  get liveChatId() {
    return this.data.livechat_id
  }

  @action
  async init() {
    const { data } = await Http.get('/web_config')
    runInAction(() => {
      Object.assign(this.data, data)
    })
  }
}
