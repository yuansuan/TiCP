import { observable, action, computed } from 'mobx'

// path history manager
export default class History {
  @observable list = []
  @observable cursor = -1

  constructor(initialValue?: any) {
    initialValue !== undefined && this.push(initialValue)
  }

  @computed
  get current() {
    return this.list.length > 0 ? this.list[this.cursor] : undefined
  }

  @computed
  get prevDisabled() {
    return this.cursor <= 0
  }

  @computed
  get nextDisabled() {
    return this.cursor >= this.list.length - 1
  }

  @action
  delete(index) {
    if (this.cursor === index) {
      this.cursor = this.list.length
    } else if (this.cursor > index) {
      this.cursor -= 1
    }
    this.list.splice(index, 1)
  }

  @action
  push(item) {
    if (!this.nextDisabled) {
      this.list = this.list.slice(0, this.cursor + 1)
    }
    this.list.push(item)
    this.cursor += 1
  }

  @action
  prev = () => {
    if (this.prevDisabled) {
      return undefined
    }

    this.cursor -= 1
    return this.list[this.cursor]
  }

  @action
  next = () => {
    if (this.nextDisabled) {
      return undefined
    }

    this.cursor += 1
    return this.list[this.cursor]
  }

  // you can us for...of and ... to iterate the history
  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
