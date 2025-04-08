/* Copyright (C) 2016-present, Yuansuan.cn */

import { Model, BaseModel } from '../index'
import { Item } from '../Item'
import { byolServer } from '@/server'

jest.mock('@/server')

const initialValue = byolServer['__data__']['data']

describe('../index', () => {
  it('constructor with no params', () => {
    const model = new Model()

    expect(model).toMatchObject(new BaseModel())
  })

  it('constructor with params ', () => {
    const model = new Model({ list: initialValue })

    expect(model.list).toEqual(initialValue.map(item => new Item(item)))
  })

  it('update', () => {
    const model = new Model()
    model.update({ list: initialValue })

    expect(model.list).toEqual(initialValue.map(item => new Item(item)))
  })

  it('fetch', async () => {
    const model = new Model()
    await model.fetch({ key: '' })

    expect(model.list).toEqual(initialValue.map(item => new Item(item)))
  })
})
