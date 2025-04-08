/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { IFile, FileFactory } from '@/utils/FileSystem'
import { observable, action, computed } from 'mobx'
import { formatByte } from '@/utils/Validator'
import { JobDirectory } from './JobDirectory'
import { UploadFileStatus } from '@/components/Uploader'

type ExtendedProps = {
  uid: string
  percent: number
  status: UploadFileStatus
  realCommonPath?: string // for standard app
}

type IJobFile = IFile<JobDirectory> & ExtendedProps

export class JobFile
  extends FileFactory<JobDirectory, IJobFile>
  implements ExtendedProps
{
  uid: string
  @observable parent: JobDirectory = null
  @observable percent: number = 100
  @observable status: UploadFileStatus = 'done'
  @observable isMain = false
  @observable realCommonPath = null
  isRoot = false

  constructor(props: any) {
    super()
    Object.assign(this, props)
  }

  @computed
  get path() {
    let paths = [this.name]

    let parent = this.parent
    while (parent) {
      paths.unshift(parent.name)
      parent = parent.parent
    }

    return paths.filter(Boolean).join('/')
  }

  @computed
  get displaySize() {
    return formatByte(this.size)
  }

  @action
  update = (props: Partial<IJobFile> = {}) => {
    Object.assign(this, props)
  }
}
