/* Copyright (C) 2016-present, Yuansuan.cn */

import { observable, action } from 'mobx'
import { Timestamp } from '@/domain/common'

export class BaseDetail {
  @observable id: string
  @observable account_id: string
  @observable bill_sign: number
  @observable amount: number
  @observable trade_type: number
  @observable trade_id: string
  @observable trade_time: Timestamp
  @observable account_balance_contain_freezed: number
  @observable remark: string
  @observable out_trade_id: string
}

export type DetailRequest = Omit<BaseDetail, 'trade_time'> & {
  trade_time: {
    seconds: number
    nanos: number
  }
}

export class Detail extends BaseDetail {
  constructor(props?: Partial<DetailRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = ({ trade_time, ...props }: Partial<DetailRequest>) => {
    Object.assign(this, props)

    if (trade_time) {
      this.trade_time = new Timestamp(trade_time)
    }
  }
}
