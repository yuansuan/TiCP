/* Copyright (C) 2016-present, Yuansuan.cn */

import { Detail, BaseDetail } from '../Detail'
import { accountServer } from '@/server'

jest.mock('@/server')

const initialValue = accountServer['__data__']['detailList'][0]

describe('domain/AccountDetail', () => {
  it('constructor with no params', () => {
    const model = new Detail()

    expect(model).toMatchObject(new BaseDetail())
  })

  it('update', () => {
    const model = new Detail()

    model.update(initialValue)

    expect(model).toMatchObject(initialValue)
  })

  it('constructor with params', () => {
    const model = new Detail(initialValue)

    expect(model).toMatchObject(initialValue)
  })
})
