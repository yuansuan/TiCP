/* Copyright (C) 2016-present, Yuansuan.cn */

import { Http, RawHttp } from '@/utils'
import { Moment } from 'moment'
export type FilterParams = Partial<{
  user_name: string
  merchandise_id: string
  resource_id: string
  job_submit_end_time: number
  job_submit_start_time: number
  page_index: number
  page_size: number
  dates: Moment[]
}>

export const billServer = {
  getList: (params: FilterParams) => {
    const { page_index, page_size, dates, ...query } = params
    return Http.post(
      '/billing/list',
      {
        filter: {
          ...query
        },
        page: {
          index: page_index,
          size: page_size
        }
      }
    )
  },
  export: (params: FilterParams) => {
    return RawHttp.post(
      '/billing/export',
      {
        filter: {
          ...params
        }
        // order_sort: {
        //   order_by: '',
        //   sort_by_asc: true
        // }
      },
      { responseType: 'blob' }
    )
  }
}
