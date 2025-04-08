/* Copyright (C) 2016-present, Yuansuan.cn */

import { Timestamp, PageCtx } from '../common'

describe('@domain/common/Timestamp', () => {
  it('constructor with no params', () => {
    const timestamp = new Timestamp()

    expect(timestamp.seconds).toBe(0)
    expect(timestamp.nanos).toBe(0)
  })

  it('constructor with initial params', () => {
    const timestamp = new Timestamp({
      nanos: 1001,
      seconds: 1000,
    })

    expect(timestamp.seconds).toBe(1000)
    expect(timestamp.nanos).toBe(1001)
  })

  it('update', () => {
    const timestamp = new Timestamp()
    timestamp.update({
      nanos: 1001,
      seconds: 1000,
    })

    expect(timestamp.seconds).toBe(1000)
    expect(timestamp.nanos).toBe(1001)
  })

  it('toString', () => {
    const timestamp = new Timestamp()
    expect(timestamp.toString()).toBe('--')

    timestamp.update({
      nanos: 0,
      seconds: 1605747445,
    })

    expect(timestamp.toString()).toBe('2020/11/19 08:57:25')
  })
})

describe('@domain/common/PageCtx', () => {
  it('constructor with no params', () => {
    const pageCtx = new PageCtx()

    expect(pageCtx.index).toBe(1)
    expect(pageCtx.size).toBe(10)
    expect(pageCtx.total).toBe(0)
  })

  it('constructor with initial params', () => {
    const pageCtx = new PageCtx({
      index: 2,
      size: 20,
      total: 100,
    })

    expect(pageCtx.index).toBe(2)
    expect(pageCtx.size).toBe(20)
    expect(pageCtx.total).toBe(100)
  })

  it('update', () => {
    const pageCtx = new PageCtx()
    pageCtx.update({
      index: 2,
      size: 20,
      total: 100,
    })

    expect(pageCtx.index).toBe(2)
    expect(pageCtx.size).toBe(20)
    expect(pageCtx.total).toBe(100)
  })
})
