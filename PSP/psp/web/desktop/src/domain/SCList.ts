/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { action, observable } from 'mobx'

type SCItem = {
  sc_id: string
  sc_name: string
  tier_name: string
}

class BaseSCList {
  @observable list: SCItem[] = []
}

export class SCList extends BaseSCList {
  constructor(props?: Partial<BaseSCList>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update(props?: Partial<BaseSCList>) {
    Object.assign(this, props)
  }

  getName(scId: string) {
    const sc = this.list.find(item => item.sc_id === scId)
    return sc ? sc.tier_name : ''
  }
}
