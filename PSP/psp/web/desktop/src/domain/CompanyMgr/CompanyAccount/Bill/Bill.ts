import { action, observable } from 'mobx'
import { formatTime } from '@/utils/formatter'

enum AccountBillSign {}

enum AccountBillTradeType {}

interface IBill {
  id?: string | null
  account_id?: string | null
  bill_sign?: AccountBillSign | null
  amount?: number | null
  trade_type?: AccountBillTradeType | null
  trade_id?: string | null
  trade_time?: any | null
  account_balance_contain_freezed?: number | null
  remark?: string | null
  out_trade_id?: string | null
  delta_normal_balance?: number | null
  delta_award_balance?: number | null
}

export class Bill implements IBill {
  @observable id?: string | null
  @observable account_id?: string | null
  @observable bill_sign?: AccountBillSign | null
  @observable amount?: number | null
  @observable trade_type?: AccountBillTradeType | null
  @observable trade_id?: string | null
  @observable trade_time?: any | null = null
  @observable account_balance_contain_freezed?: number | null
  @observable remark?: string | null
  @observable out_trade_id?: string | null
  @observable delta_normal_balance?: number | null
  @observable delta_award_balance?: number | null

  get trade_time_string() {
    return this.trade_time !== null ? formatTime(this.trade_time.seconds) : null
  }

  constructor(request?: IBill) {
    this.init(request)
  }

  @action
  public init = (request?: IBill) => {
    request && Object.assign(this, { ...request })
  }
}
