/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import BoxHttp from '@/domain/Box/BoxHttp'
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

  async fetch({
    id,
    path,
    project_id
  }: {
    id: string
    path: string
    project_id: string
  }) {
    let res = null

    try {
      res = await Http.get(`/file/ls`, {
        params: {
          path
        },
        disableErrorMessage: true // 不显示错误
      })
    } catch (e) {
      console.error(e)
    }

    let files = res?.data?.files || []
    let fileNames = files.map(file => file.name)

    let file_result_list = []

    try {
      let res = await Http.post(`/standardcompute/job/file/classify/${id}`, {
        fileNames
      })
      file_result_list = res?.data?.file_result_list || []
    } catch (e) {
      console.error(e)
    }

    let types = {}

    fileNames.map(name => {
      file_result_list.forEach(rule => {
        if (rule.rule_name !== 'all' && rule.file_name.includes(name)) {
          types[name] = rule.rule_name
        }
      })
    })

    files = files
      .map(file => ({ ...file, type: JobFileTypeEnum[types[file.name]] }))
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
