/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Tree } from '../Tree'
import { FileFactory } from './BaseFile'

export type ChildNode<
  T extends Tree = Tree,
  P extends FileFactory<T, IFile<T>> = FileFactory<T, IFile<T>>
> = T | P | null

export type IDirectory<
  T extends Tree = Tree,
  P extends FileFactory<T, IFile<T>> = FileFactory<T, IFile<T>>
> = {
  isFile: boolean
  id: string
  name: string
  type: string
  path: string
  size: number
  parent: T | null
  children: ChildNode<T, P>[]
  mtime: number
}

export type IFile<T extends Tree = Tree> = {
  isFile: boolean
  id: string
  name: string
  path: string
  type: string
  size: number
  parent: T | null
  mtime: number
}
