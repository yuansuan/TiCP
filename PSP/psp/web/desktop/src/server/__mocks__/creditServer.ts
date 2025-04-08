/* Copyright (C) 2016-present, Yuansuan.cn */

const initialValue = [
  {
    amount: 10000000,
    id: '45VE742acE1',
    normal_amount: 10000000,
    pay_type: 303,
    reward_amount: 0,
    status: 2,
    create_time: { seconds: 1605693824, nanos: 0 },
  },
  {
    amount: 10000000,
    create_time: { seconds: 1603854562, nanos: 0 },
    id: '44TtEZwwu4S',
    normal_amount: 10000000,
    pay_type: 303,
    reward_amount: 0,
    status: 1,
  },
]

export const creditServer = {
  get: async () => ({
    data: {
      list: initialValue,
      page_ctx: {
        index: 1,
        size: 10,
        total: 2,
      },
    },
  }),
  initialValue,
}
