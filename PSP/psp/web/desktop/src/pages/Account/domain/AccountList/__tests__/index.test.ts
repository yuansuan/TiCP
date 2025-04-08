/* Copyright (C) 2016-present, Yuansuan.cn */

import { DetailList } from '../index'
import { accountServer } from '@/server'
import { Detail } from '../Detail'

jest.mock('@/server')
jest.mock('../Detail')

const initialValue = accountServer['__data__']['detailList']
const initialPageCtx = {
  index: 1,
  size: 10,
  total: 0,
}

describe('@domain/AccountDetailList', () => {
  it('constructor with no params', () => {
    const model = new DetailList()

    expect(model.list).toEqual([])
    expect(model.page_ctx).toEqual(initialPageCtx)
  })

  it('update', () => {
    const model = new DetailList()

    model.update({
      list: initialValue,
      page_ctx: {
        index: 2,
        size: 20,
        total: 2,
      },
    })

    expect(Detail['mock'].calls.length).toEqual(initialValue.length)
    expect(Detail['mock'].calls[0][0]).toEqual(initialValue[0])
    expect(Detail['mock'].calls[1][0]).toEqual(initialValue[1])
    expect(model.page_ctx).toEqual({
      index: 2,
      size: 20,
      total: 2,
    })
  })

  it('constructor with params', () => {
    const model = new DetailList({
      list: initialValue,
      page_ctx: initialPageCtx,
    })

    expect(Detail['mock'].calls.length).toEqual(initialValue.length)
    expect(Detail['mock'].calls[0][0]).toEqual(initialValue[0])
    expect(Detail['mock'].calls[1][0]).toEqual(initialValue[1])
    expect(model.page_ctx).toEqual(initialPageCtx)
  })

  it('fetch will call update', async () => {
    const model = new DetailList()
    const update = jest.spyOn(model, 'update')
    await model.fetch({
      account_id: 'test',
      start_seconds: 1000,
      end_seconds: 3000,
      page_size: 10,
      page_index: 1,
    })

    expect(update).toHaveBeenCalled()
    expect(update.mock.calls[0][0]).toEqual({
      list: initialValue,
      page_ctx: {
        index: 1,
        size: 10,
        total: initialValue.length,
      },
    })
  })
})
