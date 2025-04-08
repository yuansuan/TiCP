/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { Job } from '@/domain/JobList/Job'
import { JobFileList } from '@/domain/JobList/JobFileList'
import { jobServer, jobCenterServer } from '@/server'
import { runInAction } from 'mobx'
import { getUrlParams } from '@/utils/Validator'

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
      jobFile: new JobFileList(),
      selectedKeys: [],
      setSelectedKeys(keys) {
        this.selectedKeys = keys
      },
      searchKey: '',
      setSearchKey(key) {
        this.searchKey = key
      },
      async refresh() {
        try {
          this.setLoading(true)

          const { job } = this
          const { data } = await (params?.project_id
            ? jobCenterServer.get(id, {
                project_id: params?.project_id as string
              })
            : jobServer.get(id))

          runInAction(() => {
            this.job.update(data)
          })
          await this.jobFile.fetch({
            id: job.id
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
