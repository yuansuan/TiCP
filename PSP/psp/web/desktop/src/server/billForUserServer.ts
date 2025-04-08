/* Copyright (C) 2016-present, Yuansuan.cn */

import { gql } from '@apollo/client'
import { apolloClient, RawHttp } from '@/utils'
import { env } from '@/domain'

class JobBillList {
  user_id: String
  company_id: String
  merchandise_id: String
  project_id: String
  billing_month: String
  pageIndex: Number
  pageSize: Number
  types: Number[]
}

export type FilterParams = Partial<{
  types: number[]
  merchandise_id: string
  billing_month: string
}>

export const BILLUSER_LIST = gql`
  query jobBillList($params: JobBillInput!) {
    jobBillList(params: $params) {
      list {
        bill {
          id
          billing_month
          user_id
          user_name
          update_time
          merchandise_price_unit
          merchandise_quantity
          real_amount
          out_resource_type
          out_biz_id
          job_id
          refund_amount
        }
        job {
          job {
            id
          }
        }
        merchandise_name
      }
      page_ctx {
        index
        size
        total
      }
      total_amount
      total_refund_amount
    }
  }
`

export const billUserServer = {
  async getBillUserList(params: JobBillList) {
    const {
      data: {
        jobBillList: { list, page_ctx, total_amount, total_refund_amount }
      }
    } = await apolloClient.query({
      query: BILLUSER_LIST,
      variables: { params },
      fetchPolicy: 'network-only'
    })

    return { list, page_ctx, total_amount, total_refund_amount }
  },
  export: (params: FilterParams) =>
    RawHttp.get('/bill/job/export', {
      params: {
        ...params,
        company_id: !env.isPersonal ? env?.company?.id : '1'
      },
      responseType: 'blob'
    })
}
