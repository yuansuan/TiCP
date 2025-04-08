/* Copyright (C) 2016-present, Yuansuan.cn */

import { observable, action } from 'mobx'

export class BaseBillUser {
  @observable bill_id: string
  @observable job_id: string
  @observable billing_month: string
  @observable user_info: string
  @observable user_id: string
  @observable update_time: number
  @observable merchandise_name: string
  @observable merchandise_price_unit: number
  @observable merchandise_quantity: number
  @observable real_amount: number
  @observable refund_amount: number
  @observable out_resource_type: number
  @observable out_biz_id: string
  //计算作业的作业编号
  @observable bill_job_id: string
}

export type IRequest = BaseBillUser

export class BillUser extends BaseBillUser {
  constructor(props?: Partial<IRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update(props: Partial<IRequest>) {
    Object.assign(this, props)
  }
}
