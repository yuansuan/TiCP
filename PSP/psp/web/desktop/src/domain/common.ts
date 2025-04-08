/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { formatUnixTime } from '@/utils'

class BaseTimestamp {
  @observable seconds = 0
  @observable nanos = 0
}

export class Timestamp extends BaseTimestamp {
  constructor(props?: Partial<BaseTimestamp>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = (props: Partial<BaseTimestamp>) => {
    Object.assign(this, props)
  }

  toString() {
    return this.seconds ? formatUnixTime(this.seconds) : '--'
  }
}

class BasePageCtx {
  @observable index = 1
  @observable size = 10
  @observable total = 0
}

export class PageCtx extends BasePageCtx {
  constructor(props?: Partial<BasePageCtx>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = (props: Partial<BasePageCtx>) => {
    Object.assign(this, props)
  }
}
export const routeMapConfig = {
  visualization: ['visualmgr', 'cloudApps']
}

export const tabMapConfig = {
  visualization: ['remote']
}

//核二院 监控与报表
export const PRECISION = 2
export const VISUAL_PRECISION = 6
