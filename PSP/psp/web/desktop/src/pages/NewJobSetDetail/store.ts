/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { jobServer } from '@/server'
import { runInAction } from 'mobx'
import { getUrlParams } from '@/utils'
import { JobSet } from '@/domain/JobList/JobSet'
import { JobList } from '@/domain/JobList'
export type ModelOptions = {
  id: string
}

export function useModel({ id }: Partial<ModelOptions>) {
  return useLocalStore(() => {
    const params = getUrlParams()
    return {
      loading: false,
      setLoading(flag) {
        this.loading = flag
      },
      jobSet: new JobSet(),
      model: new JobList(),
      selectedKeys: [],
      setSelectedKeys(keys) {
        this.selectedKeys = keys
      },
      expandedRowKeys: [],
      setExpandedRowKeys(val) {
        this.expandedRowKeys = val
      },
      searchKey: '',
      setSearchKey(key) {
        this.searchKey = key
      },
      pageIndex: 1,
      setPageIndex(index) {
        this.pageIndex = index
      },
      pageSize: 10,
      setPageSize(size) {
        this.pageSize = size
      },
      params: params,
      async refresh() {
        if (!id) return
        try {
          this.setLoading(true)

          const { data } = await jobServer.getJobSet(id)

          runInAction(() => {
            this.jobSet.update(data?.job_set_info)
            this.model.update({ list: data?.job_list })
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
