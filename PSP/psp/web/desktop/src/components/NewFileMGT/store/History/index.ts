/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, computed } from 'mobx'
import { Record, BaseRecord } from './Record'

export class BaseHistory {
  @observable list: Record[] = []
  @observable cursor = -1
}

export class History extends BaseHistory {
  constructor(props?: Partial<BaseHistory>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = (props: Partial<BaseHistory>) => {
    Object.assign(this, props)
  }

  @computed
  get current() {
    return this.list.length > 0 ? this.list[this.cursor] : undefined
  }

  @computed
  get prevable() {
    return this.cursor > 0
  }

  @computed
  get nextable() {
    return this.cursor < this.list.length - 1
  }

  @action
  delete(index) {
    this.list.splice(index, 1)
  }

  @action
  push(item: Partial<BaseRecord>) {
    if (this.nextable) {
      this.list = this.list.slice(0, this.cursor + 1)
    }

    this.list.push(new Record(item))
    this.cursor += 1
  }

  @action
  prev = () => {
    if (!this.prevable) {
      return undefined
    }

    this.cursor -= 1
    return this.list[this.cursor]
  }

  @action
  next = () => {
    if (!this.nextable) {
      return undefined
    }

    this.cursor += 1
    return this.list[this.cursor]
  }
}
