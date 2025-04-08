/* Copyright (C) 2016-present, Yuansuan.cn */

import { CompanyList } from '../index'
import { Company } from '../Company'
import { companyServer } from '@/server'

jest.mock('@/server')

const initialValue = companyServer['initialValue']

describe('@domain/CompanyList', () => {
  it('constructor with no params', () => {
    const model = new CompanyList()

    expect(model.list).toEqual([])
    expect(model.currentId).toBeUndefined()
  })

  it('update', () => {
    const model = new CompanyList()

    model.update({
      list: initialValue,
      currentId: initialValue[0].id,
    })

    expect(model.list).toEqual(initialValue.map(item => new Company(item)))
    expect(model.currentId).toEqual(initialValue[0].id)
  })

  it('constructor with params call update', () => {
    const model = new CompanyList({
      list: initialValue,
      currentId: initialValue[0].id,
    })

    expect(model.list).toEqual(initialValue.map(item => new Company(item)))
    expect(model.currentId).toEqual(initialValue[0].id)
  })

  it('fetch', async () => {
    const model = new CompanyList()
    await model.fetch()

    expect(model.list).toEqual(initialValue.map(item => new Company(item)))
  })

  it('current', () => {
    const model = new CompanyList({
      list: initialValue,
      currentId: initialValue[0].id,
    })

    expect(model.current).toEqual(model.list[0])
    model.update({
      currentId: initialValue[1].id,
    })
    expect(model.current).toEqual(model.list[1])
  })
})
