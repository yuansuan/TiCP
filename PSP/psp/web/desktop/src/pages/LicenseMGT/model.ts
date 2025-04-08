/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { useLocalStore } from 'mobx-react-lite'
import { createStore } from '@/utils/reducer'

export const useModel = () =>
  useLocalStore(() => ({
    token: '',
    setToken(token) {
      this.token = token
    },
    url: '',
    setURL(url) {
      this.url = url
    },
  }))

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
