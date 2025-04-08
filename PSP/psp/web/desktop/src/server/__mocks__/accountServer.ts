/* Copyright (C) 2016-present, Yuansuan.cn */

export const detailList = [
  {
    account_balance_contain_freezed: 99899931297,
    account_id: '3N5AecLTkqd',
    amount: 2222,
    bill_sign: 2,
    id: '45UPxPMMsrW',
    out_trade_id: '45UPvYqHY35',
    remark: '作业账单: build',
    trade_id: '45UPxPJrZ4E',
    trade_time: { seconds: 1605666495, nanos: 0 },
    trade_type: 1
  },
  {
    account_balance_contain_freezed: 99899933519,
    account_id: '3N5AecLTkqd',
    amount: 1111,
    bill_sign: 2,
    id: '45Fm74rwTxN',
    out_trade_id: '45Fm4YUAhbL',
    remark: '作业账单: fetch',
    trade_id: '45Fm74ocqaw',
    trade_time: { seconds: 1605255122, nanos: 0 },
    trade_type: 1
  }
]

export const account = {
  id: 'id',
  customer_id: 'customer_id',
  real_customer_id: 'real_customer_id',
  name: 'name',
  currency: 'currency',
  account_balance: 2000000,
  freezed_amount: 1000000,
  normal_balance: 0,
  award_balance: 0,
  withdraw_enabled: false,
  credit_quota: 0,
  status: 0,
  create_time: {
    nanos: 0,
    seconds: 0
  },
  update_time: {
    nanos: 0,
    seconds: 0
  },
  account_balance_contain_freezed: 3000000
}

export const accountServer = {
  __data__: {
    account,
    detailList
  },
  get: async () => ({ data: { account } }),
  getDetailList: async () => ({
    data: {
      list: detailList,
      page_ctx: {
        size: 10,
        index: 1,
        total: 2
      }
    }
  })
}
