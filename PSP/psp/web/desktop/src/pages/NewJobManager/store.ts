/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { JobList } from '@/domain/JobList'
import { jobServer } from '@/server'
import { ListParams } from '@/server/jobServer'
import { runInAction } from 'mobx'
// import { currentUser } from '@/domain'

export const initialQuery = {
  user_names: [],
  app_names: [],
  states: [],
  queues: [],
  job_types: [],
  start_time: null,
  end_time: null,
  job_id: '',
  job_name: '',
  // is_admin: currentUser.hasSysMgrPerm, // 后端同学自行处理了
  project_ids: [],
  job_set_id: '',
  job_set_names: []
}

export function useModel() {
  return useLocalStore(() => {
    const query = initialQuery

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
      query,
      setQuery(query) {
        this.query = query
      },
      pageIndex: 1,
      setPageIndex(index) {
        this.pageIndex = index
      },
      pageSize: 10,
      setPageSize(size) {
        this.pageSize = size
      },
      order: '',
      setOrder(order) {
        this.order = order
      },

      async fetch(params: ListParams) {
        const { data } = await jobServer.list(params)

        runInAction(() => {
          this.model.update({
            list: data.jobs,
            page_ctx: data.page
          })
        })
      },

      async refresh() {
        try {
          this.setLoading(true)
          await this.fetch({
            ...this.query,
            page_index: this.pageIndex,
            page_size: this.pageSize
          })
        } finally {
          this.setLoading(false)
        }
      }
    }
  })
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
