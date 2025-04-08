/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { useLocalStore } from 'mobx-react-lite'
import { runInAction } from 'mobx'
import { createStore } from '@/utils/reducer'
import { DepartmentList } from '@/domain/DepartmentList'

export const useModel = () =>
  useLocalStore(() => ({
    loading: true,
    setLoading(bool) {
      this.loading = bool
    },
    list: new DepartmentList(),
    async fetch() {
      try {
        this.setLoading(true)
        this.list.fetch()
      } catch (e) {
        runInAction(() => {
          this.setLoading(false)
        })
      }
    },
  }))

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
