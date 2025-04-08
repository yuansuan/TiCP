import { observable } from 'mobx'

class Paginations {
  @observable page_index: number = 1
  @observable page_size: number = 20

  constructor(props: any) {
    this.page_index = props?.index || 1
    this.page_size = props?.size || 20
  }
}

export class OutcomingPageAware extends Paginations {}

export class IncomingPageAware {
  @observable total: number = 0

  constructor(props: any) {
    this.total = props?.total
  }
}
