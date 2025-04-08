/* Copyright (C) 2016-present, Yuansuan.cn */

import { BaseFile } from '../BaseFile'
import { BaseDirectory } from '../BaseDirectory'

const initialValue = {
  id: 'test_id',
  type: 'txt',
  name: 'test_name',
  path: 'test_path',
  size: 1024,
  parent: new BaseDirectory(),
  m_date: 1000
}

describe('utils/BaseFile', () => {
  it('BaseFile constructor with no params', () => {
    const baseFile = new BaseFile()

    expect(baseFile.id).toBeTruthy()
    expect(baseFile.parent).toBeNull()
    expect(baseFile.isFile).toBeTruthy()
    expect(baseFile.type).toBeUndefined()
    expect(baseFile.name).toBeUndefined()
    expect(baseFile.path).toBeUndefined()
    expect(baseFile.size).toBeUndefined()
    expect(baseFile.m_date).toBeUndefined()
  })

  it('BaseFile constructor with params', () => {
    const baseFile = new BaseFile(initialValue)

    expect(baseFile.isFile).toBeTruthy()
    expect(baseFile.type).toEqual(initialValue.type)
    expect(baseFile.name).toEqual(initialValue.name)
    expect(baseFile.path).toEqual(initialValue.path)
    expect(baseFile.size).toEqual(initialValue.size)
    expect(baseFile.parent).toEqual(initialValue.parent)
    expect(baseFile.m_date).toEqual(initialValue.m_date)
  })

  it('update', () => {
    const baseFile = new BaseFile()

    baseFile.update(initialValue)

    expect(baseFile.isFile).toBeTruthy()
    expect(baseFile.type).toEqual(initialValue.type)
    expect(baseFile.name).toEqual(initialValue.name)
    expect(baseFile.path).toEqual(initialValue.path)
    expect(baseFile.size).toEqual(initialValue.size)
    expect(baseFile.parent).toEqual(initialValue.parent)
    expect(baseFile.m_date).toEqual(initialValue.m_date)
  })

  it('can not override id', () => {
    const baseFile = new BaseFile({
      id: initialValue.id
    })

    expect(baseFile.id).not.toEqual(initialValue.id)
  })
})
