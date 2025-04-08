import { action, observable, runInAction } from 'mobx'
import { formatTime } from '@/utils/formatter'
import { Http } from '@/utils'
import BillList from './Bill/BillList'

export enum AccountStatus {
  DELETED = 'DELETED',
  NORMAL = 'NORMAL'
}

interface IUpdateAccountBody {
  real_customer_id?: string | null
  name?: string | null
  withdraw_enabled?: boolean | null
  credit_quota?: number | null
}

interface IAccount {
  // 账户ID
  id: string
  // 客户ID （个人为user_id, 企业为 company_id）
  customer_id: string
  // 实名认证ID （预留）
  real_customer_id: string
  name: string
  // 币种 （CNY为人民币 ISO4217）
  currency: string
  // 账户余额（不含冻结，即未结算）
  account_balance: number
  // 冻结金额
  freezed_amount: number
  // 普通余额
  normal_balance: number
  // 赠送余额
  award_balance: number
  // 是否提现
  withdraw_enabled: boolean
  // 授信额度
  credit_quota: number
  status: AccountStatus
  // 账户余额（含冻结，即未结算）
  account_balance_contain_freezed: number

  create_time: any
  update_time: any
}

export class Account implements IAccount {
  // 账户ID
  @observable
  id: string
  // 客户ID （个人为user_id, 企业为 company_id）
  @observable customer_id: string
  // 实名认证ID （预留）
  @observable real_customer_id: string
  @observable name: string
  // 币种 （CNY为人民币 ISO4217）
  @observable currency: string
  // 账户余额（不含冻结，即未结算）
  @observable account_balance: number = 0
  // 冻结金额
  @observable freezed_amount: number = 0
  // 普通余额
  @observable normal_balance: number = 0
  // 赠送余额
  @observable award_balance: number = 0
  // 是否提现
  @observable withdraw_enabled: boolean
  // 授信额度
  @observable credit_quota: number = 0
  @observable status: AccountStatus
  // 账户余额（含冻结，即未结算）
  @observable account_balance_contain_freezed: number = 0

  @observable create_time: any
  @observable update_time: any

  @observable billList = new BillList()

  get create_time_string() {
    return this.create_time !== null
      ? formatTime(this.create_time.seconds)
      : null
  }

  get update_time_string() {
    return this.update_time !== null
      ? formatTime(this.update_time.seconds)
      : null
  }

  @action
  public init = request => {
    request && Object.assign(this, { ...request })
  }

  fetch = (account_id, company_id) => {
    let url = `/company/account?accountId=${account_id}&companyId=${company_id}`

    return Http.get(url).then(res => {
      runInAction(() => {
        this.init(res.data)
        // this.billList.fetch(res.data.id)
      })
      return res
    })
  }

  updateAccount = (id, body: IUpdateAccountBody) => {
    return Http.put(`/company/account/${id}`, { ...body })
  }
}
