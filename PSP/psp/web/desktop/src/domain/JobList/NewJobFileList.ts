/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { newBoxServer } from '@/server'
import { observable, runInAction } from 'mobx'
import { JobFileTypeEnum } from '@/constant'
import { Http } from '@/utils'
import { currentUser } from '@/domain'

type File = {
  path: string
  is_dir: boolean
  mod_time: number
  size: number
  name: string
  type: JobFileTypeEnum
  children: File[]
  parent: File[]
}

class BaseJobFileList {
  @observable list: File[] = []
}

export class NewJobFileList extends BaseJobFileList {
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
    path,
    cross = false,
    is_cloud = false,
    user_name,
    filter_regexp_list = []
  }: {
    path: string
    cross?: boolean
    is_cloud: boolean
    user_name: string
    filter_regexp_list: string[]
  }) {
    const { data } = await newBoxServer.list({
      path,
      cross,
      is_cloud,
      user_name,
      filter_regexp_list
    })
    let files = data.map(file => ({
      ...file,
      key: file.path,
      children: file.is_dir ? [] : null
    }))

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
