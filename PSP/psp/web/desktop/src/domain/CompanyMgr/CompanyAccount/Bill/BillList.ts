import { observable, runInAction, action } from 'mobx'
import { Http } from '@/utils'
import { Bill } from './Bill'
import moment from 'moment'

export const AccountBillSign_MAP = {
  AccountBillUnknow: '未知',
  AccountBillAdd: '收入',
  AccountBillReduce: '支出',
  AccountBillFreeze: '冻结',
  AccountBillUnfreeze: '解冻',
}

export const AccountBillTradeType_MAP = {
  AccountBillTradeUnknow: '未知',
  AccountBillTradePay: '支付',
  AccountBillTradeCredit: '充值',
  AccountBillTradeRefund: '退款',
  AccountBillTradeWithdraw: '提现',
  AccountBillTradeFundAdd: '加款',
  AccountBillTradeFundSub: '扣款',
}
interface IBillList {
  list: Map<string, Bill>
}

export default class BillList implements IBillList {
  // 获取最近 30 天的 bill
  @observable dates = [moment().subtract(30, 'days'), moment()].map(m =>
    m.valueOf()
  )
  @observable list = new Map()
  @observable index = 1
  @observable size = 10
  @observable totals = 0

  @action
  updateIndex(current: number) {
    this.index = current
  }

  @action
  updateSize(current: number, size: number) {
    this.index = current
    this.size = size
  }

  get = id => this.list.get(id)

  fetch = account_id => {
    let url = `/company/account/bills?account_id=${account_id}&index=${this.index}&size=${this.size}&start_time=${this.dates[0]}&end_time=${this.dates[1]}`

    return Http.get(url).then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data?.list?.map(item => {
            return [item.id, new Bill(item)]
          })
        )
        this.totals = res.data?.total || 0
      })
      return res
    })
  };

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
