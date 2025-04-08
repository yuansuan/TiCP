/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'

export class BaseRecord {
  @observable path: string
}

export class Record extends BaseRecord {
  constructor(props?: Partial<BaseRecord>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = (props: Partial<BaseRecord>) => {
    Object.assign(this, props)
  }
}
