/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import { Http, history } from '@/utils'
import { Timestamp } from '../common'

export enum InvitedUserType {
  // 未知
  INVITE_TO_UNKNOW = 0,
  // 非管理员
  INVITE_NOT_ADMIN = 1,
  // 管理员
  INVITE_IS_ADMIN = 2
}

interface IRequest {
  id: string
  company_id: string
  company_name: string
  real_name: string
  phone: string
  user_id: string
  is_admin: InvitedUserType
  status: number
  create_uid: string
  create_name: string
  create_time: {
    seconds: number
    nanos: number
  }
  update_time: {
    seconds: number
    nanos: number
  }
}

export class Invitation {
  @observable id
  @observable company_id
  @observable company_name
  @observable is_admin
  @observable real_name
  @observable phone
  @observable user_id
  @observable status
  @observable create_uid
  @observable create_name
  @observable create_time: Timestamp
  @observable update_time: Timestamp

  constructor(props: IRequest) {
    if (props) {
      this.init(props)
    }
  }

  @action
  init = (props: IRequest) => {
    Object.assign(this, props)
    this.create_time = new Timestamp(props.create_time)
    this.update_time = new Timestamp(props.update_time)
  }

  confirmInvite = async status => {
    await Http.put('/platform_user/confirm_invite', {
      invite_id: this.id,
      status
    })
  }

  accept = async () => {
    await this.confirmInvite(2)
    location.replace('/')
  }

  reject = () => this.confirmInvite(3)
}
