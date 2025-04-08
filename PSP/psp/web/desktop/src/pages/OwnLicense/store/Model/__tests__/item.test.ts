/* Copyright (C) 2016-present, Yuansuan.cn */

import { Item, BaseItem } from '../Item'
import { byolServer } from '@/server'

jest.mock('@/server')

const initialValue = byolServer['__data__']['Item']

describe('../Item', () => {
  it('constructor with no params', () => {
    const model = new Item()

    expect(model).toMatchObject(new BaseItem())
  })

  it('constructor with params', () => {
    const model = new Item(initialValue)

    expect(model).toMatchObject(initialValue)
  })

  it('update', () => {
    const model = new Item()
    model.update(initialValue)

    expect(model).toMatchObject(initialValue)
  })
})
