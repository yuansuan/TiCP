/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { JobList } from '@/domain/JobList'

type Model = {
  model: JobList
  loading: Boolean
  setLoading: (flag: Boolean) => void
  pageIndex: number
  setPageIndex: (index: number) => void
  pageSize: number
  setPageSize: (size: number) => void
  selectedKeys: string[]
  setSelectedKeys: (keys: any[]) => void
  refresh: () => void
}

export function useModel(model: Model) {
  return useLocalStore(() => model)
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
