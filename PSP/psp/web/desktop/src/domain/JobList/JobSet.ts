/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action } from 'mobx'
import moment from 'moment'

export class JobSetModel {
  @observable project_id: string // 项目ID
  @observable project_name: string // 项目名称
  @observable job_set_id: string // 作业集ID
  @observable job_set_name: string // 作业集名称
  @observable job_type: string // 作业类型
  @observable app_id: string // 应用ID
  @observable app_name: string // 应用名称
  @observable user_id: string // 用户ID
  @observable user_name: string // 用户名称
  @observable job_count: number // 作业数量
  @observable success_count: number // 成功次数
  @observable failure_count: number // 失败次数
  @observable exec_duration: string // 作业执行时间
  @observable start_time: string // 开始时间
  @observable end_time: string // 结束时间
}

export type JobSetRequest = Omit<JobSetModel, 'start_time' | 'end_time'> & {
  start_time: string
  end_time: string
}

export class JobSet extends JobSetModel {
  constructor(props?: Partial<JobSetRequest>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update({ start_time, end_time, ...props }: Partial<JobSetRequest>) {
    Object.assign(this, props)

    if (start_time) {
      this.start_time = moment(start_time).format('YYYY-MM-DD HH:mm:ss')
    }
    if (end_time) {
      this.end_time = moment(end_time).format('YYYY-MM-DD HH:mm:ss')
    }
  }
}
