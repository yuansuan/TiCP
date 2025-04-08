/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'

import { Timestamp } from '@/domain/common'

export interface IRequest {
  id: string
  name: string
  code: string
  remark: string
  status: number
  create_uid: string
  create_name: string
  modify_uid: string
  modify_name: string
  update_time: {
    seconds: number
    nanos: number
  }
  create_time: {
    seconds: number
    nanos: number
  }
}

export class Perm {
  @observable id: string
  @observable name: string
  @observable code: string
  @observable remark: string
  @observable status: number
  @observable create_uid: string
  @observable create_name: string
  @observable modify_uid: string
  @observable modify_name: string
  @observable update_time: Timestamp
  @observable create_time: Timestamp

  constructor(props: IRequest) {
    if (props) {
      this.init(props)
    }
  }

  @action
  init = (props: IRequest) => {
    Object.assign(this, props)
    this.update_time = new Timestamp(props.update_time)
    this.create_time = new Timestamp(props.create_time)
  }
}
