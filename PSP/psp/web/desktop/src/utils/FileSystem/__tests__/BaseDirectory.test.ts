/* Copyright (C) 2016-present, Yuansuan.cn */

import { BaseDirectory } from '../BaseDirectory'
import { BaseFile } from '../BaseFile'

const initialValue = {
  id: 'test_id',
  type: 'txt',
  name: 'test_name',
  path: 'test_path',
  size: 1024,
  parent: new BaseDirectory(),
  children: [new BaseFile(), new BaseFile()],
  mtime: 1000,
}

describe('utils/BaseDirectory', () => {
  it('BaseDirectory constructor with no params', () => {
    const baseDirectory = new BaseDirectory()

    expect(baseDirectory.id).toBeTruthy()
    expect(baseDirectory.parent).toBeNull()
    expect(baseDirectory.isFile).toBeFalsy()
    expect(baseDirectory.children).toEqual([])
    expect(baseDirectory.type).toEqual('FOLDER')
    expect(baseDirectory.name).toBeUndefined()
    expect(baseDirectory.path).toBeUndefined()
    expect(baseDirectory.size).toBeUndefined()
    expect(baseDirectory.mtime).toBeUndefined()
  })

  it('BaseDirectory constructor with params', () => {
    const baseDirectory = new BaseDirectory(initialValue)

    expect(baseDirectory.isFile).toBeFalsy()
    expect(baseDirectory.type).toEqual(initialValue.type)
    expect(baseDirectory.name).toEqual(initialValue.name)
    expect(baseDirectory.path).toEqual(initialValue.path)
    expect(baseDirectory.size).toEqual(initialValue.size)
    expect(baseDirectory.parent).toEqual(initialValue.parent)
    expect(baseDirectory.children).toEqual(initialValue.children)
    expect(baseDirectory.mtime).toEqual(initialValue.mtime)
  })

  it('update', () => {
    const baseDirectory = new BaseDirectory()

    baseDirectory.update(initialValue)

    expect(baseDirectory.isFile).toBeFalsy()
    expect(baseDirectory.type).toEqual(initialValue.type)
    expect(baseDirectory.name).toEqual(initialValue.name)
    expect(baseDirectory.path).toEqual(initialValue.path)
    expect(baseDirectory.size).toEqual(initialValue.size)
    expect(baseDirectory.parent).toEqual(initialValue.parent)
    expect(baseDirectory.children).toEqual(initialValue.children)
    expect(baseDirectory.mtime).toEqual(initialValue.mtime)
  })

  it('can not override id', () => {
    const baseDirectory = new BaseDirectory({
      id: initialValue.id,
    })

    expect(baseDirectory.id).not.toEqual(initialValue.id)
  })

  it('children proxy _children', () => {
    const baseDirectory = new BaseDirectory()

    expect(baseDirectory._children).toEqual([])
    expect(baseDirectory.children).toEqual([])
    baseDirectory.children = initialValue.children
    expect(baseDirectory._children).toEqual(initialValue.children)
    expect(baseDirectory.children).toEqual(initialValue.children)
  })
})
