/* Copyright (C) 2016-present, Yuansuan.cn */

import { Company, BaseCompany } from '../Company'

jest.mock('@/server')

const initialValue = {
  id: '3NjBKWJ4rby',
  name: '测试企业',
  account_id: '3NjBKWNSRpW',
  roles: [],
  cloud_type: 'mixed',
  live_chat_id: '',
  box: '',
}

describe('@domain/Company', () => {
  it('constructor with no params', () => {
    const model = new Company()

    expect(model).toMatchObject(new BaseCompany())
  })

  it('update', () => {
    const model = new Company()

    model.update(initialValue)

    expect(model).toMatchObject(initialValue)
  })

  it('constructor with params call update', () => {
    const model = new Company(initialValue)

    expect(model).toMatchObject(initialValue)
  })

  it('isMixed', async () => {
    const model = new Company()

    expect(model.isMixed).toBeFalsy()
    model.update({
      cloud_type: 'mixed',
    })
    expect(model.isMixed).toBeTruthy
  })
})
