/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { v4 as uuid } from 'uuid'
import { observable, action } from 'mobx'
import { BaseDirectory } from './BaseDirectory'
import { Tree } from '../Tree'
import { IFile } from './typing'

export class FileFactory<T extends Tree = Tree, K extends IFile<T> = IFile<T>>
  implements IFile<T>
{
  readonly isFile = true
  readonly id: string = uuid()
  @observable type: string
  @observable name: string
  @observable is_dir: boolean = false
  @observable is_sym_link: boolean = false
  @observable mode: string = ''
  @observable path: string
  @observable size: number
  @observable parent: T | null = null
  @observable mtime: number

  constructor(props?: Partial<K>) {
    if (props) {
      this.update(props as any)
    }
  }

  @action
  update = (props: Partial<Omit<K, 'id'>>) => {
    // 防止外部 id 覆盖
    Reflect.deleteProperty(props, 'id')
    Object.assign(this, props)
  }
}

/**
 * 文件管理模块：文件基类
 * @class BaseFile
 */
export class BaseFile extends FileFactory<BaseDirectory> {}
