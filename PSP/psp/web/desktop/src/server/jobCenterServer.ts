/* Copyright (C) 2016-present, Yuansuan.cn */
import { Http, RawHttp } from '@/utils'

export type ListParams = {
  page_index: number
  page_size: number
  user_filters?: string[]
  app_filters?: string[]
  state_filters?: string[]
  core_filters?: string[]
  download_filters?: string[]
  jobset_filters?: string[]
  fuzzy_key?: string
  order_by?: string
  start_seconds?: string
  end_seconds?: string
  job_id?: string
}

export const jobCenterServer = {
  get(id, params: { project_id: string }) {
    return Http.get(`/jobCenter/${id}`, {
      params,
    })
  },
  list({
    page_index = 1,
    page_size = 10,
    start_seconds,
    end_seconds,
    ...params
  }: ListParams) {
    return Http.get('/jobCenter', {
      params: {
        page: {
          index: page_index,
          size: page_size,
        },
        start_seconds: start_seconds || '0',
        end_seconds: end_seconds || '0',
        ...params,
      },
    })
  },
  export(query: Omit<ListParams, 'page_size' | 'page_index'>) {
    return RawHttp.get('/jobCenter/export', {
      params: query,
      responseType: 'blob',
    }).then(response => {
      return window.URL.createObjectURL(new Blob([response.data]))
    })
  },
  getFilters() {
    return Http.get('/jobCenter/filters')
  },
}
