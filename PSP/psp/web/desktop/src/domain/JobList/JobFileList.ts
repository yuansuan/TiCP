/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { boxServer } from '@/server'
import { observable, runInAction } from 'mobx'
import { JobFileTypeEnum } from '@/constant'
import { Http } from '@/utils'

type File = {
  is_dir: boolean
  mod_time: number
  size: number
  name: string
  type: JobFileTypeEnum
}

class BaseJobFileList {
  @observable list: File[] = []
}

export class JobFileList extends BaseJobFileList {
  constructor(props?: Partial<BaseJobFileList>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  update = (props: Partial<BaseJobFileList>) => {
    Object.assign(this, props)
  }

  async fetch({ id }: { id: string }) {
    const { data } = await boxServer.list({
      path: id
    })
    let files = data.files
    const { data: types } = await Http.post(`/job/${id}/classifier`, {
      files: files.map(file => file.name),
    })
    files = files
      .map((file, index) => ({ ...file, type: types[index] }))
      .filter(file => !file.is_dir)

    runInAction(() => {
      this.update({
        list: files
      })
    })

    return files
  }

  getName(name: string) {
    const file = this.list.find(item => item.name === name)
    return file
  }
}
