/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'

export function useModel(
  props?: Partial<{
    defaultEditing: boolean
    defaultValue: string
    defaultError: string
  }>
) {
  return useLocalStore(() => ({
    editing: props?.defaultEditing || false,
    value: props?.defaultValue || '',
    error: props?.defaultError || '',
    setEditing(editing) {
      this.editing = editing
    },
    setValue(value) {
      this.value = value
    },
    setError(error) {
      this.error = error
    },
  }))
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
