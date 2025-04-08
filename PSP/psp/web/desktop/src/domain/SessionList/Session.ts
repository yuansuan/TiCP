/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, computed } from 'mobx'
import { getDisplayRunTime } from '@/utils'
import { SESSION_STATUS_MAP } from '@/domain/Vis'
import { Software, Hardware } from '@/domain/VIsIBV'
import moment from 'moment'
export class BaseSession {
  @observable id: string
  @observable status: string
  @observable create_time?: string = '--'
  @observable start_time?: string = '--'
  @observable end_time?: string = '--'
  @observable update_time?: string = '--'
  @observable software: Software
  @observable hardware: Hardware
  @observable stream_url: string
  @observable exit_reason: string
  @observable duration: string
  @observable user_name: string
  @observable project_name: string
}

export class Session extends BaseSession {
  constructor(props) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update({ create_time, start_time, end_time, duration, ...props }) {
    Object.assign(this, props)
    if (create_time) {
      this.create_time = moment(create_time).format('YYYY-MM-DD HH:mm:ss')
    }
    if (start_time) {
      this.start_time = moment(start_time).format('YYYY-MM-DD HH:mm:ss')
    }
    if (end_time) {
      this.end_time = moment(end_time).format('YYYY-MM-DD HH:mm:ss')
    }

    if (duration) {
      this.duration = getDisplayRunTime(duration)
    }
  }

  @computed
  get status_str() {
    return SESSION_STATUS_MAP[this.status] || '--'
  }

  @computed
  get software_platform_str() {
    return this.software.platform || '--'
  }
}
