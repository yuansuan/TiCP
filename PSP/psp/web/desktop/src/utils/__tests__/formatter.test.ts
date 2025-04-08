/* Copyright (C) 2016-present, Yuansuan.cn */

const { formatAmount } = require('@/utils/formatter')

describe('formatAmount', () => {
  it('formatAmount(0) 的结果为 0.00', () => {
    expect(formatAmount(0)).toBe('0.00')
  })

  it('formatAmount(100000) 的结果为 1.00', () => {
    expect(formatAmount(100000)).toBe('1.00')
  })

  it('formatAmount(100100) 的结果为 1.00', () => {
    expect(formatAmount(100100)).toBe('1.00')
  })

  it('formatAmount(101000) 的结果为 1.01', () => {
    expect(formatAmount(101000)).toBe('1.01')
  })
})
