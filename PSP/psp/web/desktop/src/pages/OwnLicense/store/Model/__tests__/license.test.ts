/* Copyright (C) 2016-present, Yuansuan.cn */

import { License, BaseLicense } from '../license'
import { byolServer } from '@/server'

jest.mock('@/server')

const initialValue = byolServer['__data__']['license']

describe('../license', () => {
  it('constructor with no params', () => {
    const model = new License()

    expect(model).toMatchObject(new BaseLicense())
  })

  it('constructor with params', () => {
    const model = new License(initialValue)

    expect(model).toMatchObject(initialValue)
  })

  it('update', () => {
    const model = new License()
    model.update(initialValue)

    expect(model).toMatchObject(initialValue)
  })

  it('fetch', async () => {
    const model = new License()
    await model.fetch()

    expect(model).toMatchObject({
      id: initialValue.id,
      ...JSON.parse(initialValue.license),
    })
  })
})
