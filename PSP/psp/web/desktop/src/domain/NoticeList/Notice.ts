/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import moment, { Moment } from 'moment'

class BaseNotice {
  @observable id: string
  @observable title: string
  @observable content: string
  @observable start_time: Moment
  @observable end_time: Moment
  @observable published: boolean
  @observable creatpr: string
  @observable create_name: string
  @observable operator: string
  @observable operator_name: string
  @observable company_ids: string
  @observable priority: number // 0: '弹窗', 1: '横幅', 2: '悬浮窗' (暂时不支持悬浮窗)
}

export type NoticeRequest = Omit<BaseNotice, 'start_time' | 'end_time'> & {
  start_time: string
  end_time: string
}

export class Notice extends BaseNotice {
  constructor(props?: Partial<NoticeRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update({ start_time, end_time, ...props }: Partial<NoticeRequest>) {
    Object.assign(this, props)

    if (start_time) {
      this.start_time = moment(start_time)
    }

    if (end_time) {
      this.end_time = moment(end_time)
    }
  }
}
