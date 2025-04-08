/* Copyright (C) 2016-present, Yuansuan.cn */

export type FilterParams = Partial<{
  user_id: string
  user_name: string
  merchandise_id: string
  month: string
  out_resource_type: number
}>

export const billList = [
  {
    account_id: '3N5AecLTkqd',
    amount: 2222,
    billing_desc:
      '{"start_time":"2020-11-18 10:27:24 +0800 CST","end_time":"2020-11-18 10:28:04 +0800 CST","cpus":1,"cpu_time":40,"app_name":"Sample"}',
    billing_month: '2020-11',
    billing_time: { seconds: 1605666495, nanos: 0 },
    company_id: '3M6WH7N7DUU',
    create_time: { seconds: 1605666449, nanos: 0 },
    discount_amount: 0,
    freeze_amount: 0,
    id: '45UPxPJrZ4E',
    job: null,
    merchandise: {
      create_time: { seconds: 1590129622, nanos: 0 },
      creator_id: '3TX3PTpcv1C',
      id: '3W8XEWhb9uY',
      last_update_id: '3TX3PTpcv1C',
      license_type: '',
      name: '测试商品',
      out_resource_id: '3W8Vhe7duf5',
      out_resource_type: 1,
      product_id: '3VGTYLv6FhU',
      remark: 'Jonathan测试',
      update_time: { seconds: 1590129622, nanos: 0 },
    },
    merchandise_id: '3W8XEWhb9uY',
    merchandise_price_des: '核时',
    merchandise_price_id: '45PYGZVJcwH',
    merchandise_price_unit: 200000,
    merchandise_quantity: 0.011111111111111112,
    name: '作业账单: build',
    out_biz_id: '45UPvYqHY35',
    out_resource_id: '3W8Vhe7duf5',
    out_resource_type: 1,
    product_id: '3VGTYLv6FhU',
    product_name: '作业',
    project_id: '3Z7pmBNmc1m',
    real_amount: 2222,
    remark: 'TODO',
    status: 2,
    update_time: { seconds: 1605666495, nanos: 0 },
    user_id: '3P7iEcshowu',
    user_name: 'bph',
  },
]

export const billServer = {
  __data__: {
    billList,
  },
  getList: () => ({
    data: {
      list: billList,
      page_ctx: {
        index: 1,
        size: 10,
        total: 1,
      },
      total_amount: 2222,
    },
  }),
  export: () => {},
}
