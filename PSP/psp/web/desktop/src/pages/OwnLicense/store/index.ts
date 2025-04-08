/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { Model, FetchParams } from './Model'

export function useModel() {
  return useLocalStore(() => ({
    model: new Model(),
    fetching: false,
    setFetching(flag) {
      this.fetching = flag
    },
    queryKey: '',
    setQueryKey(key) {
      this.queryKey = key
    },
    get params() {
      return {
        key: this.queryKey,
      }
    },
    async fetch(params?: FetchParams) {
      try {
        this.setFetching(true)
        await this.model.fetch(params)
      } finally {
        this.setFetching(false)
      }
    },
  }))
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
