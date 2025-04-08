/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { FileFactory } from '@/utils/FileSystem'
import { IRequest, formatRequest } from './common'
import { Directory } from './Directory'

export class File extends FileFactory<Directory> {
  constructor(props: Partial<IRequest>) {
    if (props.is_dir) {
      throw new Error(`Can't constructor File with props which is_dir is true`)
    }

    super(formatRequest(props))
  }

  getContent() {
    return new Promise(resolve => {
      setTimeout(() => resolve('hello world!'), 1000)
    })
  }
}
