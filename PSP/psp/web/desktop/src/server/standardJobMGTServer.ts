/* Copyright (C) 2016-present, Yuansuan.cn */

import { Http } from '@/utils'

export type ListParams = {
  id?: string
  name?: string
  state?: string[]
  submit_start_time?: number
  submit_end_time?: number
  page_index: number
  page_size: number
}

export const standardJobMGTServer = {
  list({
    page_index = 1,
    page_size = 10,
    ...params
  }: ListParams) {
    return Http.get('/standardcompute/job/list', {
      params: {
        index: page_index,
        size: page_size,
        ...params,
      },
    })
  },
  getJob(id) {
    return Http.get(`/standardcompute/job/${id}`)
  },
  getJobStateProcess(id) {
    return Http.get(`/standardcompute/job/process/${id}`)
  },
  cancel(id: string, endpoint: string) {
    return Http.put('/standardcompute/job/action', {
      action: 'CANCEL',
      id,
      endpoint,
    })
  },
  getResidualData(id) {
    return Http.get(`/standardcompute/job/${id}/residual`)
  },
}
