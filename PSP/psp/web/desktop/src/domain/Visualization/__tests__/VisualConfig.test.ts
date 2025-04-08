/* Copyright (C) 2016-present, Yuansuan.cn */

import VisualConfig from '../VisualConfig'
import { initialValues } from '@/server/__mocks__/visualServer'
import { companyList } from '@/domain'

jest.mock('@/server')

jest.mock('@/domain', () => {
  return {
    companyList: {
      list: [{ id: '1' }],
      currentId: '1',
      get current() {
        return this.list.find(item => item.id === this.currentId)
      },
    },
  }
})

let visualConfig: VisualConfig = null

describe('visual config -> show visial app', () => {
  beforeEach(() => {
    visualConfig = new VisualConfig(initialValues)
    companyList.currentId = null
  })
  it('init', () => {
    expect(visualConfig.showVisualizeApp).toBeFalsy()
  })

  it('company = null, is open = true: expect false', () => {
    visualConfig.update({
      isOpen: true,
    })
    expect(visualConfig.showVisualizeApp).toBeFalsy()
  })
})
