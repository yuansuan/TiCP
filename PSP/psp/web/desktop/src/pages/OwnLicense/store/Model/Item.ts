/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'

export class BaseItem {
  @observable app_id: string
  @observable app_name: string
  @observable version: string
  @observable active: boolean
  @observable merchandise_id: string
}

export class Item extends BaseItem {
  constructor(props?: Partial<BaseItem>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update(props: Partial<BaseItem>) {
    Object.assign(this, props)
  }
}
