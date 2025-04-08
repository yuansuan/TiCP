/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, computed } from 'mobx'

import { Timestamp } from '@/domain/common'
import { Role, IRequest as IRoleRequest } from './Role'

export interface IRequest {
  user_id: string
  company_id: string
  real_name: string
  phone: string
  email: string
  account_id: string
  status: number
  role_list: Array<IRoleRequest>
  consume_limit: number
  last_login_time: {
    seconds: number
    nanos: number
  }
  create_time: {
    seconds: number
    nanos: number
  }
  update_time: {
    seconds: number
    nanos: number
  }
  department: {
    id: string
    name: string
  }
}

export interface IUser extends Omit<IRequest, 'role_list'> {
  role_list: Role[]
  create_time: Timestamp
  update_time: Timestamp
  last_login_time: Timestamp
}

export class User implements IUser {
  @observable user_id: string
  @observable company_id: string
  @observable real_name: string
  @observable user_name: string
  @observable phone: string
  @observable email: string
  @observable account_id: string
  @observable status: number
  @observable role_list: Role[]
  @observable create_time: Timestamp
  @observable update_time: Timestamp
  @observable last_login_time: Timestamp
  @observable consume_limit: number
  @observable display_user_name: string
  @observable department: { id: string; name: string }

  constructor(props?: IRequest) {
    if (props) {
      this.init(props)
    }
  }

  @computed
  get displayName() {
    return this.real_name || this.display_user_name || this.phone
  }

  @action
  init(props: IRequest) {
    Object.assign(this, props)
    this.consume_limit = props.consume_limit && props.consume_limit / 100000
    this.role_list = props.role_list.map(item => new Role(item))
    this.create_time = new Timestamp(props.create_time)
    this.update_time = new Timestamp(props.update_time)
    this.last_login_time = new Timestamp(props.last_login_time)
  }
}
