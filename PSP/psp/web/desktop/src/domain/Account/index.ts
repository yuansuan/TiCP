/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, runInAction } from 'mobx'
import { Timestamp } from '@/domain/common'
import { accountServer } from '@/server'
import { env } from '@/domain'
import { Http } from '@/utils'

export class BaseAccount {
  @observable account_id: string
  @observable account_name: string
  @observable currency: string
  @observable account_balance: number = 0 // 账户总余额
  @observable freezed_amount: number = 0 // 冻结金额
  @observable normal_balance: number = 0 // 账户充值
  @observable award_balance: number = 0 // 账户赠送
  @observable credit_quota_amount: number = 0 //授信额度
  @observable frozen_status: number // 0 :正常 1：冻结
  @observable preferential_amount: number = 0 //优惠金额
  @observable create_time: Timestamp
  @observable update_time: Timestamp
}

type AccountRequest = Omit<BaseAccount, 'create_time' | 'update_time'> & {
  create_time: {
    seconds: number
    nanos: number
  }
  update_time: {
    seconds: number
    nanos: number
  }
}

export class Account extends BaseAccount {
  constructor(props?: Partial<AccountRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update = (props: Partial<AccountRequest>) => {
    Object.assign(this, props)
  }

  fetch = async (account_id?: string) => {
    let { data } = await accountServer.get()

    runInAction(() => {
      this.update(data)
    })
  }
}
