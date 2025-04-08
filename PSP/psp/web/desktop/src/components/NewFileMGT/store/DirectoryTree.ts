/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { BaseDirectory } from '@/utils/FileSystem'

export class DirectoryTree extends BaseDirectory {
  // @ts-ignore
  @observable children: BaseDirectory[] = []

  @action
  setChildren = (
    props: Array<
      Partial<{
        name: string
        path: string
      }>
    >
  ) => {
    this.children = [...props].map(
      item => new BaseDirectory({ ...item, parent: this })
    )
  }
}
