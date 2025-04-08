/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { Job } from '@/domain/JobList/Job'
import { NewJobFileList } from '@/domain/JobList/NewJobFileList'
import { jobServer } from '@/server'
import { runInAction } from 'mobx'
import { getUrlParams } from '@/utils'
export type ModelOptions = {
  tabKey: 'jobs' | 'jobSets'
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
      job: new Job(),
      jobFile: new NewJobFileList(),
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
      params: params,
      async refresh() {
        if (!id) return
        try {
          this.setLoading(true)

          const { job } = this
          const { data } = await jobServer.get(id)

          runInAction(() => {
            this.job.update(data)
          })
          return await this.jobFile.fetch({
            path: job.work_dir,
            is_cloud: job?.isSyncToLocal ? false : job?.isCloud,
            user_name: job.user_name,
            filter_regexp_list: job?.file_filter_regs
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
