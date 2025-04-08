/* Copyright (C) 2016-present, Yuansuan.cn */

import { renderHook } from '@testing-library/react-hooks'
import { useModel } from '../store'

describe('store', () => {
  it('init', () => {
    const { result } = renderHook(() => useModel())
    expect(result.current.unblock).toBeFalsy()
  })
})
