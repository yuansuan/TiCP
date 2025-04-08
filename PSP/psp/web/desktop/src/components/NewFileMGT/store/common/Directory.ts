/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { DirectoryFactory } from '@/utils/FileSystem'
import { IRequest, formatRequest } from './common'
import { File } from './File'

export class Directory extends DirectoryFactory<Directory, File> {
  constructor(props: Partial<IRequest>) {
    if (!props.is_dir) {
      throw new Error(
        `Can't constructor Directory with props which is_dir is false`
      )
    }

    super(formatRequest(props))
  }
}
