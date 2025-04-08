/* Copyright (C) 2016-present, Yuansuan.cn */

import { SCList } from '../SCList'

describe('domain/SCList', () => {
  it('constructor with no params', () => {
    const model = new SCList()

    expect(model.list).toHaveLength(0)
  })

  it('constructor with initial params', () => {
    const model = new SCList({
      list: [
        {
          sc_id: '1',
          tier_name: 'sc_01',
        },
        {
          sc_id: '2',
          tier_name: 'sc_02',
        },
      ],
    })

    expect(model.list).toHaveLength(2)
  })

  it('update list', () => {
    const model = new SCList()
    model.update({
      list: [
        {
          sc_id: '1',
          tier_name: 'sc_01',
        },
        {
          sc_id: '2',
          tier_name: 'sc_02',
        },
      ],
    })

    expect(model.list).toHaveLength(2)
  })

  it('get SC name by id', () => {
    const model = new SCList({
      list: [
        {
          sc_id: '1',
          tier_name: 'sc_01',
        },
        {
          sc_id: '2',
          tier_name: 'sc_02',
        },
      ],
    })

    expect(model.getName('1')).toBe('sc_01')
    expect(model.getName('2')).toBe('sc_02')
  })
})
