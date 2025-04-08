/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'

import { Timestamp } from '@/domain/common'

export interface IRequest {
  id: string
  name: string
  company_id: string
  type: number
  status: number
  create_uid: string
  create_name: string
  modify_uid: string
  modify_name: string
  create_time: {
    seconds: number
    nanos: number
  }
  update_time: {
    seconds: number
    nanos: number
  }
}

interface IRole extends IRequest {
  create_time: Timestamp
  update_time: Timestamp
}

export class Role implements IRole {
  @observable id: string
  @observable name: string
  @observable company_id: string
  @observable type: number
  @observable status: number
  @observable create_uid: string
  @observable create_name: string
  @observable modify_uid: string
  @observable modify_name: string
  @observable create_time: Timestamp
  @observable update_time: Timestamp

  constructor(props?: IRequest) {
    if (props) {
      this.init(props)
    }
  }

  @action
  init = (props: IRequest) => {
    Object.assign(this, props)
  }
}
