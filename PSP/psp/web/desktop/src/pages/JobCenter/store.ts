/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { JobList } from '@/domain/JobList'
import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { jobCenterServer } from '@/server'
import { runInAction } from 'mobx'
import moment from 'moment'
import { pageStateStore } from '@/utils'

export const initialQuery = {
  user_filters: [],
  state_filters: [],
  start_seconds: moment().startOf('month').unix() + '',
  end_seconds: moment().endOf('month').unix() + '',
  job_id: '',
}

export function useModel() {
  // restore state from pageStateStore
  const pageState = pageStateStore.getByPath<{
    query: typeof initialQuery
    pageIndex: number
    pageSize: number
  }>()

  return useLocalStore(() => {
    const query = pageState?.query || initialQuery

    return {
      model: new JobList(),
      loading: false,
      setLoading(flag) {
        this.loading = flag
      },
      selectedKeys: [],
      setSelectedKeys(keys) {
        this.selectedKeys = keys
      },
      pageIndex: pageState?.pageIndex || 1,
      setPageIndex(index) {
        this.pageIndex = index
      },
      pageSize: pageState?.pageSize || 10,
      setPageSize(size) {
        this.pageSize = size
      },
      query,
      setQuery(query) {
        this.query = query
      },
      totalAmount: 0,
      setTotalAmount(amount) {
        this.totalAmount = amount
      },
      async refresh() {
        try {
          this.setLoading(true)
          const { data } = await jobCenterServer.list({
            page_index: this.pageIndex,
            page_size: this.pageSize,
            ...this.query,
          })
          runInAction(() => {
            this.model.update({
              list: data.jobs,
              page_ctx: data.page_ctx,
            })
            this.setTotalAmount(data.total_amount)
          })
        } finally {
          this.setLoading(false)
        }
      },
    }
  })
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
