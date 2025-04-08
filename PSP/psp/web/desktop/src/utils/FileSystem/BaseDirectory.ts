/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { v4 as uuid } from 'uuid'
import { Tree } from '../Tree'
import { FileFactory, BaseFile } from './BaseFile'
import { IDirectory, IFile, ChildNode } from './typing'

export class DirectoryFactory<
    T extends Tree = Tree,
    P extends FileFactory<T, IFile<T>> = FileFactory<T, IFile<T>>,
    K extends IDirectory<T, P> = IDirectory<T, P>
  >
  extends Tree
  implements IDirectory<T, P>
{
  readonly isFile: boolean = false
  readonly id: string = uuid()
  @observable type: string = 'FOLDER'
  @observable name: string
  @observable is_dir: boolean = false
  @observable is_text: boolean = false
  @observable mode: string = ''
  @observable path: string
  @observable size: number
  @observable parent: T = null
  @observable _children: ChildNode<T, P>[] = []
  @observable mtime: number

  constructor(props?: Partial<K>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  get children(): ChildNode<T, P>[] {
    return this._children
  }

  set children(children: ChildNode<T, P>[]) {
    this._children = children.map(child => {
      child.parent = this
      return child
    })
  }

  @action
  update = (props: Partial<Omit<IDirectory<T, P>, 'id'>>) => {
    // 防止外部 id 覆盖
    Reflect.deleteProperty(props, 'id')
    Object.assign(this, props)
  }
}

/**
 * 文件管理模块：目录基类
 * @class BaseDirectory
 */
export class BaseDirectory extends DirectoryFactory<BaseDirectory, BaseFile> {}
