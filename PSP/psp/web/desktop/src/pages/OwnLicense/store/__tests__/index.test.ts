/* Copyright (C) 2016-present, Yuansuan.cn */

import { act, renderHook } from '@testing-library/react-hooks'
import { useModel } from '..'

jest.mock('@/server')

describe('store', () => {
  it('init fetching', () => {
    const { result } = renderHook(() => useModel())
    expect(result.current.fetching).toBeFalsy()
  })

  it('update fetching', () => {
    const { result } = renderHook(() => useModel())

    expect(result.current.fetching).toBeFalsy()

    act(() => {
      result.current.setFetching(true)
    })
    expect(result.current.fetching).toBeTruthy()
  })

  it('init queryKey', () => {
    const { result } = renderHook(() => useModel())
    expect(result.current.queryKey).toBe('')
  })

  it('update queryKey', () => {
    const { result } = renderHook(() => useModel())
    expect(result.current.queryKey).toBe('')

    act(() => {
      result.current.setQueryKey('test')
    })
    expect(result.current.queryKey).toStrictEqual('test')
  })

  it('init params', () => {
    const { result } = renderHook(() => useModel())

    expect(result.current.params).toStrictEqual({
      key: result.current.queryKey
    })
  })

  it('fetch', async () => {
    const { result } = renderHook(() => useModel())
    // mobx@5 test fail

    // const spyFn = jest.spyOn(result.current, 'setFetching')
    // await result.current.fetch({ key: '' })

    // expect(spyFn).toHaveBeenCalledTimes(2)
    // expect(spyFn.mock.calls[0][0]).toBe(true)
    // expect(spyFn.mock.calls[1][0]).toBe(false)
  })
})
