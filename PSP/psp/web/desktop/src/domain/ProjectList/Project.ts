/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { action, observable, runInAction } from 'mobx'
import { Timestamp } from '@/domain/common'
import { Http } from '@/utils'
import { env } from '@/domain'

type IRequest = Omit<BaseProject, 'update_time' | 'create_time'> & {
  update_time: {
    seconds: number
    nanos: number
  }
  create_time: {
    seconds: number
    nanos: number
  }
}

export class BaseProject {
  @observable id: string
  @observable name: string
  @observable creator: string
  @observable status: number
  @observable project_icon: string
  @observable type: number
  @observable company_id: string
  @observable account_id: string
  @observable user_id: string
  @observable remark: string
  @observable is_default: boolean
  @observable modify_uid: string
  @observable modify_name: string
  @observable create_name: string
  @observable create_uid: string
  @observable update_time: Timestamp
  @observable create_time: Timestamp
  @observable total_users: number
  @observable is_freeze: boolean
  @observable box_domain: string
  @observable consume_limit_amount: number
  @observable allow_submit_job_over_limit: boolean
  @observable consume_limit_enabled: boolean
  @observable owner_uid: string
  @observable owner_name: string
}

export default class Project extends BaseProject {
  constructor(props?: BaseProject) {
    super()
    Object.assign(this, props)
  }

  @action
  update = (props: Partial<IRequest>) => {
    Object.assign(this, props)

    if (props.create_time) {
      this.create_time = new Timestamp(props.create_time)
    }

    if (props.update_time) {
      this.update_time = new Timestamp(props.update_time)
    }
  }

  async updateRemark(remark) {
    await Http.put(`/project/modify_remark/${this.id}`, {
      remark,
    })

    runInAction(() => {
      this.update({
        remark,
      })
    })
  }

  async getUsers({ page_index, page_size, key = '' }) {
    return await Http.get(`/project/users_by_company`, {
      params: {
        page_index,
        page_size,
        key,
      },
    })
  }

  async modifyUsers(data: { add?: string[]; del?: string[] }) {
    if (!env.isPersonal) {
      await Http.put(`/project/modify_users_by_company/${this.id}`, data)
    } else {
      await Http.put(`/project/modify_users_by_user/${this.id}`, data)
    }
  }
}
