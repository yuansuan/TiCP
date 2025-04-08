/* Copyright (C) 2016-present, Yuansuan.cn */

export const PAGE_STATE_KEY = 'PLATFORM_PAGE_STATE'

function storeSetter(value: object): object
function storeSetter(key: string, value: any): object
function storeSetter(key: string | object, value?: any): object {
  if (typeof key === 'object') {
    value = key
    window.sessionStorage.setItem(PAGE_STATE_KEY, JSON.stringify(value))
    return value
  } else {
    let state = {}

    try {
      state = JSON.parse(window.sessionStorage.getItem(PAGE_STATE_KEY)) || {}
    } catch (err) {}
    state[key] = value

    window.sessionStorage.setItem(PAGE_STATE_KEY, JSON.stringify(state))
    return value
  }
}

type Store = {
  get: <T = any>(key?: string) => T
  set: typeof storeSetter
  remove: (key?: string) => object
  getByPath: (key?: string) => object
  setByPath: (key?: any, value?: any) => object
}

const store: Store = {
  get(key) {
    let state = null

    try {
      state = JSON.parse(window.sessionStorage.getItem(PAGE_STATE_KEY))
    } catch (err) {}

    return key ? (state || {})[key] || null : state
  },
  set(key, value?) {
    return storeSetter(key, value)
  },
  remove(key?) {
    let state = store.get()
    if (key) {
      Reflect.deleteProperty(state, key)
      store.set(state)
      return state
    } else {
      window.sessionStorage.removeItem(PAGE_STATE_KEY)
      return null
    }
  },
  getByPath(key) {
    return store.get(key)
  },
  setByPath(key, value) {

    return store.set(key, value)
  }
}

export const pageStateStore = store
