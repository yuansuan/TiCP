/* Copyright (C) 2016-present, Yuansuan.cn */
import { pageStateStore, PAGE_STATE_KEY } from '../index'

const STORE_KEY = 'PAGE_STATE_STORE_ITEM'
const testValue = { item_01: 1, item_02: 2 }

describe('utils/pageStateStore', () => {
  afterEach(() => {
    window.sessionStorage.removeItem(PAGE_STATE_KEY)
  })

  it('pageStateStore use sessionStorage', () => {
    expect(window.sessionStorage.getItem(PAGE_STATE_KEY)).toBeNull()
    pageStateStore.set({})
    expect(window.sessionStorage.getItem(PAGE_STATE_KEY)).toBeTruthy()
  })

  it('get specific item', () => {
    expect(pageStateStore.get()).toBeNull()
    window.sessionStorage.setItem(PAGE_STATE_KEY, JSON.stringify(testValue))
    expect(pageStateStore.get('item_01')).toEqual(1)
  })

  it('get all', () => {
    expect(pageStateStore.get()).toBeNull()
    window.sessionStorage.setItem(PAGE_STATE_KEY, JSON.stringify(testValue))
    expect(pageStateStore.get()).toEqual(testValue)
  })

  it('set specific item', () => {
    expect(pageStateStore.get(STORE_KEY)).toBeNull()
    pageStateStore.set(STORE_KEY, 'test_value')
    expect(pageStateStore.get(STORE_KEY)).toEqual('test_value')
  })

  it('set all', () => {
    expect(pageStateStore.get()).toBeNull()
    pageStateStore.set(testValue)
    expect(pageStateStore.get()).toEqual(testValue)
  })

  it('remove specific item', () => {
    expect(pageStateStore.get()).toBeNull()
    pageStateStore.set(testValue)
    pageStateStore.remove('item_01')
    expect(pageStateStore.get('item_01')).toBeNull()
    expect(pageStateStore.get('item_02')).toEqual(2)
  })

  it('remove all', () => {
    expect(pageStateStore.get()).toBeNull()
    pageStateStore.set(testValue)
    pageStateStore.remove()
    expect(pageStateStore.get('item_01')).toBeNull()
    expect(pageStateStore.get('item_02')).toBeNull()
  })

  it('getByPath use window.location.pathname', () => {
    expect(pageStateStore.getByPath()).toBeNull()
    pageStateStore.set(window.location.pathname, 'test_value')
    expect(pageStateStore.getByPath()).toEqual('test_value')
  })

  it('setByPath use window.location.pathname', () => {
    expect(pageStateStore.get(window.location.pathname)).toBeNull()
    pageStateStore.setByPath('test_value')
    expect(pageStateStore.get(window.location.pathname)).toEqual('test_value')
  })
})
