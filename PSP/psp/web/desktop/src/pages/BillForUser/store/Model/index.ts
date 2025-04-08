import { observable, action } from 'mobx'
import { BillUser, IRequest as IBillUserRequest } from './billUser'

export class BaseList {
  @observable list: BillUser[] = []
  @observable page_ctx: {
    index: number
    size: number
    total: number
  }
  @observable total_amount: number
  @observable total_refund_amount: number
}

type IRequest = Omit<BaseList, 'list'> & {
  list: IBillUserRequest[]
}

export class BillUserList extends BaseList {
  @action
  update({ list, ...props }: Partial<IRequest>) {
    Object.assign(this, props)

    if (list) {
      this.list = list.map(item => new BillUser(item))
    }
  }
}
