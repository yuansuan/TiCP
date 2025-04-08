/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { useParams } from 'react-router'
import { JobList } from '@/domain/JobList'
import { jobServer } from '@/server'
import { runInAction } from 'mobx'
import { JobSet } from './JobSet'

export function useModel() {
  const { id } = useParams<{ id: string }>()
  return useLocalStore(() => ({
    model: new JobList(),
    jobSet: new JobSet(),
    loading: false,
    setLoading(flag) {
      this.loading = flag
    },
    selectedKeys: [],
    setSelectedKeys(keys) {
      this.selectedKeys = keys
    },
    pageIndex: 1,
    setPageIndex(index) {
      this.pageIndex = index
    },
    pageSize: 10,
    setPageSize(size) {
      this.pageSize = size
    },
    async refresh() {
      try {
        this.setLoading(true)
        const { data } = await jobServer.getJobSetDetail({
          id,
          page_index: this.pageIndex,
          page_size: this.pageSize,
        })
        runInAction(() => {
          this.model.update({
            list: data.jobs,
            page_ctx: data.page_ctx,
          })
          this.jobSet.update(data.job_set)
        })
      } finally {
        this.setLoading(false)
      }
    },
  }))
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
