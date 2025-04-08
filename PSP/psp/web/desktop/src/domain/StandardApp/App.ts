/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable } from 'mobx'
import { Timestamp } from '../common'


const STATUS_MAP = {
  unpublished: '未发布',
  published: '已发布',
}

export interface IApp {

  id: string

  name: string

  version: string

  image: string

  endpoint: string

  command: string

  publish_status: string

  description: string

  icon_url: string

  cores_max_limit: number

  cores_placeholder: string

  create_time: Timestamp

  update_time: Timestamp
}

export default class App implements IApp {
  @observable id = ''
  @observable name = ''
  @observable version = ''
  @observable image = null
  @observable endpoint = null
  @observable command = null
  @observable description = null
  @observable publish_status = 'unpublished'
  @observable icon_url = null
  @observable cores_max_limit = null
  @observable cores_placeholder = null
  @observable create_time = null
  @observable update_time = null

  constructor(data?: Partial<IApp>) {
    this.init(data)
  }

  async init(data) {
    if (data) {
      Object.assign(this, {
        ...data,
        create_time: new Timestamp(data?.create_time),
        update_time: new Timestamp(data?.update_time)
      })
    }
  }

  get publish_status_str() {
    return STATUS_MAP[this.publish_status]
  }
}
